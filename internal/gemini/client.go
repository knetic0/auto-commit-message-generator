package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	URL             = "https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent"
	MODEL           = "gemini-2.0-flash"
	PROMPT_TEMPLATE = "Generate a commit message for the following changes:\n\n%s\n\nPlease provide a concise and descriptive commit message that summarizes the changes made.\n\nCommit Message:"
)

type GeminiClient struct {
	ApiKey string
}

func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		ApiKey: apiKey,
	}
}

func NewGeminiRequest(text string) *GeminiRequest {
	return &GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: strings.Replace(PROMPT_TEMPLATE, "%s", text, 1)},
				},
			},
		},
	}
}

func (c *GeminiClient) GenerateCommitMessage(request *GeminiRequest) (*GeminiResponse, error) {
	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	url := strings.Replace(URL, "%s", MODEL, 1)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", c.ApiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini API request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))

	var geminiResponse GeminiResponse
	if err := json.Unmarshal(body, &geminiResponse); err != nil {
		return nil, err
	}

	return &geminiResponse, nil
}
