package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tictactoe/internal/game"
	"tictactoe/internal/map_builder"
	"tictactoe/internal/map_reader"
	"tictactoe/internal/map_storage"
)

type Server struct {
	mb *map_builder.MapBuilder
	mr *map_reader.MapReader
	r  *gin.Engine
}

func NewServer() *Server {
	s := &Server{
		r:  gin.Default(),
		mb: map_builder.NewMapBuilder(),
		mr: map_reader.NewMapReader(),
	}

	s.r.Static("/static", "./static")

	s.r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s.r.GET("/api/maps/status", func(c *gin.Context) {
		gameStr := c.Query("game")

		g, err := game.FromString(gameStr)
		if err != nil {
			fmt.Println("error parsing game str", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		progress, _ := map_storage.GetProgress(g)

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data": gin.H{
				"progress": progress,
			},
		})
	})

	s.r.POST("/api/maps/build", func(c *gin.Context) {
		gameStr := c.Query("game")

		g, err := game.FromString(gameStr)
		if err != nil {
			fmt.Println("error parsing game str", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, started := map_storage.GetProgress(g)
		if started {
			c.JSON(http.StatusBadRequest, gin.H{"error": "map is already being built"})
			return
		}

		if err := s.mb.BuildWinMap(g); err != nil {
			fmt.Println("error building win map", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s.r.GET("/api/chances", func(c *gin.Context) {
		gameStr := c.Query("game")

		g, err := game.FromString(gameStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		progress, started := map_storage.GetProgress(g)
		if !started || progress != 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "map is not ready"})
			return
		}

		stats, err := s.mr.GetGameStats(g)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data":   stats,
		})
	})

	s.r.GET("/api/next-move", func(c *gin.Context) {
		gameStr := c.Query("game")

		g, err := game.FromString(gameStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		x, y, err := s.mr.GetNextMove(g)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data":   gin.H{"x": x, "y": y},
		})
	})

	return s
}

func (s *Server) Start(port int) {
	if err := s.r.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}
