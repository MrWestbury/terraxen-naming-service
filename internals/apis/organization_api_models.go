package apis

type NewOrganizationRequest struct {
	Request bool
	Name    string `json:"name"`
}

type UpdateOrganizationRequest struct {
	Name      string            `json:"name"`
	Variables map[string]string `json:"variables"`
}
