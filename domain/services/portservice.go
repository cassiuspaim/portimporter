package services

import (
	"fmt"
	"log"

	"github.com/cassiuspaim/portimporter/domain"
	"github.com/cassiuspaim/portimporter/domain/entities"
)

// PortService is a service that handle the business rules with Port entity.
type PortService struct {
	portRepository domain.PortRepository
}

// Retrieves a new PortService
func NewPortService(portRepository domain.PortRepository) PortService {
	return PortService{
		portRepository: portRepository,
	}
}

// Upsert a Port based on its ID
func (s PortService) Upsert(portEntity entities.Port) error {
	portDB, err := s.portRepository.GetByID(portEntity.ID)
	if err != nil {
		return err
	}

	if portDB != nil {
		err = s.portRepository.Update(portEntity, portEntity.ID)
		if err != nil {
			return fmt.Errorf("Error updating port %s. Error: %v", portEntity.ID, err)
		}

		log.Printf("Port updated %s.", portEntity.ID)
	} else {
		err = s.portRepository.Create(portEntity)
		if err != nil {
			return fmt.Errorf("Error creating port %s. Error: %v", portEntity.ID, err)
		}
		log.Printf("Port created %s.", portEntity.ID)
	}

	return nil
}
