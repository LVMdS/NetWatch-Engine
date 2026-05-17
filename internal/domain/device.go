package domain

import (
	"time"
	"gorm.io/gorm"
)

type Device struct {
	ID               uint           `json:"id" gorm:"primarykey"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
	Name             string         `json:"name"`
	IPAddress        string         `json:"ip_address"`
	Port             string         `json:"port"` // Opcional: Porta TCP alvo (ex: "80", "5432")
	IPRange          string         `json:"ip_range"`
	SubnetMask       string         `json:"subnet_mask"`
	Gateway          string         `json:"gateway"`
	Description      string         `json:"description"`
	Status           string         `json:"status"`
	LastSeen         *time.Time     `json:"last_seen"`
	GroupID          *uint          `json:"group_id"`
	Group            Group          `json:"group" gorm:"foreignKey:GroupID"`

	// Métricas de Observabilidade para Cálculo de SLA / Uptime
	TotalChecks      int64          `json:"total_checks"`
	SuccessfulChecks int64          `json:"successful_checks"`
}