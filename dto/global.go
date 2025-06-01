package dto

type Pagination struct {
	Page  int8 `json:"page" binding:"required"`
	Limit int8 `json:"limit" binding:"required"`
}

type PaginationResponse struct {
	Page      int8  `json:"page" binding:"required"`
	Limit     int8  `json:"limit" binding:"required"`
	Total     int64 `json:"total"`
	TotalPage float64 `json:"total_page"`
}
