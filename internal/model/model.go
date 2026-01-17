package model

type Transaction struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Amount      int64  `json:"amount" db:"amount"` // в копейках
	Type        string `json:"type" db:"type"`
	Category    string `json:"category" db:"category"`
	EventDate   string `json:"event_date" db:"event_date"`
}
