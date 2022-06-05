package apis

type NewSchemaRequest struct {
	Name string `json:"name"`
}

type UpdateSchemaRequest struct {
	Name      string            `json:"name"`
	Resources map[string]string `json:"resources"`
}

type UpdateSchemaVersionRequest struct {
	Published bool              `json:"published"`
	Resources map[string]string `json:"resources"`
}

type CreateSchemaVersionRequest struct {
	FromVersion int               `json:"from_version"`
	Resources   map[string]string `json:"resources"`
}

type ResolveSchemaVersionRequest struct {
	ResouceName string            `json:"resource"`
	Variables   map[string]string `json:"variables"`
}
