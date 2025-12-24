package errcode

// 业务错误码
type ErrCode string

func (ec ErrCode) String() string {
	return string(ec)
}

const (
	Success ErrCode = "SUCCESS" // Success 表示成功
	// 请求错误
	InvalidParams ErrCode = "INVALID_PARAMS" //   参数错误
	BadRequest    ErrCode = "BAD_REQUEST"    //   请求错误
	// 认证错误
	Unauthorized ErrCode = "UNAUTHORIZED" // 未授权（认证失败、token无效等）
	// 授权错误
	Forbidden ErrCode = "FORBIDDEN" // 禁止访问
	// 资源错误
	NotFound ErrCode = "NOT_FOUND" // 资源不存在
	Conflict ErrCode = "CONFLICT"  // 资源冲突

	// 系统错误
	InternalError ErrCode = "INTERNAL_ERROR" // 系统内部错误
	UnknownError  ErrCode = "UNKNOWN_ERROR"  // 未知错误

	// 业务错误
	UserNotFound                  ErrCode = "USER_NOT_FOUND"                    // 用户不存在
	CategoryNameExists            ErrCode = "CATEGORY_NAME_EXISTS"              // 商品分类名称已存在
	CategoryDeleteHasChildren     ErrCode = "CATEGORY_DELETE_HAS_CHILDREN"      // 商品分类有子分类
	CategoryDeleteHasProducts     ErrCode = "CATEGORY_DELETE_HAS_PRODUCTS"      // 商品分类有商品
	MerchantNameExists            ErrCode = "MERCHANT_NAME_EXISTS"              // 商户名称已存在
	StoreNameExists               ErrCode = "STORE_NAME_EXISTS"                 // 门店名称已存在
	StoreBusinessHoursConflict    ErrCode = "STORE_BUSINESS_HOURS_CONFLICT"     // 门店营业时间冲突
	StoreBusinessHoursTimeInvalid ErrCode = "STORE_BUSINESS_HOURS_TIME_INVALID" // 门店营业时间无效,开始时间不能晚于结束时间
	StoreDiningPeriodConflict     ErrCode = "STORE_DINING_PERIOD_CONFLICT"      // 门店用餐时段时间冲突
	StoreDiningPeriodTimeInvalid  ErrCode = "STORE_DINING_PERIOD_TIME_INVALID"  // 门店用餐时段时间无效,开始时间不能晚于结束时间
	StoreDiningPeriodNameExists   ErrCode = "STORE_DINING_PERIOD_NAME_EXISTS"   // 门店用餐时段名称已存在
	StoreShiftTimeConflict        ErrCode = "STORE_SHIFT_TIME_CONFLICT"         // 门店班次时间冲突
	StoreShiftTimeTimeInvalid     ErrCode = "STORE_SHIFT_TIME_TIME_INVALID"     // 门店班次时间无效,开始时间不能晚于结束时间
	StoreShiftTimeNameExists      ErrCode = "STORE_SHIFT_TIME_NAME_EXISTS"      // 门店班次名称已存在
	RemarkNameExists              ErrCode = "REMARK_NAME_EXISTS"                // 备注名称已存在
	RemarkDeleteSystem            ErrCode = "REMARK_DELETE_SYSTEM"              // 不能删除系统内置备注
)
