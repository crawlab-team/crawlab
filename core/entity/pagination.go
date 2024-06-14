package entity

import "github.com/crawlab-team/crawlab/core/constants"

type Pagination struct {
	Page int `form:"page" url:"page"`
	Size int `form:"size" url:"size"`
}

func (p *Pagination) IsZero() (ok bool) {
	return p.Page == 0 &&
		p.Size == 0
}

func (p *Pagination) IsDefault() (ok bool) {
	return p.Page == constants.PaginationDefaultPage &&
		p.Size == constants.PaginationDefaultSize
}
