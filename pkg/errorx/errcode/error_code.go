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
	UserNotFound                 ErrCode = "USER_NOT_FOUND"                   // 用户不存在
	CategoryNameExists           ErrCode = "CATEGORY_NAME_EXISTS"             // 商品分类名称已存在
	CategoryDeleteHasChildren    ErrCode = "CATEGORY_DELETE_HAS_CHILDREN"     // 商品分类有子分类
	CategoryDeleteHasProducts    ErrCode = "CATEGORY_DELETE_HAS_PRODUCTS"     // 商品分类有商品
	ProductUnitNameExists        ErrCode = "PRODUCT_UNIT_NAME_EXISTS"         // 商品单位名称已存在
	ProductUnitDeleteHasProducts ErrCode = "PRODUCT_UNIT_DELETE_HAS_PRODUCTS" // 商品单位有商品
	ProductSpecNameExists        ErrCode = "PRODUCT_SPEC_NAME_EXISTS"         // 商品规格名称已存在
	ProductSpecDeleteHasProducts ErrCode = "PRODUCT_SPEC_DELETE_HAS_PRODUCTS" // 商品规格有商品
	ProductTagNameExists         ErrCode = "PRODUCT_TAG_NAME_EXISTS"          // 商品标签名称已存在
	ProductTagDeleteHasProducts  ErrCode = "PRODUCT_TAG_DELETE_HAS_PRODUCTS"  // 商品标签有商品

	// 商品口味做法
	ProductAttrNameExists            ErrCode = "PRODUCT_ATTR_NAME_EXISTS"              // 商品口味做法名称已存在
	ProductAttrDeleteHasItems        ErrCode = "PRODUCT_ATTR_DELETE_HAS_ITEMS"         // 商品口味做法有子项
	ProductAttrItemDeleteHasProducts ErrCode = "PRODUCT_ATTR_ITEM_DELETE_HAS_PRODUCTS" // 商品口味做法项有商品

	// 商品
	ProductNameExists ErrCode = "PRODUCT_NAME_EXISTS" // 商品名称已存在
)
