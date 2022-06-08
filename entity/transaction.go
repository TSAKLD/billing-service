package entity

import "time"

type Transaction struct {
	UserID           string    `json:"user_id"`
	Sender           *int      `json:"sender"`
	Retriever        *int      `json:"retriever"`
	Amount           *int      `json:"amount"`
	Date             time.Time `json:"date"`
	UserCommentary   *string   `json:"user_commentary"`
	SystemCommentary string    `json:"system_commentary"`
}
