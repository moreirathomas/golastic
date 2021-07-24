package internal

type Pagination struct {
	Page    int             `json:"page,omitempty"`
	PerPage int             `json:"per_page,omitempty"`
	Link    PaginationLinks `json:"links,omitempty"`
}

type PaginationLinks struct {
	Prev string `json:"prev,omitempty"`
	Next string `json:"next,omitempty"`
	Self string `json:"self,omitempty"`
}

func NewPagination(size int, from int) Pagination {
	return Pagination{
		Page:    from/size + 1,
		PerPage: size,
	}
}
