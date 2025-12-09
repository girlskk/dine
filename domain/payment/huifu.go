package payment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/huifu"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ PaymentProvider = (*HuifuPaymentProvider)(nil)

type HuifuPaymentProvider struct {
	BsPay   *huifu.BsPay
	HuifuID string
}

func NewHuifuPaymentProvider(bsPay *huifu.BsPay, huifuID string) *HuifuPaymentProvider {
	return &HuifuPaymentProvider{
		BsPay:   bsPay,
		HuifuID: huifuID,
	}
}

func (p *HuifuPaymentProvider) Provider() domain.PayProvider {
	return domain.PayProviderHuifu
}

func (p *HuifuPaymentProvider) MchID() string {
	return p.HuifuID
}

func (p *HuifuPaymentProvider) Payment(ctx context.Context, seqNo string, params *PaymentParams) (res *PaymentResult, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HuifuPaymentProvider.Payment")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	logger := logging.FromContext(ctx).Named("HuifuPaymentProvider.Payment")

	expire := time.Now().Add(time.Minute)
	// 发起支付请求
	extendInfos := map[string]interface{}{
		"time_expire": expire.Format("20060102150405"),
		"notify_url":  params.NotifyURL,
	}

	huifuReq := huifu.V2TradePaymentMicropayRequest{
		ReqSeqId:      seqNo,
		ReqDate:       huifu.GetCurrentDate(),
		HuifuId:       p.HuifuID,
		TransAmt:      params.Amount.StringFixed(2),
		GoodsDesc:     params.GoodsDesc,
		AuthCode:      params.AuthCode,
		RiskCheckData: fmt.Sprintf(`{"ip_addr": "%s"}`, params.IPAddr),
		ExtendInfos:   extendInfos,
	}

	logger.Infof("huifu pay req: %+v", huifuReq)

	reqParamRaw, err := json.Marshal(huifuReq)
	if err != nil {
		err = fmt.Errorf("json marshal huifuReq: %w", err)
		return nil, err
	}

	// 发起支付请求
	resp, err := p.BsPay.V2TradePaymentMicropayRequest(ctx, huifuReq)
	if err != nil {
		err = fmt.Errorf("huifu.V2TradePaymentMicropayRequest: %w", err)
		return nil, err
	}

	logger.Infof("huifu pay sync resp: %v", resp)

	respData, ok := resp["data"].(map[string]interface{})
	if !ok {
		err = errors.New("huifu resp format error")
		return nil, err
	}

	respCode, ok := respData["resp_code"].(string)
	if !ok {
		err = errors.New("huifu resp format error")
		return nil, err
	}

	if respCode != huifu.CodeSuccess && respCode != huifu.CodeProcessing { // 发起阶段失败直接返回错误
		var errMsg string
		if msg, ok := respData["bank_desc"].(string); ok {
			errMsg = msg
		} else if msg, ok := respData["resp_desc"].(string); ok {
			errMsg = msg
		} else {
			errMsg = fmt.Sprintf("错误码: %s", respCode)
		}
		err = domain.ParamsErrorf(errMsg)
		return nil, err
	}

	// 同步返回的原数据
	respDataRaw, err := json.Marshal(respData)
	if err != nil {
		logger.Errorf("huifu respData: %v", respData)
		logger.Errorf("json marshal respData: %v", err)
		err = nil
	}

	// 交易类型
	tradeType, ok := respData["trade_type"].(string)
	if !ok {
		return nil, errors.New("trade_type not found in respData")
	}

	channel, err := huifuTradeTypeToDomain(tradeType)
	if err != nil {
		err = fmt.Errorf("failed to convert trade type to domain channel: %w", err)
		return nil, err
	}

	res = &PaymentResult{
		Req:     reqParamRaw,
		Resp:    respDataRaw,
		Channel: channel,
	}

	return
}

func huifuTradeTypeToDomain(tradeType string) (domain.PayChannel, error) {
	switch tradeType {
	case huifu.ChannelWechatPay:
		return domain.PayChannelWechatPay, nil
	case huifu.ChannelAlipay:
		return domain.PayChannelAlipay, nil
	default:
		return "", fmt.Errorf("unknown trade type: %s", tradeType)
	}
}

var _ PaymentCallbackProvider = (*HuifuPaymentCallbackProvider)(nil)

type HuifuPaymentCallbackProvider struct{}

func NewHuifuPaymentCallbackProvider() *HuifuPaymentCallbackProvider {
	return &HuifuPaymentCallbackProvider{}
}

func (p *HuifuPaymentCallbackProvider) Callback(ctx context.Context, callback *domain.PaymentCallback, _ *domain.Payment) (res *PaymentCallbackResult, err error) {
	var respData map[string]any
	if err := json.Unmarshal(callback.Raw, &respData); err != nil {
		return nil, fmt.Errorf("json.Unmarshal respData: %w", err)
	}

	res = new(PaymentCallbackResult)
	// 交易终态
	if transStat, ok := respData["trans_stat"].(string); ok {
		res.State = huifuStateToPayState(huifu.TransStat(transStat))
	} else {
		res.State = domain.PayStateFailure
	}

	if res.State == domain.PayStateFailure {
		reason, _ := respData["bank_message"].(string)
		if reason == "" {
			reason, _ = getHuifuCallbackRespDesc(respData)
		}
		res.FailReason = reason
	}

	return
}

func getHuifuCallbackRespDesc(resp map[string]any) (string, error) {
	if desc, ok := resp["resp_desc"].(string); ok {
		return desc, nil
	} else if desc, ok = resp["sub_resp_desc"].(string); ok {
		return desc, nil
	}
	return "", errors.New("not found resp_desc")
}

func huifuStateToPayState(state huifu.TransStat) domain.PayState {
	switch state {
	case huifu.TransStatSuccess:
		return domain.PayStateSuccess
	case huifu.TransStatFail:
		return domain.PayStateFailure
	case huifu.TransStatProcessing:
		return domain.PayStateProcessing
	default:
		return domain.PayStateUnknown
	}
}
