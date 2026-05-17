package domain

import (
	"time"
	"gorm.io/gorm"
)

type Group struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	Name       string         `json:"name"`
	IPRange    string         `json:"ip_range"`
	SubnetMask string         `json:"subnet_mask"`
	Gateway    string         `json:"gateway"`
}