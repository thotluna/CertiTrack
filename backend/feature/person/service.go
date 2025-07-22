package person

import (
	"certitrack/backend/feature/person/dto"
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreatePerson(ctx context.Context, dto *dto.CreatePersonRequest) (*Person, error)
	GetPerson(ctx context.Context, id uuid.UUID) (*Person, error)
	ListPersons(ctx context.Context, filters map[string]interface{}) ([]*Person, error)
	CountPersons(ctx context.Context, filters map[string]interface{}) (int64, error)
	UpdatePerson(ctx context.Context, id uuid.UUID, dto *dto.UpdatePersonRequest) (*Person, error)
	DeletePerson(ctx context.Context, id uuid.UUID) error
	ChangePersonStatus(ctx context.Context, id uuid.UUID, status string) (*Person, error)
}

type service struct {
	repo PersonRepository
}

func NewService(repo PersonRepository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePerson(ctx context.Context, dto *dto.CreatePersonRequest) (*Person, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	exists, err := s.emailExists(ctx, dto.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	person := NewPerson(
		dto.FirstName,
		dto.LastName,
		dto.Email,
		dto.Phone,
	)

	if err := s.repo.Create(ctx, person); err != nil {
		return nil, err
	}

	return person, nil
}

func (s *service) GetPerson(ctx context.Context, id uuid.UUID) (*Person, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) ListPersons(ctx context.Context, filters map[string]interface{}) ([]*Person, error) {
	cleanFilters := make(map[string]interface{})
	for k, v := range filters {
		if str, ok := v.(string); ok {
			cleanFilters[k] = strings.TrimSpace(str)
		} else {
			cleanFilters[k] = v
		}
	}

	return s.repo.FindAll(ctx, cleanFilters)
}

func (s *service) UpdatePerson(ctx context.Context, id uuid.UUID, dto *dto.UpdatePersonRequest) (*Person, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	person, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if dto.FirstName != nil {
		person.FirstName = *dto.FirstName
	}

	if dto.LastName != nil {
		person.LastName = *dto.LastName
	}

	if dto.Phone != nil {
		person.Phone = *dto.Phone
	}

	if dto.Email != nil && *dto.Email != person.Email {
		exists, err := s.emailExists(ctx, *dto.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrEmailAlreadyExists
		}
		person.Email = *dto.Email
	}

	person.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, person); err != nil {
		return nil, err
	}

	return person, nil
}

func (s *service) DeletePerson(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) CountPersons(ctx context.Context, filters map[string]interface{}) (int64, error) {
	// Eliminar claves de paginación si existen
	cleanFilters := make(map[string]interface{})
	for k, v := range filters {
		if k != "limit" && k != "offset" {
			cleanFilters[k] = v
		}
	}

	return s.repo.Count(ctx, cleanFilters)
}

func (s *service) ChangePersonStatus(ctx context.Context, id uuid.UUID, status string) (*Person, error) {
	if !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	person, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	person.ChangeStatus(status)

	if err := s.repo.Update(ctx, person); err != nil {
		return nil, err
	}

	return person, nil
}

func (s *service) emailExists(ctx context.Context, email string) (bool, error) {
	filters := map[string]interface{}{"email": strings.ToLower(strings.TrimSpace(email))}
	persons, err := s.repo.FindAll(ctx, filters)
	if err != nil {
		return false, err
	}
	return len(persons) > 0, nil
}

func isValidStatus(s string) bool {
	switch s {
	case StatusActive, StatusInactive, StatusPending:
		return true
	default:
		return false
	}
}
