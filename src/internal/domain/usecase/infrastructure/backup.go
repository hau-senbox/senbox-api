package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sen-global-api/internal/domain/response"
	"sen-global-api/pkg/monitor"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Helper function to check if an IP address is localhost
func isLocalhost(ip string) bool {
	localhostIPs := []string{"127.0.0.1", "::1"}
	for _, localhost := range localhostIPs {
		if ip == localhost {
			return true
		}
	}

	// Check for IPs in localhost subnet (IPv4 only)
	if strings.HasPrefix(ip, "127.") {
		return true
	}

	return false
}

func BackupDatabase() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		isLocal := isLocalhost(clientIP)

		if !isLocal {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Access forbidden: only localhost is allowed",
			})
			return
		}

		defer func() {
			Backup()
		}()

		c.JSON(200, response.SucceedResponse{
			Code:    200,
			Message: "Scheduled for backup",
		})
	}
}

func Backup() {
	//Remove old backup (SQL)
	cmd := exec.Command("rm", "sen_master_db.sql")
	err := cmd.Run()

	if err != nil {
		monitor.SendMessageViaTelegram(
			"[URGENT] Error when backup database: " + err.Error(),
		)
		return
	}

	// Remove old backup (tar)
	cmd = exec.Command("rm", "sen_master_db.tar.gz")
	err = cmd.Run()

	if err != nil {
		monitor.SendMessageViaTelegram(
			"[URGENT] Error when backup database: " + err.Error(),
		)
		return
	}

	// Backup database
	cmd = exec.Command("bash", "-c", "mysqldump -u sen_master sen_master_db > sen_master_db.sql")

	// Run the command
	err = cmd.Run()

	if err != nil {
		monitor.SendMessageViaTelegram(
			"[URGENT] Error when backup database: " + err.Error(),
		)
		return
	}

	// Tar backup
	cmd = exec.Command("bash", "-c", "tar -czvf sen_master_db.tar.gz sen_master_db.sql")
	err = cmd.Run()

	if err != nil {
		monitor.SendMessageViaTelegram(
			"[URGENT] Error when backup database: " + err.Error(),
		)
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		monitor.SendMessageViaTelegram(fmt.Sprintf("Error getting current directory: %s", err))
		return
	}
	srv, err := drive.NewService(context.Background(), option.WithCredentialsFile(pwd+"/credentials/google_service_account.json"))

	if err != nil {
		monitor.SendMessageViaTelegram(
			"[URGENT] Error when establishing Google Drive service for uploading backup: " + err.Error(),
		)
		return
	}

	file, err := os.Open("sen_master_db.tar.gz")
	if err != nil {
		log.Errorf("Error: %v", err)
		monitor.SendMessageViaTelegram(
			"[URGENT] Error when backup database: " + err.Error(),
		)
		return
	}
	defer file.Close()

	now := time.Now()
	nowInString := now.Format("2006-01-02 15:04:05")
	driveID := os.Getenv("SENBOX_BACKUP_GOOGLE_DRIVE_ID")

	if driveID == "" {
		log.Error("Error: SENBOX_BACKUP_GOOGLE_DRIVE_ID is empty")
		monitor.SendMessageViaTelegram(
			"[URGENT] Error when backup database: SENBOX_BACKUP_GOOGLE_DRIVE_ID is empty",
		)
		return
	}

	f := &drive.File{
		Name:    nowInString + "_sen_master_db.tar.gz",
		Parents: []string{driveID},
	}
	_, err = srv.Files.Create(f).Media(file).Do()
	if err != nil {
		log.Errorf("Error: %v", err)
		monitor.SendMessageViaTelegram(
			"[URGENT] Error when upload database: " + err.Error(),
		)
		return
	}

	log.Info("Backup database success")

	monitor.SendMessageViaTelegram(
		"Backup database success",
	)
}
