package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client
var modelsCollection *mongo.Collection

// Configurar conexión a MongoDB
func init() {
	mongoURI := "mongodb+srv://MicroserviceDev:1997999@cluster0.hdqpd.mongodb.net/CatalogServiceDB?retryWrites=true&w=majority"

	var err error
	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}

	// Verificar la conexión
	if err := mongoClient.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatalf("Error de ping a MongoDB: %v", err)
	}

	// Conectar a la colección 'models'
	modelsCollection = mongoClient.Database("CatalogServiceDB").Collection("models")

	log.Println("Conectado a MongoDB Atlas correctamente")
}

// searchModelHandler maneja la búsqueda de modelos en MongoDB
func searchModelHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	name := queryParams.Get("name")
	createdBy := queryParams.Get("created_by")

	// Si ambos campos están vacíos, no devolvemos nada
	if name == "" && createdBy == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]bson.M{}) // Devuelve una lista vacía
		return
	}

	// Construir el filtro de búsqueda flexible
	filter := bson.M{}
	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"} // Búsqueda insensible a mayúsculas/minúsculas
	}
	if createdBy != "" {
		filter["created_by"] = bson.M{"$regex": createdBy, "$options": "i"} // También búsqueda flexible por creador
	}

	// Buscar en MongoDB
	cursor, err := modelsCollection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al buscar modelos: %v", err), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	// Leer resultados
	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		http.Error(w, fmt.Sprintf("Error al leer resultados: %v", err), http.StatusInternalServerError)
		return
	}

	// Devolver la respuesta como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	// Crear el enrutador
	r := mux.NewRouter()

	// Definir los endpoints
	r.HandleFunc("/search", searchModelHandler).Methods("GET")

	// Configurar CORS para permitir solicitudes desde la interfaz
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://3.212.132.24:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Iniciar el servidor con CORS habilitado
	fmt.Println("Microservicio de búsqueda iniciado en puerto 5005 con CORS y búsqueda flexible...")
	log.Fatal(http.ListenAndServe("0.0.0.0:5005", corsHandler.Handler(r)))
}
