package controller

import (
	"bytes"
	"feed-summary-bot/logger"
	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func HandleBotEvents(c echo.Context) error {
	if err := verifyRequest(c); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	cmd := parseSlashCommand(c)
	if err := handleAppMention(cmd); err != nil {
		logger.LOG.Error("handleAppMention failed", zap.Error(err), zap.Any("command", cmd))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid slash command"})
	}

	return c.String(http.StatusOK, "successful posted")
}

func verifyRequest(c echo.Context) error {
	buf := new(bytes.Buffer)
	tee := io.TeeReader(c.Request().Body, buf)
	body, err := ioutil.ReadAll(tee)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error reading request body")
	}
	c.Request().Body = ioutil.NopCloser(buf)

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

func parseSlashCommand(c echo.Context) slack.SlashCommand {
	return slack.SlashCommand{
		Token:       c.FormValue("token"),
		TeamID:      c.FormValue("team_id"),
		TeamDomain:  c.FormValue("team_domain"),
		ChannelID:   c.FormValue("channel_id"),
		ChannelName: c.FormValue("channel_name"),
		UserID:      c.FormValue("user_id"),
		UserName:    c.FormValue("user_name"),
		Command:     c.FormValue("command"),
		Text:        c.FormValue("text"),
		ResponseURL: c.FormValue("response_url"),
		TriggerID:   c.FormValue("trigger_id"),
	}
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
