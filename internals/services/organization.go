package services

import (
	"errors"

	"github.com/google/uuid"
)

type Organisation struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type OrganizationService struct {
	tmpStore map[string]Organisation
}

func NewOrganizationService() *OrganizationService {
	orgsvc := &OrganizationService{
		tmpStore: make(map[string]Organisation),
	}

	return orgsvc
}

func (orgsvc *OrganizationService) NewOrganization(orgName string) (*Organisation, error) {
	org, err := orgsvc.getOrganizationByName(orgName)
	if err != nil {
		return nil, err
	}

	if org != nil {
		return nil, errors.New("organization already exists")
	}

	newOrg := &Organisation{
		Id:   uuid.NewString(),
		Name: orgName,
	}

	orgsvc.tmpStore[newOrg.Id] = *newOrg
	return newOrg, nil
}

func (orgsvc *OrganizationService) getOrganizationByName(orgName string) (*Organisation, error) {
	for _, o := range orgsvc.tmpStore {
		if o.Name == orgName {
			return &o, nil
		}
	}
	return nil, nil
}
