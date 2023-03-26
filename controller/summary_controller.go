package controller

import (
	"feed-summary-bot/domain/gateway/openai"
	"fmt"
	"os"
)

var openAIKey = os.Getenv("OPENAI_API_KEY")

func getArticleSummary(articleURL string) (string, error) {
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
