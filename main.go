package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"tg-gateway/api/azure_tts"
	"tg-gateway/api/tg_wrapper"
	"tg-gateway/model"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var (
	systemPath   string
	systemSecret string
	systemBotKey string
	systemTTSKey string
	allowedUsers map[string]byte
)

func init() {
	systemSecret = os.Getenv("SECRET")
	systemBotKey = os.Getenv("BOT_KEY")
	systemTTSKey = os.Getenv("TTS_KEY")
	allowedUsers = make(map[string]byte)
	systemPath = "."

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

		ttsID := fmt.Sprintf("%x", md5.Sum([]byte(up.Message.Text)))
		filePath := fmt.Sprintf("%s/%s.mpga", systemPath, ttsID)

		ttsBody, _ := os.ReadFile(filePath)
		if len(ttsBody) == 0 {
			var ttsError error
			ttsBody, ttsError = azure_tts.TTS(systemTTSKey, up.Message.Text)
			if ttsError != nil {
				c.JSON(http.StatusOK, gin.H{
					"tts_error": ttsError,
				})
				return
			}
		} else {
			fmt.Println("hit cache", filePath)
		}

		fmt.Printf("path [%s] len [%d]\n", filePath, len(ttsBody))

		if err := os.WriteFile(filePath, ttsBody, 0644); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"tts_write_file_error": err,
			})
			return
		}

		tg := tg_wrapper.New(systemBotKey)
		svBodyRaw, svErr := tg.SendVoice(&tg_wrapper.SendVoiceOption{
			ChatID:   up.Message.From.ID,
			FileName: filePath,
		})

		if svErr != nil {
			fmt.Printf("send voice error [%+v]\n", svErr)
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
