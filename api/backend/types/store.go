package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

type StoreUpdateReq struct {
	Password         string                 `json:"password"`                                           // 账号密码
	Type             domain.StoreType       `json:"type" binding:"required,oneof=restaurant cafeteria"` // 门店类型
	City             string                 `json:"city"`                                               // 省市地区
	Address          string                 `json:"address"`                                            // 详细地址
	ContactName      string                 `json:"contact_name"`                                       // 门店联系人姓名
	ContactPhone     string                 `json:"contact_phone"`                                      // 联系人电话号码
	Images           domain.StoreInfoImages `json:"images"`                                             // 门店相关图片数组，包括店标图片、门店正面图等
	BankAccount      string                 `json:"bank_account"`                                       // 银行卡账号
	BankCardName     string                 `json:"bank_card_name"`                                     // 银行账户名称
	BankName         string                 `json:"bank_name"`                                          // 银行名称
	BranchName       string                 `json:"branch_name"`                                        // 开户支行名称
	PublicAccount    string                 `json:"public_account"`                                     // 对公账号
	CompanyName      string                 `json:"company_name"`                                       // 公司名称
	PublicBankName   string                 `json:"public_bank_name"`                                   // 对公银行名称
	PublicBranchName string                 `json:"public_branch_name"`                                 // 对公开户支行名称
	CreditCode       string                 `json:"credit_code"`                                        // 统一社会信用代码
}

func ToStoreUpsetReq(req StoreUpdateReq) *domain.Store {
	return &domain.Store{
		Type: req.Type,
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

type StoreIDReq struct {
	ID int `json:"id" binding:"required"`
}
