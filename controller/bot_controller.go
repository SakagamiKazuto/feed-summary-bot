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

var api = slack.New(os.Getenv("SLACK_APP_TOKEN"))

func HandleBotEvents(c echo.Context) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request().Body)
	body := buf.String()

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("SLACK_VERIFICATION_TOKEN")}))
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
				handleAppMention(innerEventData)
			default:
				log.Printf("[INFO] unsupported inner event: %+v\n", innerEvent.Data)
			}
		default:
			log.Printf("[INFO] unsupported event: %+v\n", eventsAPIEvent.Type)
		}
	}

	return c.String(http.StatusOK, "")
}

func handleAppMention(event *slackevents.AppMentionEvent) {
	command := strings.Split(event.Text, " ")

	if len(command) < 3 {
		api.PostMessage(event.Channel, slack.MsgOptionText("Error: Invalid command format.", false))
		return
	}

	action, url := command[1], command[2]

	switch action {
	case "init":
		// ここにwebhookURLをサーバーに保存する処理を追加
		response := "Webhook URL '" + url + "' has been saved."
		api.PostMessage(event.Channel, slack.MsgOptionText(response, false))
	case "feed":
		// ここに保存されたwebhookURLと紐づく形でfeedURLを保存する処理を追加
		response := "Feed URL '" + url + "' has been saved."
		api.PostMessage(event.Channel, slack.MsgOptionText(response, false))
	default:
		api.PostMessage(event.Channel, slack.MsgOptionText("Error: Invalid command format.", false))
	}
}
