package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"rest-app/config"
	"rest-app/internal/app/ocr/port"
	"rest-app/pkg/httpclient"
	"strings"
)

type HuggingFaceRequest struct {
	Inputs     string                 `json:"inputs"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

type HuggingFaceResponse []struct {
	GeneratedText string `json:"generated_text"`
}

type huggingFaceHTTP struct {
	conf       *config.HuggingFaceAPIConf
	httpClient *httpclient.RestClient
}

func NewHuggingFaceHTTP(conf *config.HuggingFaceAPIConf, httpClient *httpclient.RestClient) port.IHuggingFaceHTTP {
	return &huggingFaceHTTP{
		conf:       conf,
		httpClient: httpClient,
	}
}

func (h *huggingFaceHTTP) ProceedTxtToJSONGeneratorPrompt(ctx context.Context, txtTarget string) (string, error) {
	const jsonFormat = `{
		"transaction_id": "",
		"amount": 0.0,
		"currency": "",
		"date": "",
		"time": "",
		"sender_name": "",
		"sender_account": "",
		"receiver_name": "",
		"receiver_account": "",
		"bank_name": "",
		"transaction_type": "",
		"reference": "",
		"status": "",
		"fee": 0.0,
		"description": ""
	}`

	const rules = `
	Rules:
		- Return ONLY the JSON object, no other text or explanation including the prompt
		- Ensure the JSON matches the provided format exactly
		- Use empty string "" for missing text fields
		- Use 0.0 for missing numeric fields
		- Extract amounts as numbers without currency symbols
		- Remove any markdown code blocks or backticks from the output
	`

	prompt := fmt.Sprintf("Parse this text below into JSON:%s \n with format %s \n and rules is %s", txtTarget, jsonFormat, rules)

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", h.conf.APIToken),
		"Content-Type":  "application/json",
	}

	reqPayload := HuggingFaceRequest{
		Inputs: prompt,
		// Parameters: map[string]interface{}{
		// 	"max_new_tokens":   200,
		// 	"temperature":      0.2, // Lower temperature for more deterministic output
		// 	"do_sample":        false,
		// 	"return_full_text": false,
		// },
	}

	url := fmt.Sprintf("%s/models/%s", strings.TrimSuffix(h.conf.URL, "/"), h.conf.Model)
	resp, err := h.httpClient.Post(url, reqPayload, headers)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	var apiResponse HuggingFaceResponse
	if err := json.Unmarshal(resp.Body(), &apiResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(apiResponse) == 0 || apiResponse[0].GeneratedText == "" {
		return "", fmt.Errorf("empty response from API")
	}

	// Clean up the response
	generatedText := strings.TrimSpace(apiResponse[0].GeneratedText)
	generatedText = strings.Trim(generatedText, "`") // Remove markdown code blocks if present

	// Validate it's valid JSON
	var js map[string]interface{}
	if err := json.Unmarshal([]byte(generatedText), &js); err != nil {
		return "", fmt.Errorf("API response is not valid JSON: %w", err)
	}

	return generatedText, nil
}
