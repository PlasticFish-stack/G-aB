package tool

import "fmt"

var _ UseLimits = new(ResponseLimits)

type UseLimits interface {
	GetTotal() int
	GetPageSize() int
	GetPageNum() int
	GetPages() int
}

type RequestLimits struct {
	PageSize int `json:"pageSize"`
	PageNum  int `json:"pageNum"`
}

type ResponseLimits struct {
	Total    int64 `json:"total"`
	PageSize int   `json:"pageSize"`
	PageNum  int   `json:"pageNum"`
	Pages    int   `json:"pages"`
}

func NewLimits(total int64, pageSize, pageNum int) *ResponseLimits {
	pages := (int(total) + pageSize - 1) / pageSize
	return &ResponseLimits{
		Total:    total,
		PageSize: pageSize,
		PageNum:  pageNum,
		Pages:    pages,
	}
}

func (l *RequestLimits) GetOffset() (offset int, err error) {
	if l.PageSize == 0 || l.PageNum == 0 {
		return 0, fmt.Errorf("PageSize或PageNum不可为0")
	}
	return (l.PageNum - 1) * l.PageSize, nil
}

func (l *ResponseLimits) GetTotal() int {
	return int(l.Total)
}

func (l *ResponseLimits) GetPageSize() int {
	return l.PageSize
}

func (l *ResponseLimits) GetPageNum() int {
	return l.PageNum
}

func (l *ResponseLimits) GetPages() int {
	return l.Pages
}
