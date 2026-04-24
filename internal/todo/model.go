package todo

import "time"

type Todo struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Status    string    `json:"status"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
