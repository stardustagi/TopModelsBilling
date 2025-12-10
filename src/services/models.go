package services

import (
	"fmt"
	"strconv"
)

type ReportType string

const (
	TextReportType  ReportType = "text"
	ImageReportType ReportType = "image"
	VideoReportType ReportType = "video"
)

// TokenUsage 记录 token 使用情况
type ImageUsage struct {
	Quality string `json:"quality"`
	Size    string `json:"size"`
}

// VideoUsage 记录 token 使用情况
type VideoUsage struct {
	Seconds float64 `json:"seconds"`
	Size    string  `json:"size"`
}

type TextUsage struct {
	InputTokens     int64   `json:"input_tokens"`
	OutputTokens    int64   `json:"output_tokens"`
	CacheTokens     int64   `json:"cache_tokens"`
	ReasoningTokens int     `json:"reasoning_tokens"`
	TokensPerSec    int     `json:"tokens_per_sec"`
	Latency         float64 `json:"latency"`
}

func (u TextUsage) ISZero() bool {
	return u.InputTokens+u.OutputTokens+u.CacheTokens == 0
}
func (u TextUsage) String() string {
	return fmt.Sprintf("<TokenUsage: input:%d, ouput:%d>", u.InputTokens, u.OutputTokens)
}

type LLMReportMessage []*LLMCallData
type LLMCallData struct {
	Id               string     `json:"id"`
	NodeId           int        `json:"node_id"`
	Model            string     `json:"model"`
	ModelId          int        `json:"model_id"`     // 模型id（计费使用）
	ActualModel      string     `json:"actual_model"` // 实际使用的模型
	Provider         string     `json:"provider"`
	ActualProvider   string     `json:"actual_provider"`    // 实际服务商
	ActualProviderId string     `json:"actual_provider_id"` // 实际服务商id
	Caller           string     `json:"caller"`
	CallerKey        string     `json:"caller_key"`
	ClientVersion    string     `json:"client_version,omitempty"`
	AgentVersion     string     `json:"agent_version,omitempty"`
	Stream           bool       `json:"stream"`
	ReportType       ReportType `json:"report_type"`
	TokenUsage       any        `json:"token_usage"`
}

func (l *LLMCallData) UserId() int64 {
	user, _ := strconv.ParseInt(l.Caller, 10, 64)
	return user
}

func (m LLMCallData) String() string {
	return fmt.Sprintf("<LLMCallData: id:%s, model:%s, caller:%s, node:%d>", m.Id, m.Model, m.Caller, m.NodeId)
}
