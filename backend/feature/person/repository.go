package person

import (
	"context"

	"github.com/google/uuid"
)

type PersonRepository interface {
	Create(ctx context.Context, person *Person) error
	FindByID(ctx context.Context, id uuid.UUID) (*Person, error)
	FindAll(ctx context.Context, filters map[string]interface{}) ([]*Person, error)
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	Update(ctx context.Context, person *Person) error
	Delete(ctx context.Context, id uuid.UUID) error
}
