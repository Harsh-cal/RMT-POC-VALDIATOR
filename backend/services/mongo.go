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

func normalizeValidationResultDoc(doc bson.M) bson.M {
	if _, ok := doc["release_id"]; !ok {
		if v, exists := doc["releaseid"]; exists {
			doc["release_id"] = v
		}
	}

	if _, ok := doc["release_name"]; !ok {
		if v, exists := doc["releasename"]; exists {
			doc["release_name"] = v
		}
	}

	if _, ok := doc["target_fleet"]; !ok {
		if v, exists := doc["targetfleet"]; exists {
			doc["target_fleet"] = v
		}
	}

	if _, ok := doc["validated_at"]; !ok {
		if v, exists := doc["validatedat"]; exists {
			doc["validated_at"] = v
		}
	}

	if rawRecs, ok := doc["recommendations"].(bson.A); ok {
		for i, item := range rawRecs {
			rec, ok := item.(bson.M)
			if !ok {
				continue
			}
			if _, hasIssueType := rec["issue_type"]; !hasIssueType {
				if legacy, exists := rec["issuetype"]; exists {
					rec["issue_type"] = legacy
				}
			}
			rawRecs[i] = rec
		}
		doc["recommendations"] = rawRecs
	}

	return doc
}

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

	doc := bson.M{
		"release_id":      result.ReleaseID,
		"release_name":    result.ReleaseName,
		"version":         result.Version,
		"target_fleet":    result.TargetFleet,
		"risk":            result.Risk,
		"status":          result.Status,
		"issues":          result.Issues,
		"insight":         result.Insight,
		"recommendations": result.Recommendations,
		"validated_at":    result.ValidatedAt,
	}

	_, err := collection.InsertOne(ctx, doc)
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
		"aircraft":     release.Aircraft,
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

	var doc bson.M
	filter := bson.M{
		"$or": bson.A{
			bson.M{"release_id": releaseID},
			bson.M{"releaseid": releaseID},
		},
	}
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return nil, err
	}

	doc = normalizeValidationResultDoc(doc)

	raw, err := bson.Marshal(doc)
	if err != nil {
		return nil, err
	}

	var result models.ValidationResult
	if err := bson.Unmarshal(raw, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetLatestValidationResultByReleaseName retrieves the latest validation result by release name.
func GetLatestValidationResultByReleaseName(releaseName string) (*models.ValidationResult, error) {
	if MongoDatabase == nil {
		return nil, fmt.Errorf("MongoDB not initialized")
	}

	collection := MongoDatabase.Collection("validation_results")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var doc bson.M
	opts := options.FindOne().SetSort(bson.D{{Key: "validated_at", Value: -1}, {Key: "validatedat", Value: -1}})
	filter := bson.M{
		"$or": bson.A{
			bson.M{"release_name": releaseName},
			bson.M{"releasename": releaseName},
		},
	}
	err := collection.FindOne(ctx, filter, opts).Decode(&doc)
	if err != nil {
		return nil, err
	}

	doc = normalizeValidationResultDoc(doc)

	raw, err := bson.Marshal(doc)
	if err != nil {
		return nil, err
	}

	var result models.ValidationResult
	if err := bson.Unmarshal(raw, &result); err != nil {
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
