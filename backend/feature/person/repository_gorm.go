package person

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type personGormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) PersonRepository {
	return &personGormRepository{db: db}
}

func (r *personGormRepository) Create(ctx context.Context, person *Person) error {
	var count int64
	err := r.db.WithContext(ctx).Model(&Person{}).
		Where("email = ?", person.Email).
		Count(&count).Error

	if err != nil {
		return fmt.Errorf("error checking for duplicate email: %w", err)
	}

	if count > 0 {
		return ErrEmailAlreadyExists
	}

	result := r.db.WithContext(ctx).Create(person)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *personGormRepository) FindByID(ctx context.Context, id uuid.UUID) (*Person, error) {
	var person Person
	result := r.db.WithContext(ctx).First(&person, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrPersonNotFound
		}
		return nil, result.Error
	}
	return &person, nil
}

func (r *personGormRepository) FindAll(ctx context.Context, filters map[string]interface{}) ([]*Person, error) {
	var persons []*Person
	query := r.db.WithContext(ctx).Model(&Person{})
	query = applyFilters(query, filters)
	if err := query.Find(&persons).Error; err != nil {
		return nil, fmt.Errorf("error finding persons: %w", err)
	}

	return persons, nil
}

func (r *personGormRepository) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&Person{})
	cleanFilters := make(map[string]interface{})
	for k, v := range filters {
		if k != "limit" && k != "offset" {
			cleanFilters[k] = v
		}
	}
	query = applyFilters(query, cleanFilters)
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error counting persons: %w", err)
	}

	return count, nil
}

func applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, value := range filters {
		switch key {
		case "limit":
			if limit, ok := value.(int); ok {
				query = query.Limit(limit)
			}
		case "offset":
			if offset, ok := value.(int); ok {
				query = query.Offset(offset)
			}
		case "name":
			searchTerm := fmt.Sprintf("%%%s%%", value.(string))
			query = query.Where("first_name ILIKE ? OR last_name ILIKE ?", searchTerm, searchTerm)
		case "status":
			query = query.Where("status = ?", value)
		default:
			query = query.Where(key, value)
		}
	}
	return query
}

func (r *personGormRepository) Update(ctx context.Context, person *Person) error {
	existing, err := r.FindByID(ctx, person.ID)
	if err != nil {
		return err
	}

	if person.Email != existing.Email {
		var count int64
		err := r.db.WithContext(ctx).Model(&Person{}).
			Where("email = ? AND id != ?", person.Email, person.ID).
			Count(&count).Error

		if err != nil {
			return err
		}

		if count > 0 {
			return ErrEmailAlreadyExists
		}
	}

	result := r.db.WithContext(ctx).Save(person)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrEmailAlreadyExists
		}
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrPersonNotFound
	}

	return nil
}

func (r *personGormRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := r.FindByID(ctx, id); err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Delete(&Person{ID: id})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrPersonNotFound
	}

	return nil
}
