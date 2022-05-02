package apis

type NewNamespaceRequest struct {
	Name          string `json:"name"`
	Schema        string `json:"schema_id"`
	SchemaVersion string `json:"schema_version"`
}

type UpdateNamespaceRequest struct {
	Name          string `json:"name"`
	SchemaVersion string `json:"schema_version"`
}

type ResolveResourceResponse struct {
	ResourceName string `json:"name"`
	Pattern      string `json:"pattern"`
	Value        string `json:"value"`
}

type NewNamespaceVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type UpdateNamespaceVariable struct {
	Value string `json:"value"`
}
