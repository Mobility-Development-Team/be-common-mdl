package pagination

import (
	"gorm.io/gorm"
	"math"
)

type Pagination struct {
	Limit         int         `json:"limit,omitempty;query:limit"`
	Page          int         `json:"page,omitempty;query:page"`
	Sort          string      `json:"sort,omitempty;query:sort"`
	TotalRows     int64       `json:"totalRows"`
	TotalPages    int         `json:"totalPages"`
	Rows          interface{} `json:"rows"`
	IsPaginate    bool        `json:"isPaginate"`
	DistinctValue string      `json:"-"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "created_at desc"
	}
	return p.Sort
}

func Paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var (
		totalRows int64
	)
	if pagination.DistinctValue != "" {
		// In format eg: count(distinct(`%s`.id))
		db.Model(value).Select(pagination.DistinctValue).Count(&totalRows)
	} else {
		db.Model(value).Count(&totalRows)
	}
	pagination.TotalRows = totalRows
	l := pagination.Limit
	if l == 0 {
		l = 10
	}
	totalPages := int(math.Ceil(float64(totalRows) / float64(l)))
	pagination.TotalPages = totalPages
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
