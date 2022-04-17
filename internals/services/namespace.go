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

type Namespace struct {
	Id             string
	Name           string
	OrganizationId string
	Variables      map[string]string
	SchemaId       string
	SchemaVersion  string
}

var (
	ErrNamespaceAlreadyExists = errors.New("namespace with name already exists in organization")
)

type NamespaceServiceInterface interface {
	CreateNamespace(orgId string, name string, schemaId string, schemaVersion string, vars map[string]string)
	GetNamespaceById(nsId string) (Namespace, error)
	ListNamespaces(orgId string) ([]Namespace, error)
	ExistsByName(orgId, nsName string) bool
}

type NamespaceService struct {
	BaseService
}

func NewNamespaceService(config *config.Config) *NamespaceService {
	nssvc := &NamespaceService{}
	nssvc.Connect(config)
	nssvc.collection = nssvc.client.Collection("namespaces")
	return nssvc
}

func (nsSvc *NamespaceService) CreateNamespace(orgId string, name string, schemaId string, schemaVersion string, vars map[string]string) (*Namespace, error) {
	ns := &Namespace{
		Id:             uuid.NewString(),
		Name:           name,
		OrganizationId: orgId,
		Variables:      vars,
		SchemaId:       schemaId,
		SchemaVersion:  schemaVersion,
	}

	exists := nsSvc.ExistsByName(orgId, name)
	if exists {
		return nil, ErrNamespaceAlreadyExists
	}
	ctx := context.Background()
	_, err := nsSvc.collection.InsertOne(ctx, ns)
	if err != nil {
		log.Printf("Failed to create namespace %v", err)
		return nil, err
	}

	return ns, nil
}

func (nsSvc *NamespaceService) ExistsByName(orgId string, nsName string) bool {

	filter := bson.M{
		"name":           nsName,
		"organizationid": orgId,
	}

	ctx := context.Background()
	result := nsSvc.collection.FindOne(ctx, filter)

	if result.Err() == mongo.ErrNoDocuments {
		return false
	}

	return true
}

func (nsSvc *NamespaceService) GetNamespaceById(orgId string, nsId string) (Namespace, error) {

}

func (nsSvc *NamespaceService) ListNamespaces(orgId string) ([]Namespace, error) {

}
