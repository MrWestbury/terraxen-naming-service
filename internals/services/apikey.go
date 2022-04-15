package services

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/MrWestbury/terraxen-naming-service/internals/config"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApiKey struct {
	Id             string            `json:"id"`
	OrganizationId string            `json:"organization_id"`
	Name           string            `json:"name"`
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
	randStr := RandString(20)
	duration, _ := time.ParseDuration("1h")
	expires := time.Now().Add(duration)
	ak := &ApiKey{
		Id:             uuid.NewString(),
		OrganizationId: orgId,
		Key:            fmt.Sprintf("API-%s", randStr),
		Expires:        expires,
		Scope:          make(map[string]string),
	}

	ctx := context.Background()
	akSvc.collection.InsertOne(ctx, ak)

	return ak, nil
}

func RandString(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	letterBytes := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[src.Int63()%int64(len(letterBytes))]
	}
	return string(b)

}

func (akSvc *ApiKeyService) ListKeys(orgId string) ([]*ApiKey, error) {
	ctx := context.Background()
	filter := bson.D{
		primitive.E{Key: "organizationid", Value: orgId},
	}
	cur, err := akSvc.collection.Find(ctx, filter)
	defer CloseCursor(ctx, cur)

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

func (akSvc *ApiKeyService) DeleteKey(apiId string) error {
	ctx := context.Background()

	filter := bson.M{"id": apiId}
	result := akSvc.collection.FindOneAndDelete(ctx, filter)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
