package controller

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func verifyRequest(c echo.Context, body string) error {
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET") // Slackアプリの設定ページから取得したSigning Secretを設定してください
	sv, err := slack.NewSecretsVerifier(c.Request().Header, signingSecret)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	_, err = sv.Write([]byte(body))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	err = sv.Ensure()
	if err != nil {
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
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid slash command"})
	}

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	switch eventsAPIEvent.Type {
	case slackevents.CallbackEvent:
		switch eventsAPIEvent.Type {
		case slackevents.CallbackEvent:
			innerEvent := eventsAPIEvent.InnerEvent
			switch innerEventData := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				handleAppMention(innerEventData, cmd)
			default:
				log.Printf("[INFO] unsupported inner event: %+v\n", innerEvent.Data)
			}
		default:
			log.Printf("[INFO] unsupported event: %+v\n", eventsAPIEvent.Type)
		}
	}

	return c.String(http.StatusOK, "")
}

func handleAppMention(event *slackevents.AppMentionEvent, cmd slack.SlashCommand) {
	api := slack.New(os.Getenv("SLACK_APP_TOKEN"))

	command := strings.Split(event.Text, " ")

	if len(command) < 3 {
		api.PostMessage(event.Channel, slack.MsgOptionText("Error: Invalid command format.", false))
		return
	}

	action, url := command[1], command[2]

	switch action {
	case "feed":
		// ここに保存されたchannelIDとfeedURLを保存する処理を追加
		// channelIDはcmdから、feedURLはurl変数から取得する。
		response := "Feed URL '" + url + "' has been saved."
		api.PostMessage(event.Channel, slack.MsgOptionText(response, false))
	default:
		api.PostMessage(event.Channel, slack.MsgOptionText("Error: Invalid command format.", false))
	}
}
