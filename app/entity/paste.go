package entity

import "time"

type TextPaste struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Content      string    `json:"content"`
	Highlighting string    `json:"highlighting"`
	ExpireDays   int       `json:"expire_days"`
	GeneratedAt  time.Time `json:"generated_at,omitempty"`
}



