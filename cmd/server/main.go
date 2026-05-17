package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/leonardo/netwatch/database"
	"github.com/leonardo/netwatch/internal/domain"
	"github.com/leonardo/netwatch/internal/handlers"
	"github.com/leonardo/netwatch/internal/monitoring"
	"github.com/leonardo/netwatch/internal/repositories"
	"github.com/leonardo/netwatch/internal/services"
	"github.com/leonardo/netwatch/internal/websocket"
)

// setAutoStart registra o .exe no Registro do Windows para iniciar com o sistema
func setAutoStart(exePath string) {
	// O comando 'reg add' insere o programa na aba "Inicializar" do Gerenciador de Tarefas
	cmd := exec.Command("reg", "add", "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run", "/v", "NetWatchPro", "/t", "REG_SZ", "/d", exePath, "/f")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // Executa de forma 100% invisível
	cmd.Run()
}

// openBrowser aguarda a inicialização do servidor e abre a URL no navegador padrão do Windows
func openBrowser(url string) {
	time.Sleep(1 * time.Second)
	cmd := exec.Command("cmd", "/c", "start", url)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Start()
}

func main() {
	// 1. GARANTIA DE DIRETÓRIO: Força o app a rodar na mesma pasta onde o .exe está salvo
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		os.Chdir(exeDir) // Trava o diretório para encontrar a pasta "frontend" e o "netwatch.db"
		
		// 2. AUTO-INICIALIZAÇÃO: Registra o executável atual para iniciar com o Windows
		setAutoStart(exePath)
	}

	fmt.Println("🚀 Iniciando Motor Autossuficiente NetWatch Pro SaaS...")

	database.Connect()
	database.RunMigrations()

	userRepo := repositories.NewUserRepository()
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService, userRepo)

	groupRepo := repositories.NewGroupRepository()
	groupHandler := handlers.NewGroupHandler(groupRepo)

	deviceRepo := repositories.NewDeviceRepository()
	deviceHandler := handlers.NewDeviceHandler(deviceRepo)

	discoveryHandler := handlers.NewDiscoveryHandler(deviceRepo)

	wsHub := websocket.NewHub()
	go wsHub.Run()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "online"}`))
	})

	mux.HandleFunc("GET /api/check-setup", func(w http.ResponseWriter, r *http.Request) {
		var count int64
		database.DB.Model(&domain.User{}).Count(&count)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"setup_needed": count == 0})
	})

	mux.HandleFunc("POST /api/shutdown", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"success", "message":"Servidor encerrado com segurança"}`))
		go func() {
			time.Sleep(1 * time.Second)
			fmt.Println("🛑 Sinal de controle ativado. Finalizando a aplicação...")
			os.Exit(0)
		}()
	})

	mux.HandleFunc("POST /api/register", authHandler.Register)
	mux.HandleFunc("POST /api/login", authHandler.Login)
	mux.HandleFunc("GET /api/users", authHandler.ListUsers)
	mux.HandleFunc("DELETE /api/users/{id}", authHandler.DeleteUser)

	mux.HandleFunc("PUT /api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var payload struct {
			Company        string `json:"company"`
			DiscordWebhook string `json:"discord_webhook"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Requisição incorreta", http.StatusBadRequest)
			return
		}
		database.DB.Model(&domain.User{}).Where("id = ?", id).Updates(map[string]interface{}{
			"company":         payload.Company,
			"discord_webhook": payload.DiscordWebhook,
		})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	})

	mux.HandleFunc("POST /api/groups", groupHandler.Create)
	mux.HandleFunc("GET /api/groups", groupHandler.List)
	mux.HandleFunc("DELETE /api/groups/{id}", groupHandler.Delete)

	mux.HandleFunc("POST /api/devices", deviceHandler.Create)
	mux.HandleFunc("GET /api/devices", deviceHandler.List)
	mux.HandleFunc("PUT /api/devices/{id}", deviceHandler.Update)
	mux.HandleFunc("DELETE /api/devices/{id}", deviceHandler.Delete)

	mux.HandleFunc("POST /api/discovery/scan", discoveryHandler.ScanNetwork)
	mux.HandleFunc("POST /api/discovery/import", discoveryHandler.ImportHosts)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(wsHub, w, r)
	})

	// Serve a interface web
	mux.Handle("/", http.FileServer(http.Dir("./frontend")))

	go func() {
		networkScanner := monitoring.NewScanner(50, func(msg []byte) {
			wsHub.Broadcast <- msg
		})

		for {
			var admin domain.User
			database.DB.Order("created_at asc").First(&admin)

			var devices []domain.Device
			if err := database.DB.Preload("Group").Find(&devices).Error; err != nil || len(devices) == 0 {
				time.Sleep(10 * time.Second)
				continue
			}

			scanResults := networkScanner.Scan(devices, admin.DiscordWebhook)

			for ip, res := range scanResults {
				deviceRepo.RecordCheck(ip, res.IsOnline, res.Status)
			}
			time.Sleep(15 * time.Second)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{Addr: ":" + port, Handler: mux}
	serverErrors := make(chan error, 1)

	// Dispara a abertura automática do navegador
	go openBrowser("http://localhost:" + port)

	go func() {
		log.Printf("📡 Servidor corporativo escutando na porta %s", srv.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("❌ Erro no servidor: %v", err)
	case sig := <-shutdown:
		log.Printf("🛑 Ordem de desligamento segura disparada (%v).", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("❌ Erro ao fechar as conexões: %v", err)
			err = srv.Close()
		}
		log.Println("✅ Acesso encerrado e serviços parados com sucesso.")
	}
}