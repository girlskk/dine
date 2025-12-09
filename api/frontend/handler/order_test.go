package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/api/frontend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/domain/mock"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type OrderTestSuite struct {
	HandlerTestSuite
}

func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

func (s *OrderTestSuite) TestCreateOrder() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	mockOrderInteractor := mock.NewMockOrderInteractor(ctrl)
	mockTableInteractor := mock.NewMockTableInteractor(ctrl)
	h := NewOrderHandler(mockOrderInteractor, mockTableInteractor)

	tests := []struct {
		name         string
		setupMock    func()
		reqBody      string
		expectedResp *types.CreateOrderResp
		code         int
	}{
		{
			name: "参数错误 - 商品数量为负数",
			reqBody: `{
				"table_id": 1,
				"people_number": 7,
				"items": [
					{
						"product_id": 1,
						"price": "24.5",
						"quantity": "-1"
					}
				]
			}`,
			code: http.StatusBadRequest,
		},
		{
			name: "参数错误 - 商品数量为0",
			reqBody: `{
				"table_id": 1,
				"people_number": 7,
				"items": [
					{
						"product_id": 1,
						"price": "24.5",
						"quantity": "0"
					}
				]
			}`,
			code: http.StatusBadRequest,
		},
		{
			name: "参数错误 - 商品数量格式错误",
			reqBody: `{
				"table_id": 1,
				"people_number": 7,
				"items": [
					{
						"product_id": 1,
						"quantity": "xx",
						"price": "24.5"
					}
				]
			}`,
			code: http.StatusBadRequest,
		},
		{
			name: "参数错误 - 商品数量未提供",
			reqBody: `{
				"table_id": 1,
				"people_number": 7,
				"items": [
					{
						"product_id": 1,
						"price": "24.5"
					}
				]
			}`,
			code: http.StatusBadRequest,
		},
		{
			name: "参数错误 - 商品ID小于等于0",
			reqBody: `{
				"table_id": 1,
				"people_number": 7,
				"items": [
					{
						"product_id": 0,
						"quantity": "1",
						"price": "24.5"
					}
				]
			}`,
			code: http.StatusBadRequest,
		},
		{
			name: "参数错误 - 空商品列表",
			reqBody: `{
				"table_id": 1,
				"people_number": 7,
				"items": []
			}`,
			code: http.StatusBadRequest,
		},
		{
			name: "参数错误 - TableID小于等于0",
			reqBody: `{
				"table_id": 0,
				"people_number": 7,
				"items": [
					{
						"product_id": 1,
						"quantity": "1",
						"price": "24.5"
					}
				]
			}`,
			code: http.StatusBadRequest,
		},
		{
			name: "业务错误",
			reqBody: `{
				"table_id": 1,
				"people_number": 7,
				"items": [
					{
						"product_id": 1,
						"quantity": "1",
						"price": "24.5"
					}
				]
			}`,
			setupMock: func() {
				mockTableInteractor.EXPECT().
					Get(gomock.Any(), 1).
					Return(&domain.Table{ID: 1, StoreID: 1}, nil)

				mockOrderInteractor.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any()).
					Return(nil, domain.ParamsError(domain.ErrProductNotExists))
			},
			code: http.StatusBadRequest,
		},
		{
			name: "创建订单成功",
			reqBody: `{
				"table_id": 1,
				"people_number": 7,
				"items": [
					{
						"product_id": 1,
						"quantity": "2",
						"price": "24.5"
					}
				]
			}`,
			setupMock: func() {
				mockTableInteractor.EXPECT().
					Get(gomock.Any(), 1).
					Return(&domain.Table{ID: 1, StoreID: 1}, nil)

				mockOrderInteractor.EXPECT().
					CreateOrder(gomock.Any(), &domain.CreateOrderParams{
						Table:   &domain.Table{ID: 1, StoreID: 1},
						Creator: &domain.FrontendUser{ID: 1},
						Store:   &domain.Store{ID: 1},
						Items: []*domain.CreateOrderItem{
							{
								ProductID: 1,
								Quantity:  decimal.NewFromFloat(2),
								Price:     decimal.NewFromFloat(24.5),
							},
						},
						PeopleNumber: 7,
					}).
					Return(&domain.Order{No: "O202503061742"}, nil)
			},
			expectedResp: &types.CreateOrderResp{
				No: "O202503061742",
			},
		},
	}

	s.r.POST("/order/create", func(c *gin.Context) {
		store := &domain.Store{ID: 1}
		user := &domain.FrontendUser{ID: 1}
		ctx := c.Request.Context()
		ctx = domain.NewStoreContext(ctx, store)
		ctx = domain.NewFrontendUserContext(ctx, user)
		c.Request = c.Request.Clone(ctx)

		h.CreateOrder()(c)
	})

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/order/create", strings.NewReader(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")

			if tt.setupMock != nil {
				tt.setupMock()
			}

			s.r.ServeHTTP(w, req)

			require.Equal(t, http.StatusOK, w.Code)

			type respType struct {
				response.Response
				Data *types.CreateOrderResp `json:"data"`
			}
			var resp respType
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			assert.Equal(t, tt.code, resp.Code)
			if tt.expectedResp != nil {
				assert.Equal(t, tt.expectedResp, resp.Data)
			}
		})
	}
}
