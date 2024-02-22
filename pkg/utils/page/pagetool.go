package page

import (
	"math"
	"strconv"
)

const defaultPageSize = 10

type Param struct {
	CurrentPage int    `json:"currentPage"`
	PageSize    int    `json:"pageSize"`
	SortFiled   string `json:"sortFiled"`
	Order       string `json:"order"`
}

type Value struct {
	TotalSize   int           `json:"totalSize"`
	TotalPage   int           `json:"totalPage"`
	CurrentPage int           `json:"currentPage"`
	PageSize    int           `json:"pageSize"`
	Records     []interface{} `json:"records"`
	NoPage      bool          `json:"noPage"`
}

func DoPage(page, size string, records []interface{}) Value {

	// 定义如果不传 第几页 则不分页
	if page == "" {
		param := NewDefaultParam()
		pageValue := NewValue(records, len(records), param)
		pageValue.NoPage = true
		return pageValue
	}

	currentPage, err := strconv.Atoi(page)
	if err != nil {
		currentPage = 1
	}

	var pageSize int
	if size == "" {
		pageSize = 10
	} else {
		pageSize, err = strconv.Atoi(size)
		if err != nil {
			pageSize = 10
		}
	}

	param := NewParam(currentPage, pageSize)
	pageValue := NewValue(records, len(records), param)
	pageValue.PaginateMemory()
	pageValue.NoPage = false

	return pageValue
}

func (p *Value) PaginateMemory() {
	if p.PageSize < 0 {
		//每页条目小于0不分页
		p.CurrentPage = 1
		return
	}
	if p.PageSize == 0 {
		p.PageSize = defaultPageSize
	}
	if p.CurrentPage > p.TotalPage {
		p.CurrentPage = p.TotalPage
	}
	if p.CurrentPage <= 0 {
		p.CurrentPage = 1
	}

	start := (p.CurrentPage - 1) * p.PageSize
	end := min(start+p.PageSize, p.TotalSize)

	p.Records = p.Records[start:end]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (p *Value) calculateTotalPage() {
	if p.PageSize > 0 {
		p.TotalPage = int(math.Ceil(float64(p.TotalSize) / float64(p.PageSize)))
	}
}

func NewValue(records []interface{}, totalSize int, param *Param) Value {
	value := Value{
		TotalSize:   totalSize,
		Records:     records,
		CurrentPage: param.CurrentPage,
		PageSize:    param.PageSize,
	}
	value.calculateTotalPage()
	return value
}

func NewDefaultParam() *Param {
	return &Param{
		CurrentPage: 1,
		PageSize:    defaultPageSize,
	}
}

func NewParam(currentPage, pageSize int) *Param {
	return &Param{
		CurrentPage: currentPage,
		PageSize:    pageSize,
	}
}
