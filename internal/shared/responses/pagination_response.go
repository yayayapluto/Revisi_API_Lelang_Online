package responses

type PaginationResponse[T any] struct {
	CurrentPage    int     `json:"current_page"`
	CurrentPageUrl string  `json:"current_page_url"`
	Data           []T     `json:"data"`
	FirstPageURL   *string `json:"first_page_url"`
	NextPageURL    *string `json:"next_page_url"`
	PerPage        int     `json:"per_page"`
	PrevPageURL    *string `json:"prev_page_url"`
}

func NewPaginationResponse[T any](currentPage int, currentPageUrl string, perPage int, data []T, firstPageURL, nextPageURL, prevPageURL *string) PaginationResponse[T] {
	return PaginationResponse[T]{
		CurrentPage:    currentPage,
		CurrentPageUrl: currentPageUrl,
		Data:           data,
		FirstPageURL:   firstPageURL,
		NextPageURL:    nextPageURL,
		PerPage:        perPage,
		PrevPageURL:    prevPageURL,
	}
}
