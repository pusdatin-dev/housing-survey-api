package models

import (
	"time"
)

type AuditLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RequestID *string   `gorm:"index" json:"request_id"`
	UserID    *string   `gorm:"type:varchar(100);index" json:"user_id"`
	Email     *string   `gorm:"type:varchar(100);index" json:"email"`
	Role      *string   `gorm:"type:varchar(50);index" json:"role"`
	IP        *string   `gorm:"type:varchar(50)" json:"ip"`
	Action    *string   `gorm:"type:varchar(100);index" json:"action"`
	Entity    *string   `gorm:"type:text;index" json:"entity"`
	Detail    *string   `gorm:"type:text" json:"detail"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
