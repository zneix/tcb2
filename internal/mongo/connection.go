package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zneix/tcb2/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoConnection creates a new instance of mongo.Connection. Keep in mind that Connect() has to be called before using it
func NewMongoConnection(bgctx context.Context, cfg *config.TCBConfig) *Connection {
	// Prepare mongo client's options
	uri := fmt.Sprintf("mongodb://%s:%s", "localhost", cfg.MongoPort)
	credentials := options.Credential{
		AuthSource: cfg.MongoAuthDB,
		Username:   cfg.MongoUsername,
		Password:   cfg.MongoPassword,
	}
	clientOptions := options.Client().ApplyURI(uri).SetAuth(credentials)

	// Actually connect to the database and test connection with a ping
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		log.Fatalf("[Mongo] Failed to create a new client: %s\n", err)
	}

	return &Connection{
		client:       client,
		ctx:          bgctx,
		databaseName: cfg.MongoDatabaseName,
	}
}

func (conn Connection) Connect() {
	ctx, cancel := context.WithTimeout(conn.ctx, 10*time.Second)
	defer cancel()

	err := conn.client.Connect(ctx)
	if err != nil {
		log.Fatalln("[Mongo] Error connecting: " + err.Error())
	}

	err = conn.client.Ping(ctx, nil)
	if err != nil {
		log.Fatalln("[Mongo] Error while executing the ping " + err.Error())
	}
	log.Println("[Mongo] connected")
}

func (conn Connection) Disconnect() {
	ctx, cancel := context.WithTimeout(conn.ctx, 10*time.Second)
	defer cancel()

	err := conn.client.Disconnect(ctx)
	if err != nil {
		log.Println("[Mongo] Error disconnecting: " + err.Error())
	}
	log.Println("[Mongo] disconnected")
}

// Collection returns a collection with given name from the bot's main database specified in the config
func (conn *Connection) Collection(name CollectionName) *mongo.Collection {
	return conn.client.Database(conn.databaseName).Collection(string(name))
}

// CollectionSubs returns a collection that stores user's subscriptions to SubEvent events
func (conn *Connection) CollectionSubs(channelID string) *mongo.Collection {
	return conn.client.Database(conn.databaseName + "-subscriptions").Collection(channelID)
}
