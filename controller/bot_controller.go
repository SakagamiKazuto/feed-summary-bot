package controller

import (
	"bytes"
	"feed-summary-bot/logger"
	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func verifyRequest(c echo.Context, body string) error {
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET") // Slackアプリの設定ページから取得したSigning Secretを設定してください
	sv, err := slack.NewSecretsVerifier(c.Request().Header, signingSecret)
	if err != nil {
		logger.LOG.Error("NewSecretsVerifier error", zap.Error(err))
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	_, err = sv.Write([]byte(body))
	if err != nil {
		logger.LOG.Error("sv.Write error", zap.Error(err))
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	err = sv.Ensure()
	if err != nil {
		logger.LOG.Error("sv.Ensure error", zap.Error(err))
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}
	return nil
}

func HandleBotEvents(c echo.Context) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request().Body)
	body := buf.String()

	if err := verifyRequest(c, body); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	cmd, err := slack.SlashCommandParse(c.Request())
	if err != nil {
		logger.LOG.Error("SlashParse failed", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid slash command"})
	}
	err = handleAppMention(cmd)
	if err != nil {
		logger.LOG.Error("handleAppMention failed", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid slash command"})
	}

	return c.String(http.StatusOK, "successful posted")
}

func handleAppMention(cmd slack.SlashCommand) (err error) {
	api := slack.New(os.Getenv("SLACK_APP_TOKEN"))

	switch cmd.Command {
	case "/summary-feed":
		// ここに保存されたchannelIDとfeedURLを保存する処理を追加
		// channelIDはcmdから、feedURLはurl変数から取得する。
		response := "Feed URL '" + cmd.Text + "' has been saved."
		_, _, err = api.PostMessage(cmd.ChannelID, slack.MsgOptionText(response, false))
	default:
		logger.LOG.Error("invalid command format")
		_, _, err = api.PostMessage(cmd.ChannelID, slack.MsgOptionText("Error: Invalid command format.", false))
	}
	return
}
