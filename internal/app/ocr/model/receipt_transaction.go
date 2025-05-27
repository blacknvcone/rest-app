package model

type ReceiptTransaction struct {
	TransactionID   string  `json:"transaction_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	Date            string  `json:"date"`
	Time            string  `json:"time"`
	SenderName      string  `json:"sender_name"`
	SenderAccount   string  `json:"sender_account"`
	ReceiverName    string  `json:"receiver_name"`
	ReceiverAccount string  `json:"receiver_account"`
	BankName        string  `json:"bank_name"`
	TransactionType string  `json:"transaction_type"`
	Reference       string  `json:"reference"`
	Status          string  `json:"status"`
	Fee             float64 `json:"fee"`
	Description     string  `json:"description"`
}
