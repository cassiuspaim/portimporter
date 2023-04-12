package domain

// TODO move this to package only for tests.
import (
	"errors"

	"github.com/cassiuspaim/portimporter/domain/entities"
)

// MockPortRepository used for tests.
type MockPortRepository struct {
	GetByIDfn func(id string) (*entities.Port, error)
	Updatefn  func(entities.Port, string) error
}

// Does what is defined at MockPortRepository.GetByIDfn.
// If MockPortRepository.GetByIDfn is not defined it retrieves an Error.
func (r MockPortRepository) GetByID(id string) (*entities.Port, error) {
	if r.GetByIDfn != nil {
		return r.GetByIDfn(id)
	}

	return nil, errors.New("No behaviour defined")
}

// Does what is defined at MockPortRepository.Updatefn.
// If MockPortRepository.Updatefn is not defined it retrieves an Error.
func (r MockPortRepository) Update(port entities.Port, filter string) error {
	if r.Updatefn != nil {
		return r.Updatefn(port, filter)
	}

	return errors.New("No behaviour defined")
}
