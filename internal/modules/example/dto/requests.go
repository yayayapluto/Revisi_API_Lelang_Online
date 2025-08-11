package dto

type CreateRequest struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

type UpdateRequest struct {
	Kode *string `json:"kode,omitempty"`
	Nama *string `json:"nama,omitempty"`
}
