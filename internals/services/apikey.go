package services

import (
	"context"
	"fmt"
	"time"

	"github.com/MrWestbury/terraxen-naming-service/internals/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApiKey struct {
	OrganizationId string            `json:"organization_id"`
	Key            string            `json:"key"`
	Expires        time.Time         `json:"expires"`
	Scope          map[string]string `json:"scope"`
}

type ApiKeyService struct {
	BaseService
}

func NewApiKeyService(cfg *config.Config) *ApiKeyService {
	apiKeySvc := &ApiKeyService{}
	apiKeySvc.Connect(cfg)
	apiKeySvc.collection = apiKeySvc.client.Collection("api_keys")
	return apiKeySvc
}

func (akSvc *ApiKeyService) GenerateNewApiKey(orgId string) (*ApiKey, error) {
	randStr := "asdasds"
	duration, _ := time.ParseDuration("1h")
	expires := time.Now().Add(duration)
	ak := &ApiKey{
		OrganizationId: orgId,
		Key:            fmt.Sprintf("API-%s", randStr),
		Expires:        expires,
		Scope:          make(map[string]string),
	}

	ctx := context.Background()
	akSvc.collection.InsertOne(ctx, ak)

	return ak, nil
}

func (akSvc *ApiKeyService) ListKeys(orgId string) ([]*ApiKey, error) {
	ctx := context.Background()
	filter := bson.D{
		primitive.E{Key: "organization_id", Value: orgId}, //TODO: fix this filter. cur.Next is returning nothing
	}
	cur, err := akSvc.collection.Find(ctx, filter)
	defer cur.Close(ctx)

	if err != nil {
		return nil, err
	}
	var result []*ApiKey
	for cur.Next(ctx) {
		var apiKey ApiKey
		err = cur.Decode(&apiKey)
		if err != nil {
			return nil, err
		}
		result = append(result, &apiKey)
	}

	return result, nil
}

func (akSvc *ApiKeyService) GetKey(key string) *ApiKey {
	ctx := context.Background()
	result := akSvc.collection.FindOne(ctx, bson.M{"key": key})
	if result.Err() == mongo.ErrNoDocuments {
		return nil
	}
	apiKey := &ApiKey{}
	result.Decode(apiKey)
	return apiKey
}