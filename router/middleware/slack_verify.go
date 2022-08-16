package middleware

import (
	"eat-and-go/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"io/ioutil"
	"log"
)

func SlackVerifyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slackConfig := config.GetConfig().Slack
		cCp := ctx.Copy()
		secret := slackConfig.SigningSecret
		sv, err := slack.NewSecretsVerifier(cCp.Request.Header, secret)
		if err != nil {
			log.Fatal("verify error:", err)
			return
		}
		body, err := ioutil.ReadAll(cCp.Request.Body)
		if err  != nil {
			log.Fatal("read err", err)
			return
		}
		if _, err := sv.Write(body); err != nil {
			log.Fatal("sv write err:", err)
			return
		}
		fmt.Println("body is :", body)
		if err := sv.Ensure(); err != nil {
			log.Fatal("ensure err:", err)
			return
		}
		ctx.Next()
	}
}