package dto

import (
	"errors"
	"strings"
)

type CreatePersonRequest struct {
	FirstName string `json:"first_name" binding:"required,min=2,max=100"`
	LastName  string `json:"last_name" binding:"required,min=2,max=100"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone,omitempty" binding:"omitempty,e164"`
}

func (r *CreatePersonRequest) Validate() error {
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Email = strings.TrimSpace(r.Email)
	r.Phone = strings.TrimSpace(r.Phone)

	if r.FirstName == "" || len(r.FirstName) < 2 {
		return errors.New("first name is required and must be at least 2 characters")
	}

	if r.LastName == "" || len(r.LastName) < 2 {
		return errors.New("last name is required and must be at least 2 characters")
	}

	if r.Email == "" {
		return errors.New("email is required")
	}

	return nil
}

type UpdatePersonRequest struct {
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,min=2,max=100"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty,e164"`
}

func (r *UpdatePersonRequest) Validate() error {
	if r.FirstName == nil && r.LastName == nil && r.Email == nil && r.Phone == nil {
		return errors.New("at least one field must be provided for update")
	}

	if r.FirstName != nil {
		trimmed := strings.TrimSpace(*r.FirstName)
		if len(trimmed) < 2 {
			return errors.New("first name must be at least 2 characters")
		}
		r.FirstName = &trimmed
	}

	if r.LastName != nil {
		trimmed := strings.TrimSpace(*r.LastName)
		if len(trimmed) < 2 {
			return errors.New("last name must be at least 2 characters")
		}
		r.LastName = &trimmed
	}

	if r.Email != nil {
		email := strings.TrimSpace(*r.Email)
		if email == "" {
			return errors.New("email cannot be empty")
		}
		r.Email = &email
	}

	if r.Phone != nil {
		phone := strings.TrimSpace(*r.Phone)
		r.Phone = &phone
	}

	return nil
}
