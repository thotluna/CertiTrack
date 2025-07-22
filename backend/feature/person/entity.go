package person

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusPending  = "pending"
)

type Person struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FirstName string         `gorm:"size:100;not null" json:"first_name"`
	LastName  string         `gorm:"size:100;not null" json:"last_name"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone,omitempty"`
	Status    string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Person) TableName() string {
	return "persons"
}

func NewPerson(firstName, lastName, email, phone string) *Person {
	return &Person{
		ID:        uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Status:    StatusActive,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (p *Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

func (p *Person) Update(firstName, lastName, phone string) {
	p.FirstName = firstName
	p.LastName = lastName
	p.Phone = phone
	p.UpdatedAt = time.Now().UTC()
}

func (p *Person) ChangeStatus(status string) {
	p.Status = status
	p.UpdatedAt = time.Now().UTC()
}
