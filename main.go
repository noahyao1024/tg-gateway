package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"tg-gateway/api/tg_wrapper"
	"tg-gateway/model"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var (
	systemSecret string
	systemBotKey string
	allowedUsers map[string]byte
)

func init() {
	systemSecret = os.Getenv("SECRET")
	systemBotKey = os.Getenv("BOT_KEY")
	allowedUsers = make(map[string]byte)

	for _, u := range strings.Split(os.Getenv("ALLOWED_USERS"), ",") {
		allowedUsers[u] = '1'
	}
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
		up := &model.Update{}
		json.Unmarshal(bytes, up)

		fmt.Printf("%+v\n", up)

		if allowedUsers[up.Message.From.Username] != '1' {
			c.JSON(http.StatusForbidden, gin.H{
				"message": fmt.Sprintf("invalid user [%+v]", up.Message.From),
			})
			return
		}

		tg := tg_wrapper.New(systemBotKey)
		svBodyRaw, svErr := tg.SendVoice(&tg_wrapper.SendVoiceOption{
			ChatID:   up.Message.From.ID,
			FileName: "/Users/noah/Downloads/response.mpga",
		})

		if svErr != nil {
			fmt.Println("send voice error", svErr)
		}

		svBody := make(map[string]interface{})
		json.Unmarshal([]byte(svBodyRaw), &svBody)

		c.JSON(http.StatusOK, gin.H{
			"send_voice_err":  svErr,
			"send_voice_body": svBody,
		})
	})

	r.Run("127.0.0.1:9182")
}
