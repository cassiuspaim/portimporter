package domain

import "github.com/cassiuspaim/portimporter/domain/entities"

// Interface to define the operations for the PortService.
type PortService interface {
	Upsert(entities.Port) error
}

// Interface to define the operations for the PortRepository.
type PortRepository interface {
	GetByID(id string) (*entities.Port, error)
	Create(entities.Port) error
	Update(entities.Port, string) error
}
