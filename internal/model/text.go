package model

type TextRequest struct {
	ID   int64  `json:"id,omitempty"`
	Text string `json:"text"`
}

type TextResponse struct {
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	UserID    int64  `json:"user_id"`
	IsDeleted bool   `json:"is_deleted,omitempty"`
}
