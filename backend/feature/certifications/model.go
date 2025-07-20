// features/certifications/model.go
package certifications

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Certification represents the certification entity in the system.
type Certification struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	PersonID            uuid.UUID      `gorm:"type:uuid;not null" json:"person_id"`
	CertificationTypeID uuid.UUID      `gorm:"type:uuid;not null" json:"certification_type_id"`
	CertificationName   string         `gorm:"type:varchar(255);not null" json:"certification_name"`
	IssueDate           time.Time      `gorm:"not null" json:"issue_date"`
	ExpirationDate      time.Time      `gorm:"not null" json:"expiration_date"`
	Description         string         `gorm:"type:text" json:"description"`
	AttachmentURL       string         `gorm:"type:varchar(500)" json:"attachment_url"`
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Certification) TableName() string {
	return "certifications"
}

func (c *Certification) IsValidDates() error {
	if c.ExpirationDate.Before(c.IssueDate) {
		return errors.New("expiration date cannot be before issue date")
	}
	return nil
}

func (c *Certification) IsExpired() bool {
	return time.Now().After(c.ExpirationDate)
}
