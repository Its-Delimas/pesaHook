package store

import "github.com/Its-Delimas/pesahook/internal/apikey"

type APIKeyStore interface {
	Save(k apikey.APIKey) error
	GetByHash(hash string) (apikey.APIKey, error)
}
