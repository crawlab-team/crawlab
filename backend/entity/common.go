package entity

import "strconv"

type Page struct {
	Skip int
	Limit int
	PageNum int
	PageSize int
}

func (p *Page)GetPage(pageNum string, pageSize string) {
	p.PageNum, _ = strconv.Atoi(pageNum)
	p.PageSize, _ = strconv.Atoi(pageSize)
	p.Skip = p.PageSize * (p.PageNum - 1)
	p.Limit = p.PageSize
}