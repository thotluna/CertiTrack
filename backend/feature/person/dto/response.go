package dto

import "time"

type PersonResponse struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewPersonResponse(id, fullName, email, phone, status string, createdAt, updatedAt time.Time) *PersonResponse {
	return &PersonResponse{
		ID:        id,
		FullName:  fullName,
		Email:     email,
		Phone:     phone,
		Status:    status,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
