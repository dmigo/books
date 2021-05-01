package books

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *Handler {
	return &Handler{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

func (handler *Handler) getBooksFromStorage() (*[]Book, error) {
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(handler.ctx)

	var books []Book
	err = cur.All(handler.ctx, &books)
	return &books, err
}

func (handler *Handler) getBooksFromCache() (*[]Book, error) {
	val, err := handler.redisClient.Get("books").Result()
	if err != nil {
		return nil, err
	}

	books := make([]Book, 0)
	json.Unmarshal([]byte(val), &books)

	return &books, nil
}

func (handler *Handler) saveBooksInCache(books *[]Book) error {
	data, _ := json.Marshal(books)
	if status := handler.redisClient.Set("books", string(data), 0); status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (handler *Handler) clearCache() error {
	if status := handler.redisClient.Del("books"); status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (handler *Handler) Get(c *gin.Context) {
	books, err := handler.getBooksFromCache()
	if err == redis.Nil {
		books, err := handler.getBooksFromStorage()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		handler.saveBooksInCache(books)
		c.JSON(http.StatusOK, books)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, books)
	}
}

func (handler *Handler) Post(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.ID = primitive.NewObjectID()
	_, err := handler.collection.InsertOne(handler.ctx, book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while adding a book"})
		return
	}

	handler.clearCache()

	c.JSON(http.StatusOK, book)
}

func (handler *Handler) PutStatus(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var newStatus NewStatus
	if err = c.ShouldBindJSON(&newStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = handler.collection.UpdateByID(handler.ctx, objectId, bson.D{
		{"$set", bson.M{"status": newStatus.Status}},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	handler.clearCache()

	c.JSON(http.StatusOK, newStatus)
}
