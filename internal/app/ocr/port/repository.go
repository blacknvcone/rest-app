package port

import "context"

type IHuggingFaceHTTP interface {
	ProceedTxtToJSONGeneratorPrompt(ctx context.Context, txtTarget string) (string, error)
}

type IGoogleAIHTTP interface {
	ProceedTxtToJSONGeneratorPrompt(ctx context.Context, txtTarget string) ([]byte, error)
}
