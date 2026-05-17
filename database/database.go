package database

import (
	"log"

	"github.com/glebarez/sqlite"
	"github.com/leonardo/netwatch/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	var err error
	// Habilita o modo WAL (Write-Ahead Logging) e timeout para concorrência nativa
	DB, err = gorm.Open(sqlite.Open("netwatch.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("❌ Falha crítica ao inicializar o banco SQLite puro: %v", err)
	}

	// 👇 INSTRUÇÃO CRÍTICA: Limita o pool a 1 conexão para extinguir erros "database is locked"
	sqlDB, err := DB.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(1)
	}

	log.Println("🗄️ Conexão com banco SQLite nativo (netwatch.db) estabelecida com sucesso.")
}

func RunMigrations() {
	if err := DB.AutoMigrate(&domain.User{}, &domain.Group{}, &domain.Device{}); err != nil {
		log.Fatalf("❌ Erro ao aplicar AutoMigrate no SQLite: %v", err)
	}
	log.Println("✅ Estruturas do banco validadas e prontas para uso.")
}