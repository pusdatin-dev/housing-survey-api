package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`       // Nullable
	Email     *string    `gorm:"type:varchar(100);index" json:"email"` // Nullable, for authenticated users
	IP        *string    `gorm:"type:varchar(50)" json:"ip"`           // For non-auth user tracking
	Action    *string    `gorm:"type:varchar(100)" json:"action"`      // Nullable
	Entity    *string    `gorm:"type:text;index" json:"entity"`        // API path
	Detail    *string    `gorm:"type:text" json:"detail"`              // Optional info
	CreatedAt time.Time  `gorm:"index" json:"created_at"`
}
