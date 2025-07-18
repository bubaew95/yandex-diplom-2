package model

type BinaryRequest struct {
	Binary []byte `json:"binary"`
}

type BinaryResponse struct {
	ID     int64  `json:"id"`
	Binary []byte `json:"binary"`
	UserID int64  `json:"user_id"`
}
