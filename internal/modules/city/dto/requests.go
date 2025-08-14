package city

type CreateRequest struct {
	Code     string `json:"code"`
	Nama     string `json:"nama"`
	FullCode string `json:"full_code"`
}

type UpdateRequest struct {
	Code     *string `json:"code,omitempty"`
	Nama     *string `json:"nama,omitempty"`
	FullCode *string `json:"full_code,omitempty"`
}
