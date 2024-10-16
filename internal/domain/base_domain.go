package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseDomain struct {
	ID        uint       `gorm:"primary_key;autoIncrement" json:"id"`
	UUID      uuid.UUID  `gorm:"type:char(36)" json:"uuid"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// Hook BeforeCreate untuk menggenerate UUID sebelum entri data
func (b *BaseDomain) BeforeCreate(tx *gorm.DB) (err error) {
	// Jika UUID belum di-set, generate UUID baru
	if b.UUID == uuid.Nil {
		b.UUID = uuid.New()
	}
	return
}
