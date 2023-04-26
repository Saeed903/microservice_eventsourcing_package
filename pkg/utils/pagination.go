package utils

import (
	"fmt"
	"math"
	"strconv"
)

const (
	defaultSize = 10
	defaultPage = 1
)

type PaginationResponse struct {
	TotalCount int64 `json:"totalCount"`
	TotalPages int64 `json:"totalPages"`
	Page       int64 `json:"page"`
	Size       int64 `json:"size"`
	HasMore    bool  `json:"hasMore"`
}

func (p *PaginationResponse) String() string {
	return fmt.Sprintf("TotalCount: %d, TotalPages: %d, Page: %d, size: %d, HasMore: %v", p.TotalCount, p.TotalPages, p.Page, p.Size, p.HasMore)
}

func NewPaginationResponse(totalCount int64, pg *Pagination) *PaginationResponse {
	return &PaginationResponse{
		TotalCount: totalCount,
		TotalPages: int64(pg.GetTotalPages(int(totalCount))),
		Page:       int64(pg.GetPage()),
		Size:       int64(pg.GetSize()),
		HasMore:    pg.GetHasMore(int(totalCount)),
	}
}

type Pagination struct {
	Size    int    `json:"size,omitempty"`
	Page    int    `json:"page,omitempty"`
	OrderBy string `json:"orderBy,omitempty"`
}

func NewPagination(size, page int) *Pagination {
	if size == 0 {
		return &Pagination{Size: defaultSize, Page: defaultPage}
	}
	return &Pagination{Size: size, Page: page}
}

func NewPaginationFromQueryParams(size string, page string) *Pagination {
	p := &Pagination{
		Size: defaultSize, Page: 1,
	}

	if sizeNum, err := strconv.Atoi(size); err == nil && sizeNum != 0 {
		p.Size = sizeNum
	}

	if pageNum, err := strconv.Atoi(page); err == nil && pageNum != 0 {
		p.Page = pageNum
	}

	return p
}

func (p *Pagination) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		p.Size = defaultSize
		return nil
	}
	n, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return err
	}
	p.Size = n

	return nil
}

func (p *Pagination) SetPage(pageQuery string) error {
	if pageQuery == "" {
		p.Page = defaultPage
		return nil
	}

	m, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	p.Page = m
	return nil
}

func (p *Pagination) SetOrderBy(orderByQuery string) {
	p.OrderBy = orderByQuery
}

func (p *Pagination) GetOffSet() int {
	if p.Page == 0 {
		return 0
	}
	return (p.Page - 1) * p.Size
}

func (p *Pagination) GetLimit() int {
	return p.Size
}

func (p *Pagination) GetOrderBy() string {
	return p.OrderBy
}

func (p *Pagination) GetPage() int {
	return p.Page
}

func (p *Pagination) GetSize() int {
	return p.Size
}

func (p *Pagination) GetQueryString() string {
	return fmt.Sprintf("page=%d&size=%d&orderBy=%s", p.Page, p.Size, p.OrderBy)
}

func (p *Pagination) GetTotalPages(totalCount int) int {
	d := float64(totalCount) / float64(p.GetSize())
	return int(math.Ceil(d))
}

func (p *Pagination) GetHasMore(totalCount int) bool {
	return p.GetPage() < totalCount/p.GetSize()
}
