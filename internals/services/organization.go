package services

import (
	"context"
	"errors"
	"log"

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

type OrganizationServiceInterface interface {
	NewOrganization(organizationName string) (*Organization, error)
	GetOrganizationById(orgId string) (*Organization, error)
	ExistsById(orgId string) bool
	ExistsByName(orgName string) bool
	UpdateOrganization(orgId string, orgName string, orgVars map[string]string) (*Organization, error)
	DeleteOrganization(organizationId string) error
}

type OrganizationService struct {
	BaseService
}

func NewOrganizationService(cfg *config.Config) *OrganizationService {
	orgsvc := &OrganizationService{}
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
	return result.Err() != mongo.ErrNoDocuments
}

func (orgsvc *OrganizationService) ExistsByName(orgName string) bool {
	ctx := context.Background()
	result := orgsvc.collection.FindOne(ctx, bson.M{"name": orgName})
	return result.Err() != mongo.ErrNoDocuments
}

func (orgSvc *OrganizationService) UpdateOrganization(orgId string, orgName string, orgVars map[string]string) (*Organization, error) {
	org, err := orgSvc.GetOrganizationById(orgId)
	if err != nil {
		log.Printf("Failed to update organization: %v", err)
		return nil, err
	}

	org.Name = orgName
	org.OrgVars = orgVars

	ctx := context.Background()

	filter := bson.M{
		"id": org.Id,
	}

	orgSvc.collection.FindOneAndUpdate(ctx, filter, org)
	return org, nil
}

func (orgSvc OrganizationService) DeleteOrganization(organizationId string) error {
	org, err := orgSvc.GetOrganizationById(organizationId)
	if err != nil {
		log.Printf("Failed to update organization: %v", err)
		return err
	}

	ctx := context.Background()

	filter := bson.M{
		"id": org.Id,
	}

	result, err := orgSvc.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no organization deleted")
	}
	return nil
}
