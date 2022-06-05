package services

import "errors"

var (
	ErrNamespaceAlreadyExists = errors.New("namespace with name already exists in organization")
	ErrNamespaceNotFound      = errors.New("namespace not found")
	ErrSchemaNotFound         = errors.New("schema not found")
)
