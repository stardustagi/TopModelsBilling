package services

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/deepissue/fee_server/config"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type Consumer interface {
	Do(LLMReportMessage) (bool, error)
}

type Handler struct {
	ctx      context.Context
	consumer Consumer
	stopChan chan struct{}
}

func (m *Handler) decode(data []byte) []*LLMCallData {
	var raw []*LLMCallData

	err := json.Unmarshal(data, &raw)
	if nil != err {
		logrus.Errorf("Failed to unmarshal data: %s", string(data))
		return nil
	}
	return raw
}

func (m *Handler) do(ch <-chan *nats.Msg) {
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-m.stopChan:
			return
		case msg, ok := <-ch:
			if !ok {
				continue
			}
			report := m.decode(msg.Data)
			skipped := report == nil || len(report) == 0
			// logrus.Debug("received a message: ", string(msg.Data), "skipped: ", skipped)
			if skipped {
				msg.Ack()
				continue
			}
			redeliver, err := m.consumer.Do(report)
			if nil == err {
				logrus.Infof("ack message: %v", report)
				msg.Ack()
			} else {
				if redeliver {
					msg.NakWithDelay(time.Minute * 5)
				} else {
					msg.Term()
				}
			}
		}
	}
}

type NatsMQ struct {
	ctx          context.Context
	cancel       context.CancelFunc
	config       *config.NatsMQConfig
	client       *nats.Conn
	handlers     map[string]*Handler
	subscription *nats.Subscription
	cacheChan    chan *nats.Msg
	mu           sync.RWMutex
}

func NewNatsMQ(ctx context.Context, config *config.NatsMQConfig) (*NatsMQ, error) {
	nc, err := nats.Connect(config.Url, nats.UserInfo(config.User, config.Pass))
	if nil != err {
		return nil, err
	}
	mq := &NatsMQ{
		ctx:       ctx,
		config:    config,
		client:    nc,
		handlers:  make(map[string]*Handler),
		cacheChan: make(chan *nats.Msg, config.BufferSize),
	}
	return mq, nil
}

func (m *NatsMQ) AddConsumer(name string, c Consumer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[name] = &Handler{ctx: m.ctx, consumer: c}
}

func (m *NatsMQ) RemoveConsumer(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[name].stopChan <- struct{}{}
	delete(m.handlers, name)
}

func (m *NatsMQ) Subscribe() error {
	if len(m.handlers) == 0 {
		return errors.New("No consumers registered")
	}

	js, err := m.client.JetStream()
	if err != nil {
		return err
	}

	sub, err := js.QueueSubscribe(m.config.Topic, m.config.WorkerGroup, func(msg *nats.Msg) {
		logrus.Debugf("DsstributeMessage: %s ", msg.Data)
		// 将消息分发给所有consumer
		m.distributeMessage(msg)
	},
		nats.Durable(m.config.Consumer), nats.MaxDeliver(5),
		nats.ManualAck(), nats.MaxAckPending(m.config.BufferSize),
		nats.AckWait(time.Minute*time.Duration(m.config.AckWaitMintues)),
	)

	if err != nil {
		logrus.Errorf("Failed to subscribe to topic %s: with queue: %s,  %v ", m.config.Topic, m.config.WorkerGroup, err)
		return err
	}
	logrus.Infof("Topic: %s subscribed, using consumer: %s, with queue: %s", m.config.Topic, m.config.Consumer, m.config.WorkerGroup)
	m.subscription = sub
	return nil
}

// distributeMessage 将消息分发给所有注册的consumer (FOUT模式)
func (m *NatsMQ) distributeMessage(msg *nats.Msg) {

	m.cacheChan <- msg
}

func (m *NatsMQ) Close() error {
	if m.subscription != nil {
		// m.subscription.Unsubscribe()
	}
	if m.client != nil {
		m.client.Close()
	}
	return nil
}

func (m *NatsMQ) Start() {
	for name, handler := range m.handlers {
		go handler.do(m.cacheChan)
		logrus.Infof("Consumer: %s started", name)
	}
}

func (m *NatsMQ) Publish(data interface{}) error {
	js, err := m.client.JetStream()
	if err != nil {
		logrus.Errorf("Failed to connect to JetStream: %v", err)
		return err
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// 设置发布选项
	opts := []nats.PubOpt{
		nats.AckWait(30 * time.Second),
	}
	_, err = js.Publish("billing.userConsume", payload, opts...)
	if err != nil {
		logrus.Errorf("Failed to publish message to topic %s: %v", "billing.userConsume", err)
		return err
	}

	logrus.Debugf("Published message to topic %s: %s", "billing.userConsume", string(payload))
	return nil
}
