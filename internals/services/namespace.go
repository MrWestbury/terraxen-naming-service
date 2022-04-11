package services

import (
	"github.com/MrWestbury/terraxen-naming-service/internals/config"
	"github.com/google/uuid"
)

type Namespace struct {
	Id             string
	Name           string
	OrganizationId string
	Variables      map[string]string
	SchemaId       string
	SchemaVersion  string
}

type NamespaceService struct {
	tmpStore map[string]Namespace
}

func NewNamespaceService(config *config.Config) *NamespaceService {
	nssvc := &NamespaceService{
		tmpStore: make(map[string]Namespace),
	}
	return nssvc
}

func (nsSvc *NamespaceService) CreateNamespace(orgId string, name string, schemaId string, schemaVersion string, vars map[string]string) {
	ns := Namespace{
		Id:             uuid.NewString(),
		Name:           name,
		OrganizationId: orgId,
		Variables:      vars,
		SchemaId:       schemaId,
		SchemaVersion:  schemaVersion,
	}

	nsSvc.tmpStore[ns.Id] = ns
}

func (nsSvc *NamespaceService) GetNamespaceById(nsId string) (Namespace, error) {
	return nsSvc.tmpStore[nsId], nil
}
