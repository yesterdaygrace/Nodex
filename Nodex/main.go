package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"Nodex/controllers"
	"Nodex/models"

	"github.com/gin-gonic/gin"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	// Production mode when APP_ENV=production
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS middleware if configured (DEPLOY.md §29)
	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		allowed := parseCSV(origins)
		router.Use(corsMiddleware(allowed))
	}

	// Security headers (DEPLOY.md §56)
	router.Use(securityHeadersMiddleware())

	// Health / readiness - must be unauthenticated and lightweight
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/ready", func(c *gin.Context) {
		// Verify DB connectivity
		if models.DB == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "reason": "db not initialized"})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		sqlDB, err := models.DB.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "reason": "db handle"})
			return
		}
		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "reason": "db ping failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// DB connect with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := models.ConnectDatabase(ctx); err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	// Optional: configure pool per DEPLOY.md §53
	if sqlDB, err := models.DB.DB(); err == nil {
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
	}

	// API routes - keep Hello World at /api for health UX, root serves frontend if present
	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello World!"})
	})
	router.GET("/api/posts", controllers.FindPosts)
	router.POST("/api/posts", controllers.StorePost)
	router.GET("/api/posts/trashed", controllers.FindTrashedPosts)
	router.GET("/api/posts/:id", controllers.FindPostById)
	router.PUT("/api/posts/:id", controllers.UpdatePost)
	router.DELETE("/api/posts/:id", controllers.DeletePost)
	router.PUT("/api/posts/:id/restore", controllers.RestorePost)
	router.DELETE("/api/posts/:id/permanent", controllers.DeletePermanentPost)
	router.PUT("/api/posts/:id/archive", controllers.ArchivePost)
	router.PUT("/api/posts/:id/unarchive", controllers.UnarchivePost)

	// Serve frontend static if built (single-domain deploys per DEPLOY.md §50)
	frontendDist := filepath.Join("frontend", "dist")
	hasFrontend := false
	if _, err := os.Stat(frontendDist); err == nil {
		if _, err2 := os.Stat(filepath.Join(frontendDist, "index.html")); err2 == nil {
			hasFrontend = true
		}
	}
	if hasFrontend {
		slog.Info("serving frontend", "dir", frontendDist)
		router.Static("/assets", filepath.Join(frontendDist, "assets"))
		router.StaticFile("/vite.svg", filepath.Join(frontendDist, "vite.svg"))
		// Root serves SPA index.html, /api keeps Hello World JSON
		router.GET("/", func(c *gin.Context) {
			c.File(filepath.Join(frontendDist, "index.html"))
		})
		// SPA fallback: serve index.html for non-API, non-health routes
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			// API / health / ready should return JSON 404, not SPA
			if strings.HasPrefix(path, "/api/") || path == "/health" || path == "/ready" || strings.HasPrefix(path, "/assets/") {
				c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "not found", "data": nil})
				return
			}
			c.File(filepath.Join(frontendDist, "index.html"))
		})
	} else {
		// No frontend build - keep simple Hello World at root
		router.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello World!"})
		})
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	// Gin already listens on 0.0.0.0 when using :port
	addr := ":" + port

	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		slog.Info("starting server", "addr", addr, "env", os.Getenv("APP_ENV"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM (DEPLOY.md §24)
	ctxSig, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctxSig.Done()
	slog.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}
	// Close DB
	if models.DB != nil {
		if sqlDB, err := models.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
}

// parseCSV splits comma-separated string and trims spaces
func parseCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	allowedSet := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowedSet[o] = struct{}{}
	}
	allowAll := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}
		allowed := allowAll
		if !allowed {
			_, allowed = allowedSet[origin]
		}
		// Handle wildcard disallowed for auth APIs per spec: we still
		// respond with specific origin if matched, never "*"
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// CSP minimal - allow self, inline styles needed for Vue
		// Do not break app with overly strict policy
		c.Next()
	}
}
