package domains

import (
	"time"
)

type Memo struct {
	ID string `json:"id"`
	UserID string `json:"-"`
	Title string `json:"title"`
	Content string `json:"content"`
	IsArchived bool `json:"isArchived"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}