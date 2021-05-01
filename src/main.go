package main

import (
	"context"
	"log"
	"os"

	books "github.com/dmigo/books/books"
	health "github.com/dmigo/books/health"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var booksHandler *books.Handler
var healthHandler *health.Handler

func connectToStorage(ctx context.Context) (collection *mongo.Collection) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("[MongoDB] status: connected")
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("books")
	return
}

func connectToCache() (redisClient *redis.Client) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping()
	log.Printf("[Redis] %v", status)
	return
}

func init() {
	ctx := context.Background()

	collection := connectToStorage(ctx)
	redisClient := connectToCache()

	booksHandler = books.NewHandler(ctx, collection, redisClient)
	healthHandler = health.NewHandler(ctx)
}

func main() {
	router := gin.Default()
	router.GET("/health", healthHandler.Get)
	router.GET("/books", booksHandler.Get)
	router.POST("/books", booksHandler.Post)
	router.PUT("/books/:id/status", booksHandler.PutStatus)
	router.Run()
}
