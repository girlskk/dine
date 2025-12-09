package types

import (
	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

// StoreCreateReq 用于创建门店时传递的请求参数
type StoreCreateReq struct {
	Name                string                      `json:"name" binding:"required"`                            // 门店名称
	UserName            string                      `json:"username" binding:"required"`                        // 门店账号
	Password            string                      `json:"password" binding:"required"`                        // 账号密码
	Type                domain.StoreType            `json:"type" binding:"required,oneof=restaurant cafeteria"` // 门店类型
	CooperationType     domain.StoreCooperationType `json:"cooperation_type" binding:"required,oneof=join"`     // 合作类型
	NeedAudit           bool                        `json:"need_audit"`                                         // 商品是否需要总部审核
	Enabled             bool                        `json:"enabled"`                                            // 门店状态（开启/关闭）
	PointSettlementRate decimal.Decimal             `json:"point_settlement_rate" binding:"d_lte=1"`            // 积分结算费率（单位：百分比，例如 0.1234 表示 12.34%）
	PointWithdrawalRate decimal.Decimal             `json:"point_withdrawal_rate" binding:"d_lte=1"`            // 积分提现费率（单位：百分比，例如 0.1234 表示 12.34%）
	HuifuID             string                      `json:"huifu_id"`                                           // 汇付ID
	ZxhID               string                      `json:"zxh_id"`                                             // 知心话ID
	ZxhSecret           string                      `json:"zxh_secret"`                                         // 知心话密钥
	City                string                      `json:"city"`                                               // 省市地区
	Address             string                      `json:"address"`                                            // 详细地址
	ContactName         string                      `json:"contact_name"`                                       // 门店联系人姓名
	ContactPhone        string                      `json:"contact_phone"`                                      // 联系人电话号码
	Images              domain.StoreInfoImages      `json:"images"`                                             // 门店相关图片数组，包括店标图片、门店正面图等
	BankAccount         string                      `json:"bank_account"`                                       // 银行卡账号
	BankCardName        string                      `json:"bank_card_name"`                                     // 银行账户名称
	BankName            string                      `json:"bank_name"`                                          // 银行名称
	BranchName          string                      `json:"branch_name"`                                        // 开户支行名称
	PublicAccount       string                      `json:"public_account"`                                     // 对公账号
	CompanyName         string                      `json:"company_name"`                                       // 公司名称
	PublicBankName      string                      `json:"public_bank_name"`                                   // 对公银行名称
	PublicBranchName    string                      `json:"public_branch_name"`                                 // 对公开户支行名称
	CreditCode          string                      `json:"credit_code"`                                        // 统一社会信用代码
}

func ToStoreCreateReq(req StoreCreateReq) *domain.Store {
	return &domain.Store{
		Name:                req.Name,
		Type:                req.Type,
		CooperationType:     req.CooperationType,
		NeedAudit:           req.NeedAudit,
		Enabled:             req.Enabled,
		PointSettlementRate: req.PointSettlementRate,
		PointWithdrawalRate: req.PointWithdrawalRate,
		HuifuID:             req.HuifuID,
		ZxhID:               req.ZxhID,
		ZxhSecret:           req.ZxhSecret,
		Info: &domain.StoreInfo{
			City:         req.City,
			Address:      req.Address,
			ContactName:  req.ContactName,
			ContactPhone: req.ContactPhone,
			Images:       req.Images,
		},
		Finance: &domain.StoreFinance{
			BankAccount:      req.BankAccount,
			BankCardName:     req.BankCardName,
			BankName:         req.BankName,
			BranchName:       req.BranchName,
			PublicAccount:    req.PublicAccount,
			CompanyName:      req.CompanyName,
			PublicBankName:   req.PublicBankName,
			PublicBranchName: req.PublicBranchName,
			CreditCode:       req.CreditCode,
		},
	}
}

type StoreUpdateReq struct {
	ID                  int                         `json:"id" binding:"required"`                              // 门店ID
	Name                string                      `json:"name" binding:"required"`                            // 门店名称
	Password            string                      `json:"password"`                                           // 账号密码
	Type                domain.StoreType            `json:"type" binding:"required,oneof=restaurant cafeteria"` // 门店类型
	CooperationType     domain.StoreCooperationType `json:"cooperation_type" binding:"required,oneof=join"`     // 合作类型
	NeedAudit           bool                        `json:"need_audit"`                                         // 商品是否需要总部审核
	Enabled             bool                        `json:"enabled"`                                            // 门店状态（开启/关闭）
	PointSettlementRate decimal.Decimal             `json:"point_settlement_rate" binding:"d_lte=1"`            // 积分结算费率（单位：百分比，例如 0.1234 表示 12.34%）
	PointWithdrawalRate decimal.Decimal             `json:"point_withdrawal_rate" binding:"d_lte=1"`            // 积分提现费率（单位：百分比，例如 0.1234 表示 12.34%）
	HuifuID             string                      `json:"huifu_id"`                                           // 汇付ID
	ZxhID               string                      `json:"zxh_id"`                                             // 知心话ID
	ZxhSecret           string                      `json:"zxh_secret"`                                         // 知心话密钥
	City                string                      `json:"city"`                                               // 省市地区
	Address             string                      `json:"address"`                                            // 详细地址
	ContactName         string                      `json:"contact_name"`                                       // 门店联系人姓名
	ContactPhone        string                      `json:"contact_phone"`                                      // 联系人电话号码
	Images              domain.StoreInfoImages      `json:"images"`                                             // 门店相关图片数组，包括店标图片、门店正面图等
	BankAccount         string                      `json:"bank_account"`                                       // 银行卡账号
	BankCardName        string                      `json:"bank_card_name"`                                     // 银行账户名称
	BankName            string                      `json:"bank_name"`                                          // 银行名称
	BranchName          string                      `json:"branch_name"`                                        // 开户支行名称
	PublicAccount       string                      `json:"public_account"`                                     // 对公账号
	CompanyName         string                      `json:"company_name"`                                       // 公司名称
	PublicBankName      string                      `json:"public_bank_name"`                                   // 对公银行名称
	PublicBranchName    string                      `json:"public_branch_name"`                                 // 对公开户支行名称
	CreditCode          string                      `json:"credit_code"`                                        // 统一社会信用代码
}

func ToStoreUpdateReq(req StoreUpdateReq) *domain.Store {
	return &domain.Store{
		ID:                  req.ID,
		Name:                req.Name,
		Type:                req.Type,
		CooperationType:     req.CooperationType,
		NeedAudit:           req.NeedAudit,
		Enabled:             req.Enabled,
		PointSettlementRate: req.PointSettlementRate,
		PointWithdrawalRate: req.PointWithdrawalRate,
		HuifuID:             req.HuifuID,
		ZxhID:               req.ZxhID,
		ZxhSecret:           req.ZxhSecret,
		Info: &domain.StoreInfo{
			City:         req.City,
			Address:      req.Address,
			ContactName:  req.ContactName,
			ContactPhone: req.ContactPhone,
			Images:       req.Images,
		},
		Finance: &domain.StoreFinance{
			BankAccount:      req.BankAccount,
			BankCardName:     req.BankCardName,
			BankName:         req.BankName,
			BranchName:       req.BranchName,
			PublicAccount:    req.PublicAccount,
			CompanyName:      req.CompanyName,
			PublicBankName:   req.PublicBankName,
			PublicBranchName: req.PublicBranchName,
			CreditCode:       req.CreditCode,
		},
	}
}

type StoreListReq struct {
	Page int    `json:"page"`
	Size int    `json:"size"`
	Name string `json:"name"` // 门店名称
	City string `json:"city"` // 省市地区
}

type StoreIDReq struct {
	ID int `json:"id" binding:"required"`
}
