package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client
var modelsCollection *mongo.Collection

// Middleware CORS sin encabezados, solo responde a OPTIONS
func corsPassthroughMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Configurar conexión a MongoDB
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("❌ Error cargando archivo .env: %v", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("❌ MONGO_URI no definido en .env")
	}

	mongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ Error al conectar a MongoDB: %v", err)
	}

	if err := mongoClient.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Fatalf("❌ Error de ping a MongoDB: %v", err)
	}

	modelsCollection = mongoClient.Database("CatalogServiceDB").Collection("models")
	log.Println("✅ Conectado a MongoDB correctamente")
}

// searchModelHandler maneja la búsqueda de modelos en MongoDB
func searchModelHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	name := queryParams.Get("name")
	createdBy := queryParams.Get("created_by")

	if name == "" && createdBy == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]bson.M{})
		return
	}

	filter := bson.M{}
	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if createdBy != "" {
		filter["created_by"] = bson.M{"$regex": createdBy, "$options": "i"}
	}

	cursor, err := modelsCollection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, fmt.Sprintf("❌ Error al buscar modelos: %v", err), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		http.Error(w, fmt.Sprintf("❌ Error al leer resultados: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/search", searchModelHandler).Methods("GET", "OPTIONS") // OJO: incluir OPTIONS

	// Envolver router con middleware de CORS vacío (solo permite preflight sin headers)
	handler := corsPassthroughMiddleware(r)

	fmt.Println("🚀 Microservicio de búsqueda iniciado en puerto 5005...")
	log.Fatal(http.ListenAndServe("0.0.0.0:5005", handler))
}
