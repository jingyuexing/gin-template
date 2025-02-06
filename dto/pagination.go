package dto

type Pagination struct {
	Page int `json:"page" form:"page" validate:"required,numeric,min=1" validateMsg:"the minimum number of pages is 1"`
	Size int `json:"size" form:"size" validate:"required,numeric,gte=1,lte=20" validateMsg:"the size out of range, must between 1 and 20"`
}
type PaginationOption struct {
	Page int `json:"page" form:"page" validate:"omitempty,numeric,min=1" validateMsg:"the minimum number of pages is 1"`
	Size int `json:"size" form:"size" validate:"omitempty,numeric,gte=1,lte=20" validateMsg:"the size out of range, must between 1 and 20"`
}
