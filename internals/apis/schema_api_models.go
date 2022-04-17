package apis

type NewSchemaRequest struct {
	Name string `json:"name"`
}

type UpdateSchemaRequest struct {
	Name string `json:"name"`
}

type ResolveSchemaVersionRequest struct {
	ResouceName string            `json:"resource"`
	Varibales   map[string]string `json:"variables"`
}
