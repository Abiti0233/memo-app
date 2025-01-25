package bookmark

import (
	"time"
)

type Bookmark struct {
	ID        string    `json:"id"`
	UserID    string    `json:"-"`
	MemoID    string    `json:"memoId"`
	CreatedAt time.Time `json:"createdAt"`
}
