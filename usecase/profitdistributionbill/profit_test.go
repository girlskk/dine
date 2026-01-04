package profitdistributionbill

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/repository"
)

type ProfitTestSuite struct {
	UsecaseTestSuite
	interactor *ProfitDistributionBillInteractor
	ctx        context.Context
	ds         domain.DataStore
}

func TestProfitTestSuite(t *testing.T) {
	suite.Run(t, new(ProfitTestSuite))
}

func (s *ProfitTestSuite) SetupTest() {
	// 调用父类的 SetupTest 来初始化数据库
	s.UsecaseTestSuite.SetupTest()

	// 创建 DataStore 并赋值给 s.ds
	s.ds = repository.New(s.client)

	// 创建 ProductInteractor
	s.interactor = NewProfitDistributionBillInteractor(s.ds)
	s.ctx = context.Background()
}

func TestProfitDistributionBillInteractor_GenerateProfitDistributionBills(t *testing.T) {
	suite.Run(t, new(ProfitTestSuite))
}

func (s *ProfitTestSuite) TestProfitDistributionBillInteractor_GenerateProfitDistributionBills() {
	err := s.interactor.GenerateProfitDistributionBills(s.ctx)
	if err != nil {
		s.T().Errorf("generate profit distribution bills error: %v", err)
		return
	}
}
