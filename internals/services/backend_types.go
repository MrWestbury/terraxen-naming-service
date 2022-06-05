package services

import "time"

type ApiKey struct {
	Id             string            `json:"id"`
	OrganizationId string            `json:"organization_id"`
	Name           string            `json:"name"`
	Key            string            `json:"key"`
	Expires        time.Time         `json:"expires"`
	Scope          map[string]string `json:"scope"`
}

type Namespace struct {
	Id             string
	Name           string
	OrganizationId string
	SchemaId       string
	SchemaVersion  string
}

type NamespaceVar struct {
	Key         string
	Value       string
	OrgId       string
	NamespaceId string
}

type Organization struct {
	Id      string            `json:"id"`
	Name    string            `json:"name"`
	OrgVars map[string]string `json:"vars"`
}

type OrganizationVar struct {
	Id    string `json:"id"`
	OrgId string `json:"orgid"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Schema struct {
	Id             string
	OrganizationId string
	Name           string
}

type SchemaVersion struct {
	Id        int               `json:"id"`
	Published bool              `json:"published"`
	SchemaId  string            `json:"schema_id"`
	Resources map[string]string `json:"resources"`
}
