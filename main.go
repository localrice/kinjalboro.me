package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"kinjalboro.me/internal/discord"
	"kinjalboro.me/internal/markdown"
)

var (
	discordStatus string
	statusMu      sync.RWMutex
)

var posts []markdown.Post

func main() {
	// initial fetch
	discordStatus = discord.GetOnlineStatus()

	posts, _ = markdown.LoadPosts()

	// background updater
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			<-ticker.C

			newStatus := discord.GetOnlineStatus()

			statusMu.Lock()
			discordStatus = newStatus
			statusMu.Unlock()
		}
	}()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	// home page
	router.GET("/", func(ctx *gin.Context) {
		statusMu.RLock()
		status := discordStatus
		statusMu.RUnlock()

		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
			"onlineStatus": status,
		})
	})

	router.GET("/posts", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "posts.tmpl", gin.H{
			"posts": posts,
		})
	})

	router.GET("/posts/:slug", func(ctx *gin.Context) {
		slug := ctx.Param("slug")

		for _, post := range posts {
			if post.Slug == slug {
				ctx.HTML(http.StatusOK, "post.tmpl", gin.H{
					"post": post,
				})
				return
			}
		}
		ctx.Status(http.StatusNotFound)
	})

	// test
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Static("/static", "./static")
	router.Run()
}
