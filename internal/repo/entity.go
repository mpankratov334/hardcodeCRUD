package repo

import "time"

// DataObject - шаблонная структура, для хранения
type DataObject struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Data      string    `json:"data"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Task struct {
	DataObject
	UserID string `json:"user_id"`
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
