package monitoring

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	"github.com/leonardo/netwatch/internal/domain"
)

type CheckResult struct {
	Status   string
	IsOnline bool
}

type Scanner struct {
	State     map[string]string
	Failures  map[string]int
	broadcast func([]byte)
	client    *http.Client
}

func NewScanner(maxConcurrent int, wsBroadcast func([]byte)) *Scanner {
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Scanner{
		State:     make(map[string]string),
		Failures:  make(map[string]int),
		broadcast: wsBroadcast,
		client:    &http.Client{Timeout: 10 * time.Second, Transport: customTransport},
	}
}

func (s *Scanner) Scan(devices []domain.Device, webhookURL string) map[string]CheckResult {
	results := make(map[string]CheckResult)

	for _, dev := range devices {
		start := time.Now()
		var isOnline bool

		if dev.Port != "" {
			target := net.JoinHostPort(dev.IPAddress, dev.Port)
			conn, err := net.DialTimeout("tcp", target, 2*time.Second)
			if err == nil {
				conn.Close()
				isOnline = true
			}
		} else {
			isOnline = pingTarget(dev.IPAddress)
		}

		elapsed := time.Since(start)
		latencyStr := fmt.Sprintf("%dms", elapsed.Milliseconds())
		if latencyStr == "0ms" {
			latencyStr = "< 15ms"
		}

		statusStr := "ONLINE"
		if isOnline {
			s.Failures[dev.IPAddress] = 0
		} else {
			s.Failures[dev.IPAddress]++
			if s.Failures[dev.IPAddress] >= 3 {
				statusStr = "OFFLINE"
			} else {
				if prev, exists := s.State[dev.IPAddress]; exists {
					statusStr = prev
				}
			}
		}

		prevState, hasHistory := s.State[dev.IPAddress]
		if !hasHistory || prevState != statusStr {
			if hasHistory {
				log.Printf("🔄 Transição de Estado Validada [%s]: %s -> %s", dev.Name, prevState, statusStr)
				go s.sendDiscordAlert(webhookURL, dev, statusStr)
			}
		}

		s.State[dev.IPAddress] = statusStr
		results[dev.IPAddress] = CheckResult{Status: statusStr, IsOnline: isOnline}

		wsMsg, _ := json.Marshal(map[string]string{
			"ip":      dev.IPAddress,
			"status":  statusStr,
			"latency": latencyStr,
		})
		s.broadcast(wsMsg)
	}

	return results
}

func pingTarget(ip string) bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "1", ip)
	}
	
	// 👇 INSTRUÇÃO CRÍTICA: Bloqueia a criação de interface visual no Windows
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	
	return cmd.Run() == nil
}

func (s *Scanner) sendDiscordAlert(webhookURL string, dev domain.Device, status string) {
	if webhookURL == "" {
		return
	}

	var color int
	var statusIcon string

	if status == "ONLINE" {
		color = 3066993
		statusIcon = "🟢"
	} else {
		color = 15158332
		statusIcon = "🔴"
	}

	groupName := "Sem Grupo"
	if dev.Group.Name != "" {
		groupName = dev.Group.Name
	}

	serviceType := "ICMP Echo"
	if dev.Port != "" {
		serviceType = "TCP Port " + dev.Port
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       statusIcon + " Alerta de Conectividade",
				"description": fmt.Sprintf("O dispositivo **%s** alterou o seu estado operacional.", dev.Name),
				"color":       color,
				"fields": []map[string]interface{}{
					{"name": "Estado", "value": status, "inline": true},
					{"name": "Alvo", "value": fmt.Sprintf("`%s`", dev.IPAddress), "inline": true},
					{"name": "Protocolo", "value": serviceType, "inline": true},
					{"name": "Segmento", "value": groupName, "inline": false},
				},
				"footer":    map[string]interface{}{"text": "NetWatch Pro - Enterprise Engine"},
				"timestamp": time.Now().Format(time.RFC3339),
			},
		},
	}

	jsonPayload, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	s.client.Do(req)
}