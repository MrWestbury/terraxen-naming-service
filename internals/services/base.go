package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MrWestbury/terraxen-naming-service/internals/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseService struct {
	client     *mongo.Database
	collection *mongo.Collection
}

func (bs *BaseService) Connect(cfg *config.Config) {

	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", cfg.Username, cfg.Password, cfg.DBHost, cfg.DBName)

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(cfg.DBName)
	if db == nil {
		log.Panic(errors.New("database not selected"))
	}
	bs.client = db
}

func CloseCursor(ctx context.Context, cur *mongo.Cursor) {
	err := cur.Close(ctx)
	if err != nil {
		log.Printf("Failed closing cursor: %v", err)
	}
}
