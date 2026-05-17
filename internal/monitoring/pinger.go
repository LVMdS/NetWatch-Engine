package monitoring

import (
	"time"

	"github.com/go-ping/ping"
)

// PingResult guarda o resultado da sondagem de um dispositivo
type PingResult struct {
	IP      string
	IsUp    bool
	Latency time.Duration
}

// PingHost tenta enviar um pacote ICMP para o IP especificado
func PingHost(ip string) PingResult {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		return PingResult{IP: ip, IsUp: false, Latency: 0}
	}

	// Configuração obrigatória para envio de pacotes ICMP no Windows
	pinger.SetPrivileged(true) 
	
	// Regras do Ping
	pinger.Count = 1                  // Apenas 1 pacote para ser rápido
	pinger.Timeout = time.Second * 2  // Desiste após 2 segundos
	
	err = pinger.Run() // Esta função é bloqueante, ela espera o ping terminar
	if err != nil {
		return PingResult{IP: ip, IsUp: false, Latency: 0}
	}

	stats := pinger.Statistics()
	
	// Se recebeu pelo menos 1 pacote de volta, está online
	if stats.PacketsRecv > 0 {
		return PingResult{IP: ip, IsUp: true, Latency: stats.AvgRtt}
	}

	return PingResult{IP: ip, IsUp: false, Latency: 0}
}