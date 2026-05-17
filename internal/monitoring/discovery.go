package monitoring

import (
	"context"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

type DiscoveredHost struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	MAC      string `json:"mac"`
}

func RunDiscovery(targetIPs []string) []DiscoveredHost {
	var results []DiscoveredHost
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 👇 REDUZIDO DE 100 PARA 20: Suaviza o uso de CPU e evita travamentos no Windows
	sem := make(chan struct{}, 20)

	for _, ip := range targetIPs {
		wg.Add(1)
		sem <- struct{}{}

		go func(target string) {
			defer wg.Done()
			defer func() { <-sem }()

			if pingTarget(target) {
				host := DiscoveredHost{
					IP:       target,
					Hostname: resolveHostname(target),
					MAC:      getMACAddress(target),
				}
				mu.Lock()
				results = append(results, host)
				mu.Unlock()
			}
		}(ip)
	}

	wg.Wait()
	return results
}

func resolveHostname(ip string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	names, err := net.DefaultResolver.LookupAddr(ctx, ip)
	if err == nil && len(names) > 0 {
		return strings.TrimSuffix(names[0], ".")
	}
	return "Host Desconhecido"
}

func getMACAddress(ip string) string {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("arp", "-a", ip)
	} else {
		cmd = exec.Command("arp", "-n", ip)
	}

	// 👇 INSTRUÇÃO CRÍTICA: Oculta completamente a janela do console no Windows
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	out, err := cmd.Output()
	if err != nil {
		return "N/A"
	}

	re := regexp.MustCompile(`([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})`)
	match := re.FindString(string(out))
	if match != "" {
		return strings.ToUpper(strings.ReplaceAll(match, "-", ":"))
	}
	return "N/A"
}