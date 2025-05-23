package handlers

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SystemInfo struct {
	Hostname string `json:"hostname"`
	Uptime   string `json:"uptime"`
	OS       string `json:"os"`
	Kernel   string `json:"kernel"`
}

type ResourceUsage struct {
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
	Disk   float64 `json:"disk"`
}

type LogEntry struct {
	Line string `json:"line"`
}

func GetSystemInfo(c *gin.Context) {
	info := SystemInfo{}

	// 호스트네임 가져오기
	if hostname, err := exec.Command("hostname").Output(); err == nil {
		info.Hostname = strings.TrimSpace(string(hostname))
	}

	// 업타임 가져오기
	if uptime, err := exec.Command("uptime", "-p").Output(); err == nil {
		info.Uptime = strings.TrimSpace(string(uptime))
	}

	// OS 정보 가져오기
	if osInfo, err := exec.Command("cat", "/etc/os-release").Output(); err == nil {
		lines := strings.Split(string(osInfo), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				info.OS = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
				break
			}
		}
	}

	// 커널 버전 가져오기
	if kernel, err := exec.Command("uname", "-r").Output(); err == nil {
		info.Kernel = strings.TrimSpace(string(kernel))
	}

	c.JSON(http.StatusOK, info)
}

func GetResourceUsage(c *gin.Context) {
	// CPU 사용률
	cpuUsage := 0.0
	if out, err := exec.Command("sh", "-c", "top -bn1 | grep 'Cpu(s)' | awk '{print $2 + $4}'").Output(); err == nil {
		cpuUsage, _ = strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	}

	// 메모리 사용률
	memUsage := 0.0
	if out, err := exec.Command("sh", "-c", "free | grep Mem | awk '{print $3/$2 * 100.0}'").Output(); err == nil {
		memUsage, _ = strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	}

	// 디스크 사용률
	diskUsage := 0.0
	if out, err := exec.Command("sh", "-c", "df --total | grep total | awk '{print $5}' | tr -d '%' ").Output(); err == nil {
		diskUsage, _ = strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	}

	c.JSON(http.StatusOK, ResourceUsage{
		CPU:    cpuUsage,
		Memory: memUsage,
		Disk:   diskUsage,
	})
}

func GetRecentLogs(c *gin.Context) {
	logPath := "/var/log/syslog"
	if _, err := exec.LookPath("journalctl"); err == nil {
		// systemd 기반이면 journalctl 사용
		out, err := exec.Command("journalctl", "-n", "50", "--no-pager").Output()
		if err == nil {
			lines := strings.Split(string(out), "\n")
			entries := []LogEntry{}
			for _, line := range lines {
				if line != "" {
					entries = append(entries, LogEntry{Line: line})
				}
			}
			c.JSON(http.StatusOK, entries)
			return
		}
	}
	// 일반 syslog 파일 사용
	data, err := ioutil.ReadFile(logPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "로그 파일을 읽을 수 없습니다."})
		return
	}
	lines := strings.Split(string(data), "\n")
	entries := []LogEntry{}
	for i := len(lines) - 51; i < len(lines)-1; i++ {
		if i >= 0 {
			entries = append(entries, LogEntry{Line: lines[i]})
		}
	}
	c.JSON(http.StatusOK, entries)
}
