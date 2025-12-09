package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/zxh"
)

var _ PaymentProvider = (*ZxhWalletPaymentProvider)(nil)

type ZxhWalletPaymentProvider struct {
	mgr         *zxh.Manager
	credentials *zxh.Credentials
}

func NewZxhWalletPaymentProvider(mgr *zxh.Manager, mchID, accessKey string) *ZxhWalletPaymentProvider {
	return &ZxhWalletPaymentProvider{
		mgr:         mgr,
		credentials: zxh.NewCredentials(mchID, accessKey),
	}
}

func (p *ZxhWalletPaymentProvider) Provider() domain.PayProvider {
	return domain.PayProviderZhiXinHuaWallet
}

func (p *ZxhWalletPaymentProvider) MchID() string {
	return p.credentials.MchID
}

func (p *ZxhWalletPaymentProvider) Payment(ctx context.Context, seqNo string, params *PaymentParams) (res *PaymentResult, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ZxhWalletPaymentProvider.Payment")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	logger := logging.FromContext(ctx).Named("ZxhWalletPaymentProvider.Payment")

	client, release := p.mgr.GetClient(p.credentials)
	defer release()

	axhUserInfo, err := client.GetUserInfo(ctx, params.AuthCode)
	if err != nil {
		logger.Errorf("zxh.GetUserInfo: %v", err)
		return nil, transPointPayRespErr(err, "获取用户信息失败")
	}
	logger.Infof("zxh.GetUserInfo resp: %+v", axhUserInfo)

	respRaw, err := json.Marshal(axhUserInfo)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal axhUserInfo: %w", err)
	}

	amt := params.Amount.StringFixed(2)
	pointLimits := []zxh.PointLimit{
		{
			Code:  DineCode,
			Limit: amt,
		},
	}

	merchantName := params.Get("merchant_name").(string)
	payParam := &zxh.PaymentParam{
		PayCode:      params.AuthCode,
		MerchantName: merchantName,
		OutOrderID:   seqNo,
		Desc:         params.GoodsDesc,
		NotifyURL:    params.NotifyURL,
		Amount:       amt,
		Points:       pointLimits,
	}
	logger.Infof("zxh.PaymentParam: %+v", payParam)

	reqRaw, err := json.Marshal(payParam)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	if err := client.Payment(ctx, payParam); err != nil {
		logger.Errorf("zxh.Payment: %v", err)
		return nil, transPointPayRespErr(err, "积分支付失败")
	}

	return &PaymentResult{
		Req:     reqRaw,
		Resp:    respRaw,
		Channel: domain.PayChannelPointWallet,
	}, nil
}

var _ PaymentCallbackProvider = (*ZxhWalletPaymentCallbackProvider)(nil)

type ZxhWalletPaymentCallbackProvider struct{}

func NewZxhWalletPaymentCallbackProvider() *ZxhWalletPaymentCallbackProvider {
	return &ZxhWalletPaymentCallbackProvider{}
}

func (p *ZxhWalletPaymentCallbackProvider) Callback(ctx context.Context, callback *domain.PaymentCallback, payment *domain.Payment) (res *PaymentCallbackResult, err error) {
	var data zxh.ZhixinhuaPointPayCallBack
	if err := json.Unmarshal(callback.Raw, &data); err != nil {
		return nil, fmt.Errorf("json.Unmarshal resp: %w", err)
	}

	res = new(PaymentCallbackResult)
	switch zxh.Status(data.Status) {
	case zxh.StatusPending:
		return nil, fmt.Errorf("payOrder[%s] pending", data.OutOrderID)
	case zxh.StatusWaiting:
		res.State = domain.PayStateWaiting
		res.FailReason = data.ErrMsg
	case zxh.StatusSuccess:
		res.State = domain.PayStateSuccess
	case zxh.StatusFail:
		res.State = domain.PayStateFailure
		res.FailReason = data.ErrMsg
	default:
		return nil, fmt.Errorf("unknown status[%d]", data.Status)
	}

	if res.State != domain.PayStateSuccess {
		return
	}

	var axhUserInfo zxh.UserInfo
	if err := json.Unmarshal(payment.Resp, &axhUserInfo); err != nil {
		return nil, fmt.Errorf("json.Unmarshal axhUserInfo: %w", err)
	}

	res.MemberInfo = domain.PaymentMemberInfo{
		ID:    int(axhUserInfo.ID),
		Name:  axhUserInfo.Name,
		Phone: axhUserInfo.PhoneNumber,
	}

	return
}
