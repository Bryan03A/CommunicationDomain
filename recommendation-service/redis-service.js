const express = require("express");
const { createClient } = require("redis");
const cors = require("cors");

const app = express();
const port = 5006;

// Configurar Redis con credenciales en la nube
const redisClient = createClient({
    username: "default",
    password: "TD0nSdyb1joTFrilTfO5xZ2L6rmFqvNp",
    socket: {
        host: "redis-15780.c14.us-east-1-2.ec2.redns.redis-cloud.com",
        port: 15780
    }
});

// Conectar a Redis con manejo de errores
(async () => {
    try {
        await redisClient.connect();
        console.log("✅ Conectado a Redis Cloud");
    } catch (err) {
        console.error("❌ Error al conectar a Redis:", err);
    }
})();

// Middleware para manejar errores de Redis
redisClient.on("error", (err) => console.error("❌ Redis Error:", err));

// Configurar CORS
app.use(cors({
    origin: "*",  // Permitir solicitudes de cualquier origen
    methods: ["GET", "POST"],
    allowedHeaders: ["Content-Type", "Authorization"]
}));

// Configurar el body parser
app.use(express.json());

// Ruta de Health Check
app.get("/health", (req, res) => {
    res.status(200).send("Healthy");
});

// Guardar búsqueda en Redis para un usuario específico
app.post("/save-search", async (req, res) => {
    const { query, creator, username, firstModelName } = req.body;

    if (!query || !username) {
        return res.status(400).send("Falta la consulta de búsqueda o el nombre de usuario");
    }

    try {
        const userSearchKey = `search-history:${username}`;
        const searchData = {
            query,
            creator,
            firstModelName,
            timestamp: Date.now()
        };
        await redisClient.lPush(userSearchKey, JSON.stringify(searchData));

        res.status(200).send("Búsqueda guardada en Redis");
    } catch (error) {
        console.error("Error al guardar búsqueda en Redis:", error);
        res.status(500).send("Error al guardar búsqueda");
    }
});

// Obtener las últimas búsquedas guardadas en Redis para un usuario específico
app.get("/recent-searches", async (req, res) => {
    const { username } = req.query;

    if (!username) {
        return res.status(400).send("Falta el nombre de usuario");
    }

    try {
        const userSearchKey = `search-history:${username}`;
        const searches = await redisClient.lRange(userSearchKey, 0, 9);
        const parsedSearches = searches.map((search) => JSON.parse(search));

        res.json(parsedSearches);
    } catch (error) {
        console.error("Error al obtener búsquedas recientes:", error);
        res.status(500).send("Error al obtener búsquedas recientes");
    }
});

// Iniciar el servidor en todas las interfaces
app.listen(port, '0.0.0.0', () => {
    console.log(`Servidor Redis-Service escuchando en http://34.200.3.211:${port}`);
});