package services

import (
	"context"
	"errors"

	"github.com/MrWestbury/terraxen-naming-service/internals/config"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Organization struct {
	Id      string            `json:"id"`
	Name    string            `json:"name"`
	OrgVars map[string]string `json:"vars"`
}

type OrganizationService struct {
	BaseService
	tmpStore map[string]Organization
}

func NewOrganizationService(cfg *config.Config) *OrganizationService {
	orgsvc := &OrganizationService{
		tmpStore: make(map[string]Organization),
	}
	orgsvc.Connect(cfg)
	orgsvc.collection = orgsvc.client.Collection("organizations")

	return orgsvc
}

func (orgsvc *OrganizationService) NewOrganization(orgName string) (*Organization, error) {
	org, err := orgsvc.getOrganizationByName(orgName)
	if err != nil {
		return nil, err
	}

	if org != nil {
		return nil, errors.New("organization already exists")
	}

	newOrg := &Organization{
		Id:      uuid.NewString(),
		Name:    orgName,
		OrgVars: make(map[string]string),
	}

	ctx := context.Background()
	_, err = orgsvc.collection.InsertOne(ctx, newOrg)
	if err != nil {
		return nil, err
	}
	return newOrg, err
}

func (orgsvc *OrganizationService) getOrganizationByName(orgName string) (*Organization, error) {
	ctx := context.Background()
	result := orgsvc.collection.FindOne(ctx, bson.M{"name": orgName})
	if result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	org := &Organization{}
	result.Decode(org)
	return org, nil
}

func (orgsvc *OrganizationService) GetOrganizationById(orgId string) (*Organization, error) {
	ctx := context.Background()
	result := orgsvc.collection.FindOne(ctx, bson.M{"id": orgId})
	if result.Err() == mongo.ErrNoDocuments {
		return nil, nil
	}
	org := &Organization{}
	result.Decode(org)
	return org, nil
}

func (orgsvc *OrganizationService) ExistsById(orgId string) bool {
	ctx := context.Background()
	result := orgsvc.collection.FindOne(ctx, bson.M{"id": orgId})
	if result.Err() == mongo.ErrNoDocuments {
		return false
	}
	return true
}

func (orgsvc *OrganizationService) ExistsByName(orgName string) bool {
	ctx := context.Background()
	result := orgsvc.collection.FindOne(ctx, bson.M{"name": orgName})
	if result.Err() == mongo.ErrNoDocuments {
		return false
	}
	return true
}
