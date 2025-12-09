package types

type AreaCreateReq struct {
	Name string `json:"name" binding:"required"` // 区域名称
}

type AreaUpdateReq struct {
	ID   int    `json:"id" binding:"required"`   // 区域ID
	Name string `json:"name" binding:"required"` // 区域名称
}

type AreaDeleteReq struct {
	ID int `json:"id" binding:"required"` // 区域ID
}

type AreaListReq struct {
	Page int `json:"page"` // 页码
	Size int `json:"size"` // 每页数量
}

type TableCreateReq struct {
	AreaID    int    `json:"area_id" binding:"required"` // 区域ID
	Name      string `json:"name" binding:"required"`    // 台桌名称
	SeatCount int    `json:"seat_count" binding:"min=1"` // 座位数
}

type TableUpdateReq struct {
	ID        int    `json:"id" binding:"required"`      // 台桌ID
	AreaID    int    `json:"area_id" binding:"required"` // 区域ID
	Name      string `json:"name" binding:"required"`    // 台桌名称
	SeatCount int    `json:"seat_count" binding:"min=1"` // 座位数
}

type TableDeleteReq struct {
	ID int `json:"id" binding:"required"` // 台桌ID
}

type TableListReq struct {
	AreaID int `json:"area_id"` // 区域ID
	Page   int `json:"page"`    // 页码
	Size   int `json:"size"`    // 每页数量
}
