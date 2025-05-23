package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/linux-web-panel/backend/internal/handlers"
	"github.com/linux-web-panel/backend/internal/middleware"
	"github.com/linux-web-panel/backend/pkg/database"
)

func main() {
	// 환경 변수 로드
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// MongoDB 연결
	db, err := database.ConnectMongoDB(os.Getenv("MONGODB_URI"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Gin 라우터 초기화
	r := gin.Default()

	// CORS 미들웨어 설정
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Refresh-Token")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 핸들러 초기화
	authHandler := handlers.NewAuthHandler(db)
	shellHandler := handlers.NewShellHandler()

	// API 라우트
	api := r.Group("/api")
	{
		// 인증 라우트
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 보호된 라우트
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// 시스템 정보
			protected.GET("/system/info", handlers.GetSystemInfo)

			// 유저 관리
			protected.GET("/users", authHandler.ListUsers)
			protected.DELETE("/users/:id", authHandler.DeleteUser)
			protected.PATCH("/users/:id/password", authHandler.ChangePassword)

			// 시스템 모니터링
			protected.GET("/monitor/resources", handlers.GetResourceUsage)
			protected.GET("/monitor/logs", handlers.GetRecentLogs)

			// 쉘 접근
			protected.GET("/ws/shell", shellHandler.HandleWebSocket)

			// 데몬 관리
			protected.GET("/daemons", handlers.ListDaemons)
			protected.GET("/daemons/:name", handlers.GetDaemonStatus)
			protected.POST("/daemons/:name/start", handlers.StartDaemon)
			protected.POST("/daemons/:name/stop", handlers.StopDaemon)
			protected.POST("/daemons/:name/restart", handlers.RestartDaemon)
			protected.GET("/daemons/:name/logs", handlers.DaemonLogs)

			// 파일 시스템 관리
			protected.GET("/files", handlers.ListFiles)
			protected.POST("/files", handlers.CreateFileOrDir)
			protected.DELETE("/files", handlers.DeleteFileOrDir)
			protected.PATCH("/files/rename", handlers.RenameFileOrDir)
			protected.GET("/files/read", handlers.ReadFile)
			protected.PUT("/files/write", handlers.WriteFile)
			protected.POST("/files/upload", handlers.UploadFile)
			protected.GET("/files/download", handlers.DownloadFile)
		}

		// 관리자 전용 라우트
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
		{
			// 관리자 전용 API 라우트 추가 예정
		}
	}

	// 기본 라우트
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 서버 포트 설정
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 서버 시작
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
