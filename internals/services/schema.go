package services

import (
	"context"

	"github.com/MrWestbury/terraxen-naming-service/internals/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Schema struct {
	Id             string
	OrganizationId string
	Name           string
}

type Resource struct {
	Name    string `json:"name"`
	Pattern string `json:"pattern"`
}

type SchemaVersion struct {
	Id        string              `json:"id"`
	SchemaId  string              `json:"schema_id"`
	Resources map[string]Resource `json:"resources"`
}

type SchemaService struct {
	BaseService
}

func NewSchemaService(cfg *config.Config) *SchemaService {
	sSvc := &SchemaService{}
	sSvc.Connect(cfg)
	sSvc.collection = sSvc.client.Collection("schemas")

	return sSvc
}

func (sSvc *SchemaService) GetSchemaVersionById(schemaId string, schemaVersionId string) (*SchemaVersion, error) {
	ctx := context.Background()
	result := sSvc.collection.FindOne(ctx, bson.M{"id": schemaVersionId, "schema_id": schemaId})
	if result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	schemaVersion := &SchemaVersion{}
	result.Decode(schemaVersion)
	return schemaVersion, nil

}
