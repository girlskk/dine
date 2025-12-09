package domain

// Operator 操作者接口
type Operator interface {
	GetOperatorID() int
	GetOperatorName() string
	GetOperatorType() OperatorType
	GetOperatorStoreID() int
}

type OperatorType string

const (
	OperatorTypeFrontend OperatorType = "frontend" // 前台用户
	OperatorTypeBackend  OperatorType = "backend"  // 后台用户
	OperatorTypeAdmin    OperatorType = "admin"    // 管理员
	OperatorTypeSystem   OperatorType = "system"   // 系统
	OperatorTypeCustomer OperatorType = "customer" // 顾客
)

func (OperatorType) Values() []string {
	return []string{
		string(OperatorTypeFrontend),
		string(OperatorTypeBackend),
		string(OperatorTypeAdmin),
		string(OperatorTypeSystem),
		string(OperatorTypeCustomer),
	}
}

type OperatorInfo struct {
	Type OperatorType
	ID   int
	Name string
}

func ExtractOperatorInfo(operator any) *OperatorInfo {
	if operator == nil {
		return &OperatorInfo{
			Type: OperatorTypeSystem,
			ID:   0,
			Name: "系统",
		}
	}
	switch v := operator.(type) {
	default:
		panic("invalid operator type")
	case *FrontendUser:
		return &OperatorInfo{
			Type: OperatorTypeFrontend,
			ID:   v.ID,
			Name: v.Nickname,
		}
	case *BackendUser:
		return &OperatorInfo{
			Type: OperatorTypeBackend,
			ID:   v.ID,
			Name: v.Nickname,
		}
	case *AdminUser:
		return &OperatorInfo{
			Type: OperatorTypeAdmin,
			ID:   v.ID,
			Name: v.Nickname,
		}
	case *Customer:
		return &OperatorInfo{
			Type: OperatorTypeCustomer,
			ID:   v.ID,
			Name: v.Phone, // 顾客手机号作为昵称
		}
	}
}
