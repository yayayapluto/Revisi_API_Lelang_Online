package country

type CreateRequest struct {
	Kode  string `json:"kode"`
	Nama  string `json:"nama"`
	Nomor string `json:"nomor"`
}

type UpdateRequest struct {
	Kode  *string `json:"kode,omitempty"`
	Nama  *string `json:"nama,omitempty"`
	Nomor *string `json:"nomor,omitempty"`
}
