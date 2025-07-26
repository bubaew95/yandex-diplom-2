package model

type CardRequest struct {
	ID     int64  `json:"id,omitempty"`
	Number string `json:"number"`
}

type CardResponse struct {
	ID        int64  `json:"id"`
	Number    string `json:"number"`
	UserID    int64  `json:"user_id"`
	IsDeleted bool   `json:"is_deleted,omitempty"`
}
