package services

import (
	"fmt"
	"time"

	"github.com/deepissue/core/server"
	"github.com/deepissue/core/utils"
	"github.com/deepissue/fee_server/config"
	"github.com/deepissue/fee_server/models"
	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

type FeeInstance struct {
	userId    int64
	data      LLMCallData
	priceInfo PriceInfo
}
type FeeService struct {
	xorm  xorm.EngineInterface
	mq    *NatsMQ
	price *PriceService
}

func NewFeeService(srv *server.Server, xorm xorm.EngineInterface, c *config.Config) (*FeeService, error) {
	f := &FeeService{
		xorm: xorm,
	}
	mq, err := NewNatsMQ(srv.Ctx, &c.Nats)
	if nil != err {
		return nil, err
	}
	f.mq = mq
	f.price = NewPriceService(srv.Ctx, xorm)
	return f, nil
}

func (m *FeeService) Start() error {

	m.mq.AddConsumer("fee", m)
	if err := m.mq.Subscribe(); err != nil {
		return err
	}
	m.mq.Start()

	return nil
}

func (m *FeeService) Stop() {
	m.mq.Close()
}

func (m *FeeService) Do(report LLMReportMessage) (bool, error) {

	logrus.Tracef("Received message: %v", report)
	var instances []FeeInstance
	for _, usage := range report {
		if usage.TokenUsage.ISZero() {
			continue
		}

		priceInfo, has := m.price.FetchProviderPrice(usage.ModelId)
		if !has {
			return false, fmt.Errorf("model price not found: %s, %s", usage.ModelId, usage.Model)
		}

		logrus.Infof("consume info: user: %s, provider: %s, model: %s, price: %v, usage: %s", usage.Caller, usage.Provider, usage.Model, priceInfo, usage.TokenUsage)
		instances = append(instances, FeeInstance{userId: usage.UserId(), data: *usage, priceInfo: priceInfo})
	}

	if len(instances) == 0 {
		return false, nil
	}

	consumes, err := m.deductFees(instances)
	if err != nil {
		logrus.Errorf("Failed to deduct fees: %v, error: %v", utils.EncodeToString(report), err)
		return true, err
	}

	m.mq.Publish(consumes)

	return false, nil
}

func (m *FeeService) deductFees(instances []FeeInstance) ([]*models.UserConsumeRecord, error) {
	session := m.xorm.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return nil, err
	}
	var consumes []*models.UserConsumeRecord
	for _, inst := range instances {
		var details []models.UserConsumeRecord
		balance := models.UserWallet{UserId: inst.userId}
		if has, err := session.Get(&balance); err != nil {
			return nil, err
		} else if !has {
			return nil, fmt.Errorf("user wallet not found: %d", inst.userId)
		}

		inputCost := CalculateTokenCostMicro(inst.data.TokenUsage.InputTokens, float64(inst.priceInfo.InputPrice))
		outputCost := CalculateTokenCostMicro(inst.data.TokenUsage.InputTokens, float64(inst.priceInfo.InputPrice))

		remainingCost := inputCost + outputCost
		balance.Balance -= remainingCost

		rows, err := session.ID(balance.Id).Update(&balance)
		if err != nil {
			logrus.Errorf("update user balance failed: %d, cost: %d", inst.userId, remainingCost)
			return nil, err
		}
		if rows == 0 {
			return nil, fmt.Errorf("failed to update user balance: %d, cost: %d", inst.userId, remainingCost)
		}

		//保存扣费记录
		record := models.UserConsumeRecord{
			UserId:           inst.userId,
			Model:            inst.data.Model,
			ModelId:          inst.data.NodeId,
			NodeId:           inst.data.NodeId,
			TotalConsumed:    remainingCost,
			ActualProvider:   inst.data.ActualProvider,
			ActualProviderId: inst.data.ActualProviderId,
			InputTokens:      inst.data.TokenUsage.InputTokens,
			OutputTokens:     inst.data.TokenUsage.OutputTokens,
			CacheTokens:      inst.data.TokenUsage.CacheTokens,
			InputPrice:       inst.priceInfo.InputPrice,
			OutputPrice:      inst.priceInfo.OutputPrice,
			CachePrice:       inst.priceInfo.CachePrice,
			CreatedAt:        time.Now().Unix(),
		}
		if _, err := session.InsertOne(&record); err != nil {
			logrus.Errorf("insert record: %v", err)
			return nil, err
		}
		if len(details) > 0 {
			if _, err := session.InsertMulti(&details); err != nil {
				logrus.Errorf("insert detail records: %v", err)
				return nil, err
			}
		}
		consumes = append(consumes, &record)
	}
	if err := session.Commit(); err != nil {
		return nil, err
	}

	return consumes, nil
}
