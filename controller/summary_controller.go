package controller

import (
	"feed-summary-bot/domain/gateway/openai"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/slack-go/slack"
	"os"
	"time"
)

func postSummaryToSlack(summary, slackChannelID, articleURL string) error {
	slackBotToken := os.Getenv("SLACK_APP_TOKEN")
	api := slack.New(slackBotToken)

	_, _, err := api.PostMessage(
		slackChannelID,
		slack.MsgOptionText(fmt.Sprintf("Summary of the article at %s:\n%s", articleURL, summary), false),
	)
	if err != nil {
		return err
	}
	return nil
}

func getRecentFeedEntries(feedURL string) ([]*gofeed.Item, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, err
	}

	updatedSince := time.Now().Add(-6 * time.Hour)
	var recentEntries []*gofeed.Item
	for _, item := range feed.Items {
		published, _ := time.Parse(time.RFC1123Z, item.Published)
		if published.After(updatedSince) {
			recentEntries = append(recentEntries, item)
		}
	}
	return recentEntries, nil
}

func getArticleSummary(articleURL string) (string, error) {
	openAIKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(openAIKey)

	prompt := fmt.Sprintf(`Please summarize the article at the following URL: %s`, articleURL)

	req := &openai.CompletionRequest{
		Prompt:      prompt,
		MaxTokens:   100,
		Temperature: 0.3,
		TopP:        1,
		N:           1,
	}

	completion, err := client.CreateCompletion(req)
	if err != nil {
		return "", err
	}

	if len(completion.Choices) > 0 {
		return completion.Choices[0].Text, nil
	}

	return "", fmt.Errorf("no summary available")
}
