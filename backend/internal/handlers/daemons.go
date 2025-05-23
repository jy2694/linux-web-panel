package handlers

import (
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

type DaemonInfo struct {
	Name   string `json:"name"`
	Load   string `json:"load"`
	Active string `json:"active"`
	Sub    string `json:"sub"`
	Desc   string `json:"desc"`
}

// 서비스 목록 조회
func ListDaemons(c *gin.Context) {
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--no-legend")
	out, err := cmd.Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "서비스 목록을 가져올 수 없어요."})
		return
	}
	lines := strings.Split(string(out), "\n")
	var daemons []DaemonInfo
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 5 {
			daemons = append(daemons, DaemonInfo{
				Name:   fields[0],
				Load:   fields[1],
				Active: fields[2],
				Sub:    fields[3],
				Desc:   strings.Join(fields[4:], " "),
			})
		}
	}
	c.JSON(http.StatusOK, daemons)
}

// 서비스 상태 조회
func GetDaemonStatus(c *gin.Context) {
	name := c.Param("name")
	cmd := exec.Command("systemctl", "status", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": string(out)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": string(out)})
}

// 서비스 시작
func StartDaemon(c *gin.Context) {
	name := c.Param("name")
	cmd := exec.Command("systemctl", "start", name)
	err := cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "시작 실패"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "시작됨"})
}

// 서비스 중지
func StopDaemon(c *gin.Context) {
	name := c.Param("name")
	cmd := exec.Command("systemctl", "stop", name)
	err := cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "중지 실패"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "중지됨"})
}

// 서비스 재시작
func RestartDaemon(c *gin.Context) {
	name := c.Param("name")
	cmd := exec.Command("systemctl", "restart", name)
	err := cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "재시작 실패"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "재시작됨"})
}

// 서비스 로그
func DaemonLogs(c *gin.Context) {
	name := c.Param("name")
	cmd := exec.Command("journalctl", "-u", name, "-n", "100", "--no-pager")
	out, err := cmd.Output()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "로그를 가져올 수 없어요."})
		return
	}
	lines := strings.Split(string(out), "\n")
	c.JSON(http.StatusOK, lines)
}
