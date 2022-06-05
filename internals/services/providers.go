package services

type ApiKeyProvider interface {
	GenerateNewApiKey(orgId string) (*ApiKey, error)
	ListKeys(orgId string) ([]*ApiKey, error)
	GetKey(key string) *ApiKey
	DeleteKey(apiId string) error
}

type NamespaceServiceProvider interface {
	CreateNamespace(orgId string, name string, schemaId string, schemaVersion string, vars map[string]string) (*Namespace, error)
	GetNamespaceById(orgId string, nsId string) (*Namespace, error)
	ListNamespaces(orgId string) ([]*Namespace, error)
	ExistsByName(orgId, nsName string) bool
	UpdateNamespace(orgId string, nsId string, nsName string, schemaVersion string) error
	DeleteNamespace(orgId string, nsId string) error
	ListNamespaceVars(orgId string, nsId string) ([]*NamespaceVar, error)
	GetNamespaceVariable(orgId string, nsId string, varId string) (*NamespaceVar, error)
	CreateNamespaceVariable(orgId string, nsId string, varId string, value string) (*NamespaceVar, error)
	UpdateNamespaceVariable(orgId string, nsId string, varId string, value string) (*NamespaceVar, error)
	DeleteNamespaceVariable(orgId string, nsId string, varId string) error
	GetVariablesAsMap(orgId string, nsId string) (map[string]string, error)
	NamespaceVariableExists(orgId string, nsId string, varId string) bool
}

type OrganizationServiceProvider interface {
	NewOrganization(organizationName string) (*Organization, error)
	GetOrganizationById(orgId string) (*Organization, error)
	ExistsById(orgId string) bool
	ExistsByName(orgName string) bool
	UpdateOrganization(orgId string, orgName string, orgVars map[string]string) (*Organization, error)
	DeleteOrganization(organizationId string) error
}

type SchemaServiceProvider interface {
	CreateSchema(orgId string, name string) (*Schema, error)
	ListSchemaInOrganization(orgId string) ([]*Schema, error)
	GetSchemaById(orgId string, schemaId string) (*Schema, error)
	UpdateSchema(schema Schema) error
	DeleteSchema(orgId string, schemaId string) error
	ListSchemaVersions(orgId string, schemaId string) ([]*SchemaVersion, error)
	CreateSchemaVersion(orgId string, schemaId string, resources map[string]string, published bool) (*SchemaVersion, error)
	GetSchemaVersion(orgId string, schemaId string, schemaVersionId string) (*SchemaVersion, error)
	UpdateSchemaVersion(orgId string, schemaId string, schemaVersionId string, resources map[string]string, published bool) (*SchemaVersion, error)
}
