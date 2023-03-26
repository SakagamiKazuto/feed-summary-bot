package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	OpenAIEndpoint = "https://api.openai.com/v1/engines/davinci-codex/completions"
)

type OpenAIClient struct {
	APIKey string
	Client *http.Client
}

type CompletionRequest struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	N           int     `json:"n"`
	Echo        bool    `json:"echo"`
	Stop        string  `json:"stop"`
}

type CompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func NewClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		APIKey: apiKey,
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (c *OpenAIClient) CreateCompletion(req *CompletionRequest) (*CompletionResponse, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", OpenAIEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to create completion")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var completionResponse CompletionResponse
	err = json.Unmarshal(body, &completionResponse)
	if err != nil {
		return nil, err
	}

	return &completionResponse, nil
}
