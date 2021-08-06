package response

type EmptyObject struct{}

// EmptyErrorRes example
type EmptyErrorRes struct {
	Success   bool        `json:"success" example:"false"`
	Status    string      `json:"status" example:"Status string corresponding to the error"`
	Message   string      `json:"message" example:"failure message"`
	Timestamp string      `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      EmptyObject `json:"data"`
}

// EmptyListErrorRes example
type EmptyListErrorRes struct {
	Success   bool          `json:"success" example:"false"`
	Status    string        `json:"status" example:"Status string corresponding to the error"`
	Message   string        `json:"message" example:"failure message"`
	Timestamp string        `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      []EmptyObject `json:"data"`
}

// EmptySuccessRes example
type EmptySuccessRes struct {
	Success   bool        `json:"success" example:"false"`
	Status    string      `json:"status" example:"OK"`
	Message   string      `json:"message" example:"success message"`
	Timestamp string      `json:"timestamp" example:"2006-01-02T15:04:05.000Z"`
	Data      EmptyObject `json:"data"`
}

type ListMeta struct {
	Page  int64 `json:"page"`
	Pages int64 `json:"pages"`
	Limit int64 `json:"Limit"`
	Count int64 `json:"count"`
}

func GetListMeta(page, limit, count int64) *ListMeta {
	pages := count / limit
	if count%limit != 0 {
		pages++
	}

	return &ListMeta{
		Page:  page,
		Pages: pages,
		Limit: limit,
		Count: count,
	}
}
