package apis

type NewNamespaceRequest struct {
	OrgId         string `json:"organization_id"`
	Name          string `json:"name"`
	Schema        string `json:"schema_id"`
	SchemaVersion string `json:"schema_version"`
}

type ResolveResourceResponse struct {
	ResourceName string `json:"name"`
	Pattern      string `json:"pattern"`
	Value        string `json:"value"`
}
