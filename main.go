package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"tg-gateway/model"

	"github.com/gin-gonic/gin"
)

var systemSecret string

func init() {
	systemSecret, _ = os.LookupEnv("SECRET")
}

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		secret := c.GetHeader("X-Telegram-Bot-Api-Secret-Token")
		if secret == systemSecret {
			c.Next()
			return
		}

		c.Abort()
		c.JSON(http.StatusForbidden, gin.H{
			"message": "invalid secret",
		})
	})

	r.POST("/receivingUpdates", func(c *gin.Context) {
		bytes, _ := io.ReadAll(c.Request.Body)
		fmt.Println(string(bytes))

		up := &model.Update{}
		json.Unmarshal(bytes, up)

		c.JSON(http.StatusOK, gin.H{
			"pong":    string(bytes),
			"pong_up": up,
		})
	})

	r.Run("127.0.0.1:9182")
}
