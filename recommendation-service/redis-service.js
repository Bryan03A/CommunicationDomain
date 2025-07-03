require("dotenv").config(); // ✅ Cargar variables del .env

const express = require("express");
const { createClient } = require("redis");

const app = express();
const port = 5006;

// Configurar Redis desde .env
const redisClient = createClient({
    username: process.env.REDIS_USERNAME || undefined,
    password: process.env.REDIS_PASSWORD || undefined,
    socket: {
        host: process.env.REDIS_HOST,
        port: Number(process.env.REDIS_PORT)
    }
});

// Conectar a Redis
(async () => {
    try {
        await redisClient.connect();
        console.log("✅ Conectado a Redis");
    } catch (err) {
        console.error("❌ Error al conectar a Redis:", err);
    }
})();

// Manejar errores de Redis
redisClient.on("error", (err) => console.error("❌ Redis Error:", err));

app.use(express.json());

// Health Check
app.get("/health", (req, res) => {
    res.status(200).send("Healthy");
});

// Guardar búsqueda
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

// Obtener búsquedas recientes
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

// Iniciar servidor
app.listen(port, '0.0.0.0', () => {
    console.log(`🚀 Servidor Redis-Service escuchando en http://0.0.0.0:${port}`);
});