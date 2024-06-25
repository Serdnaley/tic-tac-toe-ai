package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tictactoe/game"
	"tictactoe/map_builder"
	"tictactoe/map_reader"
	"tictactoe/util"
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

	s.r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s.r.GET("/maps/:map/status", func(c *gin.Context) {
		mapKey := c.Params.ByName("map")

		w, h, l, err := util.ParseMapKey(mapKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "data": s.mb.GetMapStatus(w, h, l)})
	})

	s.r.POST("/maps/build", func(c *gin.Context) {
		gameStr := c.Query("game")

		g, err := game.FromString(gameStr)
		if err != nil {
			fmt.Println("error parsing game str", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.mb.BuildWinMap(g); err != nil {
			fmt.Println("error building win map", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	s.r.GET("/chances", func(c *gin.Context) {
		gameStr := c.Query("game")

		g, err := game.FromString(gameStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		stats, err := s.mr.GetGameStats(g)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(stats)

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data":   stats,
		})
	})

	return s
}

func (s *Server) Start(port int) {
	if err := s.r.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}
