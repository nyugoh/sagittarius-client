package api

import (
	"github.com/gin-gonic/gin"
	io "github.com/googollee/go-socket.io"
	"github.com/jinzhu/gorm"
	"github.com/nyugoh/sagittarius-client/utils"
)

type App struct {
	DB    *gorm.DB
	Name  string
	Port  string
	Folders []string
	SocketServer *io.Server

}

func (app *App) Index(c *gin.Context) {
	utils.SendJson(c, gin.H{
		"success": "success",
		"message": "Hello, am a sagittarius client",
	})
}