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

func (i FeeInstance) ImageUsage() ([]ImageUsage, error) {
	var usage []ImageUsage
	if err := utils.Swap(i.data.TokenUsage, &usage); err != nil {
		return nil, err
	}

	return usage, nil
}

func (i FeeInstance) VideoUsage() (*VideoUsage, error) {
	var usage VideoUsage
	if err := utils.Swap(i.data.TokenUsage, &usage); err != nil {
		return nil, err
	}
	return &usage, nil
}

func (i FeeInstance) TextUsage() (*TextUsage, error) {
	var usage TextUsage
	if err := utils.Swap(i.data.TokenUsage, &usage); err != nil {
		return nil, err
	}
	return &usage, nil
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

	var textInstances, imageInstances, videoInstances []FeeInstance
	for _, usage := range report {
		inst := FeeInstance{userId: usage.UserId(), data: *usage}

		switch usage.ReportType {
		case ImageReportType:
			imageInstances = append(imageInstances, inst)
		case VideoReportType:
			videoInstances = append(videoInstances, inst)
		default:
			priceInfo, has := m.price.FetchProviderPrice(usage.ModelId)
			if !has {
				return false, fmt.Errorf("model price not found: %s, %s", usage.ModelId, usage.Model)
			}
			inst.priceInfo = priceInfo
			textInstances = append(textInstances, inst)
		}
		logrus.Infof("consume info: user: %s, provider: %s, model: %s, type: %s, usage: %v", usage.Caller, usage.Provider, usage.Model, usage.ReportType, usage.TokenUsage)
	}

	var allConsumes []*models.UserConsumeRecord
	if len(textInstances) > 0 {
		consumes, err := m.deductTextFees(textInstances)
		if err != nil {
			logrus.Errorf("Failed to deduct text fees: %v", err)
			return true, err
		}
		allConsumes = append(allConsumes, consumes...)
	}
	if len(imageInstances) > 0 {
		consumes, err := m.deductImageFees(imageInstances)
		if err != nil {
			logrus.Errorf("Failed to deduct image fees: %v", err)
			return true, err
		}
		allConsumes = append(allConsumes, consumes...)
	}
	if len(videoInstances) > 0 {
		consumes, err := m.deductVideoFees(videoInstances)
		if err != nil {
			logrus.Errorf("Failed to deduct video fees: %v", err)
			return true, err
		}
		allConsumes = append(allConsumes, consumes...)
	}

	if len(allConsumes) > 0 {
		m.mq.Publish(allConsumes)
	}
	return false, nil
}

func (m *FeeService) deductTextFees(instances []FeeInstance) ([]*models.UserConsumeRecord, error) {
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

		usage, _ := inst.TextUsage()
		inputCost := CalculateTokenCostMicro(usage.InputTokens, float64(inst.priceInfo.InputPrice))
		outputCost := CalculateTokenCostMicro(usage.InputTokens, float64(inst.priceInfo.InputPrice))

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
			ModelId:          inst.data.ModelId,
			NodeId:           inst.data.NodeId,
			TotalConsumed:    remainingCost,
			ActualProvider:   inst.data.ActualProvider,
			ActualProviderId: inst.data.ActualProviderId,
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

func (m *FeeService) deductImageFees(instances []FeeInstance) ([]*models.UserConsumeRecord, error) {
	session := m.xorm.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return nil, err
	}
	var consumes []*models.UserConsumeRecord
	for _, inst := range instances {
		balance := models.UserWallet{UserId: inst.userId}
		if has, err := session.Get(&balance); err != nil {
			return nil, err
		} else if !has {
			return nil, fmt.Errorf("user wallet not found: %d", inst.userId)
		}

		totalPrice, _, err := m.calculateImageActualCost(inst)
		if err != nil {
			return nil, err
		}
		balance.Balance -= totalPrice

		if _, err := session.ID(balance.Id).Update(&balance); err != nil {
			return nil, err
		}

		record := models.UserConsumeRecord{
			UserId:           inst.userId,
			Model:            inst.data.Model,
			ModelId:          inst.data.ModelId,
			NodeId:           inst.data.NodeId,
			TotalConsumed:    totalPrice,
			ConsumeType:      "image",
			ActualProvider:   inst.data.ActualProvider,
			ActualProviderId: inst.data.ActualProviderId,
			CreatedAt:        time.Now().Unix(),
		}
		if _, err := session.InsertOne(&record); err != nil {
			return nil, err
		}

		usage, _ := inst.ImageUsage()
		var details []models.UserConsumeDetailImage
		for _, img := range usage {
			details = append(details, models.UserConsumeDetailImage{
				ConsumeId: record.ID,
				Quality:   img.Quality,
				Size:      img.Size,
				CreatedAt: time.Now().Unix(),
			})
		}
		if len(details) > 0 {
			if _, err := session.InsertMulti(&details); err != nil {
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

func (m *FeeService) deductVideoFees(instances []FeeInstance) ([]*models.UserConsumeRecord, error) {
	session := m.xorm.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return nil, err
	}
	var consumes []*models.UserConsumeRecord
	for _, inst := range instances {
		balance := models.UserWallet{UserId: inst.userId}
		if has, err := session.Get(&balance); err != nil {
			return nil, err
		} else if !has {
			return nil, fmt.Errorf("user wallet not found: %d", inst.userId)
		}

		totalPrice, _, err := m.calculateVideoActualCost(inst)
		if err != nil {
			return nil, err
		}
		balance.Balance -= totalPrice

		if _, err := session.ID(balance.Id).Update(&balance); err != nil {
			return nil, err
		}

		record := models.UserConsumeRecord{
			UserId:           inst.userId,
			Model:            inst.data.Model,
			ModelId:          inst.data.ModelId,
			NodeId:           inst.data.NodeId,
			TotalConsumed:    totalPrice,
			ConsumeType:      "video",
			ActualProvider:   inst.data.ActualProvider,
			ActualProviderId: inst.data.ActualProviderId,
			CreatedAt:        time.Now().Unix(),
		}
		if _, err := session.InsertOne(&record); err != nil {
			return nil, err
		}

		usage, _ := inst.VideoUsage()
		detail := models.UserConsumeDetailVideo{
			ConsumdId: record.ID,
			Seconds:   usage.Seconds,
			Size:      usage.Size,
			CreatedAt: time.Now().Unix(),
		}
		if _, err := session.InsertOne(&detail); err != nil {
			return nil, err
		}
		consumes = append(consumes, &record)
	}
	if err := session.Commit(); err != nil {
		return nil, err
	}
	return consumes, nil
}
