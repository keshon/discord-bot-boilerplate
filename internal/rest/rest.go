package rest

import (
	"io"
	"strconv"

	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
	"github.com/keshon/discord-bot-boilerplate/internal/botsdef"
)

type Rest struct {
	Bots []map[string]map[string]botsdef.Discord
}

// NewRest initializes a new Rest object with the given botInstances.
//
// botInstances: a map of bot instances
// Returns a pointer to the newly initialized Rest object
func NewRest(bots []map[string]map[string]botsdef.Discord) *Rest {
	return &Rest{
		Bots: bots,
	}
}

// Start starts the REST API routes.
//
// router: *gin.Engine
func (r *Rest) Start(router *gin.Engine) {
	slog.Info("REST API routes started")

	router.GET("/", func(ctx *gin.Context) {
		toc := generateTableOfContents(router)
		ctx.JSON(http.StatusOK, gin.H{"api_methods": toc})
	})

	logRoutes := router.Group("/log")
	{
		r.registerLogRoutes(logRoutes)
	}

	guildRoutes := router.Group("/guild")
	{
		r.registerGuildRoutes(guildRoutes)
	}

	avatarRoutes := router.Group("/avatar")
	{
		r.registerAvatarRoutes(avatarRoutes)
	}
}

type GuildInfo struct {
	GuildID string
}

type GuildSession struct {
	GuildID     string
	GuildActive bool
	BotStatus   string
}

// generateTableOfContents generates a table of contents for the given gin router.
//
// router *gin.Engine - The gin router to generate the table of contents for.
// []map[string]string - The generated table of contents.
func generateTableOfContents(router *gin.Engine) []map[string]string {
	var toc []map[string]string

	for _, routeInfo := range router.Routes() {
		route := map[string]string{
			"method": routeInfo.Method,
			"path":   routeInfo.Path,
		}
		toc = append(toc, route)
	}

	return toc
}

// Examples:
// http://localhost:8080/log
// http://localhost:8080/log/download
// http://localhost:8080/log/clear

// registerLogRoutes registers routes to handle logging in the Rest struct.
//
// router: pointer to the gin RouterGroup where the routes will be registered.
// No return value.
func (r *Rest) registerLogRoutes(router *gin.RouterGroup) {
	router.GET("/", func(ctx *gin.Context) {
		file, err := os.Open("./logs/all-levels.log")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		b, err := io.ReadAll(file)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.Data(http.StatusOK, "text/plain", b)
	})

	router.GET("/download", func(ctx *gin.Context) {
		file, err := os.Open("./logs/all-levels.log")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		ctx.File("./logs/all-levels.log")
	})

	router.GET("/clear", func(ctx *gin.Context) {
		logFilePath := "./logs/all-levels.log"

		err := os.Truncate(logFilePath, 0)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = slog.Flush()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, "Log file cleared")
	})
}

// Examples:
// http://localhost:8080/guild/info/897053062030585916
// http://localhost:8080/guild/playing/897053062030585916

// registerGuildRoutes registers the guild routes for the Rest struct.
//
// It takes a pointer to a gin.RouterGroup as a parameter and has no return type.
func (r *Rest) registerGuildRoutes(router *gin.RouterGroup) {
	router.GET("/ids", func(ctx *gin.Context) {
		activeSessions := []GuildInfo{}

		for guildID := range r.Bots {
			guildID := strconv.Itoa(guildID)
			activeSessions = append(activeSessions, GuildInfo{GuildID: guildID})
		}

		ctx.JSON(http.StatusOK, activeSessions)
	})
}

// Examples:
// http://localhost:8080/avatar
// http://localhost:8080/avatar/random

// registerAvatarRoutes registers routes for avatar handling.
//
// router: The gin router group to register the avatar routes.
// None.
func (r *Rest) registerAvatarRoutes(router *gin.RouterGroup) {
	router.GET("/", func(ctx *gin.Context) {
		const folderPath = "./assets/avatars"
		imageList, err := getImageList(folderPath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, imageList)
	})

	router.GET("/random", func(ctx *gin.Context) {
		const folderPath = "./assets/avatars"
		randomImage, err := getRandomImage(folderPath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		imagePath := filepath.Join(folderPath, randomImage)
		ctx.File(imagePath)
	})
}
