package port

import (
	"context"
	"rest-app/internal/app/ocr/model"
)

type IOCRService interface {
	ReceiptDataGenerator(ctx context.Context, imgBytes []byte) (*model.ReceiptTransaction, error)
}
