package upagination

const (
	DefaultPage int = 1
	DefaultSize int = 10
	MaxSize     int = 5000 // 防止查询过大数据
)

type Pagination struct {
	Total int `json:"total"` // 总页数
	Page  int `json:"page"`  // 页码
	Size  int `json:"size"`  // 每页数量
}

func New(page, size int) *Pagination {
	// 边界处理
	if page < 1 {
		page = DefaultPage
	}
	if size < 1 {
		size = DefaultSize
	}
	if size > MaxSize {
		size = MaxSize
	}
	return &Pagination{
		Page: page,
		Size: size,
	}
}

// Offset 获取数据库查询的偏移量
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Size
}

// SetTotal 设置总数
func (p *Pagination) SetTotal(total int) {
	p.Total = total
}

type RequestPagination struct {
	Page int `json:"page" form:"page" binding:"omitempty,gt=0"`         // 页码
	Size int `json:"size" form:"size" binding:"omitempty,gt=0,lt=1000"` // 每页数量
}

func (p *RequestPagination) ToPagination() *Pagination {
	return New(p.Page, p.Size)
}

func TotalPages(total, size int) int {
	return (total + size - 1) / size
}
