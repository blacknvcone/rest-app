package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"rest-app/config"
	"rest-app/internal/app/ocr/port"

	"rest-app/pkg/httpclient"
)

type GenerationConfig struct {
	ResponseMimeType string          `json:"responseMimeType,omitempty"`
	ResponseSchema   *ResponseSchema `json:"responseSchema,omitempty"`
}

type ResponseSchema struct {
	Type             string               `json:"type"`
	Items            *SchemaItems         `json:"items,omitempty"`
	Properties       map[string]*Property `json:"properties,omitempty"`
	PropertyOrdering []string             `json:"propertyOrdering,omitempty"`
}

type SchemaItems struct {
	Type             string               `json:"type"`
	Properties       map[string]*Property `json:"properties,omitempty"`
	PropertyOrdering []string             `json:"propertyOrdering,omitempty"`
}

type Property struct {
	Type  string       `json:"type"`
	Items *SchemaItems `json:"items,omitempty"`
}

type Part struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []Part `json:"parts"`
	Role  string `json:"role,omitempty"` // Added role field
}

type Candidate struct {
	Content      Content `json:"content"` // Single content object, not array
	FinishReason string  `json:"finishReason"`
	AvgLogprobs  float64 `json:"avgLogprobs"`
}

type GoogleTextGenerationRequest struct {
	Contents         []Content         `json:"contents"`
	GenerationConfig *GenerationConfig `json:"generationConfig,omitempty"`
}

type GoogleTextGenerationResponse struct {
	Candidates    []Candidate `json:"candidates"` // Fixed: Use Candidate struct
	UsageMetadata interface{} `json:"usageMetadata,omitempty"`
	ModelVersion  string      `json:"modelVersion,omitempty"`
	ResponseId    string      `json:"responseId,omitempty"`
}

type googleaiTextGenerationHTTP struct {
	conf       *config.GoogleAIAPIConf
	httpClient *httpclient.RestClient
}

func NewGoogleAIHTTP(conf *config.GoogleAIAPIConf, httpClient *httpclient.RestClient) port.IGoogleAIHTTP {
	return &googleaiTextGenerationHTTP{
		conf:       conf,
		httpClient: httpClient,
	}
}

func (h *googleaiTextGenerationHTTP) ProceedTxtToJSONGeneratorPrompt(ctx context.Context, txtTarget string) ([]byte, error) {

	const rules = `
	Rules:
		- Ensure the JSON matches the provided format exactly
		- Use empty string "" for missing text fields
		- Use 0.0 for missing numeric fields
		- Extract amounts as numbers without currency symbols
		- Remove any markdown code blocks or backticks from the output
	`

	prompt := fmt.Sprintf("Parse this text below into JSON:%s \n and rules is %s", txtTarget, rules)

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	reqPayload := GoogleTextGenerationRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{
						Text: prompt,
					},
				},
			},
		},
		GenerationConfig: &GenerationConfig{
			ResponseMimeType: "application/json",
			ResponseSchema: &ResponseSchema{
				Type: "object",
				Properties: map[string]*Property{
					"transaction_id": {
						Type: "string",
					},
					"amount": {
						Type: "number",
					},
					"currency": {
						Type: "string",
					},
					"date": {
						Type: "string",
					},
					"time": {
						Type: "string",
					},
					"sender_name": {
						Type: "string",
					},
					"sender_account": {
						Type: "string",
					},
					"receiver_name": {
						Type: "string",
					},
					"receiver_account": {
						Type: "string",
					},
					"bank_name": {
						Type: "string",
					},
					"transaction_type": {
						Type: "string",
					},
					"reference": {
						Type: "string",
					},
					"status": {
						Type: "string",
					},
					"fee": {
						Type: "number",
					},
					"description": {
						Type: "string",
					},
				},
				PropertyOrdering: []string{
					"transaction_id",
					"amount",
					"currency",
					"date",
					"time",
					"sender_name",
					"sender_account",
					"receiver_name",
					"receiver_account",
					"bank_name",
					"transaction_type",
					"reference",
					"status",
					"fee",
					"description",
				},
			},
		},
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", h.conf.URL, h.conf.Model, h.conf.APIToken)
	resp, err := h.httpClient.Post(url, reqPayload, headers)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	var finalResp GoogleTextGenerationResponse

	err = json.Unmarshal(resp.Body(), &finalResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(finalResp.Candidates) == 0 || len(finalResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no valid content found in response")
	}

	jsonText := finalResp.Candidates[0].Content.Parts[0].Text

	var parsedJSON interface{}
	err = json.Unmarshal([]byte(jsonText), &parsedJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON from response text: %w", err)
	}

	finalRespBytes, err := json.Marshal(parsedJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal final response: %w", err)
	}

	return finalRespBytes, nil
}
