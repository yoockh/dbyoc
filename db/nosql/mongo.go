package nosql

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBClient struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongoDBClient(uri, database, collection string) (*MongoDBClient, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	return &MongoDBClient{
		client:     client,
		database:   database,
		collection: collection,
	}, nil
}

func (m *MongoDBClient) Insert(document interface{}) error {
	collection := m.client.Database(m.database).Collection(m.collection)
	_, err := collection.InsertOne(context.Background(), document)
	return err
}

func (m *MongoDBClient) Find(filter interface{}) (*mongo.Cursor, error) {
	collection := m.client.Database(m.database).Collection(m.collection)
	return collection.Find(context.Background(), filter)
}

func (m *MongoDBClient) Update(filter interface{}, update interface{}) error {
	collection := m.client.Database(m.database).Collection(m.collection)
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (m *MongoDBClient) Close() error {
	return m.client.Disconnect(context.Background())
}