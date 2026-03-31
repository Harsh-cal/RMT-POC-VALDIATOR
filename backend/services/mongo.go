package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/harsh-cal/rmt-poc-validator/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

// InitMongo initializes MongoDB connection
func InitMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background()) // cleanup on failure
		return err
	}

	MongoClient = client

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "rmt_validator"
	}

	MongoDatabase = client.Database(dbName)

	fmt.Println("MongoDB connected successfully")
	return nil
}

// CloseMongo closes MongoDB connection
func CloseMongo() error {
	if MongoClient != nil {
		return MongoClient.Disconnect(context.Background())
	}
	return nil
}

// SaveValidationResult saves validation result to MongoDB
func SaveValidationResult(result *models.ValidationResult) error {
	if MongoDatabase == nil {
		return fmt.Errorf("MongoDB not initialized")
	}

	collection := MongoDatabase.Collection("validation_results")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, result)
	return err
}

// SaveRelease saves release snapshot to MongoDB
func SaveRelease(release *models.ValidateRequest, releaseID string) error {
	if MongoDatabase == nil {
		return fmt.Errorf("MongoDB not initialized")
	}

	collection := MongoDatabase.Collection("releases")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := bson.M{
		"release_id":   releaseID,
		"release_name": release.ReleaseName,
		"version":      release.Version,
		"target_fleet": release.TargetFleet,
		"containers":   release.Containers,
		"created_at":   time.Now(),
	}

	_, err := collection.InsertOne(ctx, doc)
	return err
}

// GetValidationResult retrieves a validation result by release ID
func GetValidationResult(releaseID string) (*models.ValidationResult, error) {
	if MongoDatabase == nil {
		return nil, fmt.Errorf("MongoDB not initialized")
	}

	collection := MongoDatabase.Collection("validation_results")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result models.ValidationResult
	err := collection.FindOne(ctx, bson.M{"release_id": releaseID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAllValidationResults retrieves all validation results
func GetAllValidationResults() ([]models.ValidationResult, error) {
	if MongoDatabase == nil {
		return nil, fmt.Errorf("MongoDB not initialized")
	}

	collection := MongoDatabase.Collection("validation_results")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.ValidationResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}