package domain

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID             uint           `json:"id" gorm:"primarykey"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	Name           string         `json:"name"`
	Email          string         `json:"email" gorm:"unique"`
	Password       string         `json:"-"`
	Company        string         `json:"company"`
	DiscordWebhook string         `json:"discord_webhook"`
}