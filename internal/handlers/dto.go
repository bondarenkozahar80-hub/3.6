package handlers

//структура, которая определяет, какие данные приходят от клиента и какие уходят в ответ, отдельно от самой бизнес-модели или таблицы БД. (применил для приема времени от клиента в нужном формате)

type TransactionDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount" binding:"required"` // от клиента в рублях
	Type        string  `json:"type" binding:"required,oneof=income expense"`
	Category    string  `json:"category"`
	EventDate   string  `json:"event_date" binding:"required"` // "2025-11-26"
}
