package model

type DataType string

const (
	LoginPassword DataType = "auth"
	TextData      DataType = "text"
	BinaryData    DataType = "byte"
	CardData      DataType = "card"
)

type Data struct {
	Text string   `json:"text"`
	Type DataType `json:"type"`
}

type BinaryRequest struct {
	Binary []byte `json:"binary"`
}

type BinaryResponse struct {
	ID     int64  `json:"id"`
	Binary []byte `json:"binary"`
	UserID int64  `json:"user_id"`
}

type CardResponse struct {
	ID        int64  `json:"id"`
	Number    string `json:"number"`
	UserID    int64  `json:"user_id"`
	IsDeleted bool   `json:"is_deleted,omitempty"`
}

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

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CardDataContent struct {
	CardNumber string `json:"number"`
	CardHolder string `json:"holder"`
	ExpiryDate string `json:"expiry_date"`
	CVV        string `json:"cvv"`
}

type BinaryDataContent struct {
	FileName string `json:"file_name"`
	Data     []byte `json:"data"`
}
