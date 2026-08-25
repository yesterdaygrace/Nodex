package main

import (
	"crudin/controllers"
	"crudin/models"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Configure structured logging (JSON on most systems) at program start
	// so app-level lifecycle events are emitted as structured records rather
	// than ad-hoc lines. Gin's own access logging via gin.Default() is left
	// untouched, so there is no duplicate request-logging middleware here.
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	//inisialiasai Gin
	router := gin.Default()

	//panggil koneksi database
	models.ConnectDatabase()

	//membuat route dengan method GET
	router.GET("/", func(c *gin.Context) {

		//return response JSON
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	//membuat route get all posts
	router.GET("/api/posts", controllers.FindPosts)
	router.POST("/api/posts", controllers.StorePost)
	// /trashed must be registered before /:id so the static segment is matched
	// first instead of being captured as the :id param.
	router.GET("/api/posts/trashed", controllers.FindTrashedPosts)
	router.GET("/api/posts/:id", controllers.FindPostById)
	router.PUT("/api/posts/:id", controllers.UpdatePost)
	router.DELETE("/api/posts/:id", controllers.DeletePost)
	router.PUT("/api/posts/:id/restore", controllers.RestorePost)
	router.DELETE("/api/posts/:id/permanent", controllers.DeletePermanentPost)
	router.PUT("/api/posts/:id/archive", controllers.ArchivePost)
	router.PUT("/api/posts/:id/unarchive", controllers.UnarchivePost)

	//mulai server dengan port 3001
	router.Run(":3001")
}
