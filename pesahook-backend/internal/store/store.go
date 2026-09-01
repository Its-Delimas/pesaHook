package store

import (
	"errors"

	"github.com/Its-Delimas/pesahook/internal/endpoint"
	"github.com/Its-Delimas/pesahook/internal/event"
)

var ErrNotFound = errors.New("not found")

type EndpointStore interface {
	Save(e endpoint.Endpoint) error
	GetByID(id string) (endpoint.Endpoint, error)
	ListByAccount(accountID string) ([]endpoint.Endpoint, error)
}

type EventStore interface {
	Save(e event.NormalizedEvent) error
	GetByID(id string) (event.NormalizedEvent, error)
	ListByEndpoint(endpointID string) ([]event.NormalizedEvent, error)
}
