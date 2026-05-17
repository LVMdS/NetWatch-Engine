package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TotalPings = promauto.NewCounter(prometheus.CounterOpts{
		Name: "netwatch_pings_total",
		Help: "Número total de sondagens ICMP executadas",
	})

	TotalFailures = promauto.NewCounter(prometheus.CounterOpts{
		Name: "netwatch_device_failures_total",
		Help: "Número total de quedas de dispositivos detectadas",
	})

	CurrentOnlineDevices = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "netwatch_online_devices_current",
		Help: "Número atual de dispositivos respondendo ao ping",
	})

	// 👇 NOVA MÉTRICA 1: Totalizador de dispositivos offline
	CurrentOfflineDevices = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "netwatch_offline_devices_current",
		Help: "Número atual de dispositivos que não estão respondendo",
	})

	// 👇 NOVA MÉTRICA 2: Rastreio de status individual por IP (1 = Online, 0 = Offline)
	DeviceStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "netwatch_device_status",
		Help: "Status em tempo real de cada IP monitorado",
	}, []string{"ip"})
)