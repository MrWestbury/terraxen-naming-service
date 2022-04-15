package services

import (
	"context"
	"errors"
	"log"

	"github.com/MrWestbury/terraxen-naming-service/internals/config"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	Id        int                 `json:"id"`
	SchemaId  string              `json:"schema_id"`
	Resources map[string]Resource `json:"resources"`
}

type SchemaService struct {
	BaseService
	versionCollection *mongo.Collection
}

func NewSchemaService(cfg *config.Config) *SchemaService {
	sSvc := &SchemaService{}
	sSvc.Connect(cfg)
	sSvc.collection = sSvc.client.Collection("schemas")
	sSvc.versionCollection = sSvc.client.Collection("schemaversions")

	return sSvc
}

func (sSvc *SchemaService) CreateSchema(orgId string, name string) (*Schema, error) {
	schema, err := sSvc.getSchemaByName(orgId, name)
	if err != nil {
		return nil, err
	}

	if schema != nil {
		return nil, errors.New("schema already exists")
	}

	newSchema := &Schema{
		Id:             uuid.NewString(),
		OrganizationId: orgId,
		Name:           name,
	}

	newSchemaVersion := &SchemaVersion{
		Id:        1,
		SchemaId:  newSchema.Id,
		Resources: make(map[string]Resource),
	}

	ctx := context.Background()
	result, err := sSvc.collection.InsertOne(ctx, newSchema)
	if err != nil {
		return nil, err
	}

	_, err = sSvc.versionCollection.InsertOne(ctx, newSchemaVersion)
	if err != nil {
		sSvc.collection.FindOneAndDelete(ctx, bson.M{"_id": result.InsertedID}) // Cleanup
		return nil, err
	}

	return newSchema, nil
}

func (sSvc *SchemaService) GetSchemaVersionById(orgId string, schemaId string, schemaVersionId string) (*SchemaVersion, error) {
	ctx := context.Background()

	filter := bson.M{
		"id":             schemaId,
		"organizationid": orgId,
		"schema_id":      schemaId,
	}

	opts := options.FindOne()
	opts.SetSort(bson.D{primitive.E{Key: "id", Value: 1}})

	result := sSvc.collection.FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	schemaVersion := &SchemaVersion{}
	result.Decode(schemaVersion)
	return schemaVersion, nil
}

func (sSvc *SchemaService) getSchemaByName(orgId string, name string) (*Schema, error) {
	ctx := context.Background()
	filter := bson.M{
		"organizationid": orgId,
		"name":           name,
	}
	result := sSvc.collection.FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	} else if result.Err() != nil {
		log.Printf("Failed to get schema by name: %v", result.Err())
		return nil, result.Err()
	}

	var sch *Schema
	err := result.Decode(sch)
	if err != nil {
		log.Printf("Failed to get schema by name: %v", err)
		return nil, err
	}

	return sch, nil
}

func (sSvc *SchemaService) ListSchemaInOrganization(orgId string) ([]*Schema, error) {
	ctx := context.Background()
	filter := bson.M{
		"organizationid": orgId,
	}

	opts := options.Find()
	opts.SetLimit(50)
	opts.SetSort(bson.D{primitive.E{Key: "name", Value: 1}})

	cur, err := sSvc.collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Error finding schemas in org: %v", err)
		return nil, err
	}
	defer CloseCursor(ctx, cur)

	results := make([]*Schema, 0)
	for cur.Next(ctx) {
		var schema Schema
		err = cur.Decode(&schema)
		if err != nil {
			log.Printf("Unable to decode schema: %v", err)
			continue
		}
		results = append(results, &schema)
	}

	return results, nil
}

func (sSvc *SchemaService) GetSchemaById(orgId string, schemaId string) (*Schema, error) {
	ctx := context.Background()

	filter := bson.M{
		"organizationid": orgId,
		"schema_id":      schemaId,
	}

	result := sSvc.collection.FindOne(ctx, filter)
	if result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	} else if result.Err() != nil {
		log.Printf("Failed to get schema by name: %v", result.Err())
		return nil, result.Err()
	}

	var sch *Schema
	err := result.Decode(sch)
	if err != nil {
		log.Printf("Failed to get schema by name: %v", err)
		return nil, err
	}

	return sch, nil

}

func (sSvc *SchemaService) UpdateSchema(schema Schema) error {
	ctx := context.Background()

	filter := bson.M{
		"organizationid": schema.OrganizationId,
		"schema_id":      schema.Id,
	}

	results := sSvc.collection.FindOneAndReplace(ctx, filter, schema)
	if results.Err() != nil {
		log.Printf("failed to update schema: %v", results.Err())
		return results.Err()
	}
	return nil
}

func (sSvc *SchemaService) DeleteSchema(orgId string, schemaId string) error {
	ctx := context.Background()

	filter := bson.M{
		"organizationid": orgId,
		"schema_id":      schemaId,
	}

	result, err := sSvc.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("failed to delete schema: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("schema not found")
	}

	versionFilter := bson.M{
		"schema_id": schemaId,
	}

	_, err = sSvc.versionCollection.DeleteMany(ctx, versionFilter)
	if err != nil {
		log.Printf("failed to delete schema version: %v", err)
		return err
	}

	return nil
}
