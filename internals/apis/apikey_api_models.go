package apis

type NewApiKeyRequest struct {
	Name   string            `json:"name"`
	Scopes map[string]string `json:"scopes"`
}
