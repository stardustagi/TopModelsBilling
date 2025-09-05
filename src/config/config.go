package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/sirupsen/logrus"
)

// NatsMQConfig
// url       = "nats://47.128.253.184:4222"
// user      = "agentcp-mq"
// pass      = ""
// topic     = "modelgate"
type NatsMQConfig struct {
	Url            string `json:"url" hcl:"url"`
	User           string `json:"user" hcl:"user"`
	Pass           string `json:"pass" hcl:"pass"`
	Topic          string `json:"topic" hcl:"topic"`
	Consumer       string `json:"consumer" hcl:"consumer"`
	BufferSize     int    `json:"buffer_size" hcl:"buffer_size"`
	WorkerGroup    string `json:"worker_group" hcl:"worker_group"`
	AckWaitMintues int    `json:"ack_wait_mintues" hcl:"ack_wait_mintues"`
}

type XormConfig struct {
	ShowSql    string   `json:"show_sql" hcl:"show_sql"`
	Datasource []string `json:"datasource" hcl:"datasource"`
	Driver     string   `json:"driver" hcl:"driver"`
}

type Config struct {
	Nats NatsMQConfig `json:"natsmq" hcl:"natsmq,block"`
	Xorm XormConfig   `json:"xorm" hcl:"xorm,block"`
}

func LoadConfig(configPath string) *Config {

	c := &Config{}

	parser := hclparse.NewParser()
	file, err := parser.ParseHCLFile(configPath)
	if err != nil {
		// 处理错误
		logrus.Error(err)
	}
	ctx := &hcl.EvalContext{}
	diags := gohcl.DecodeBody(file.Body, ctx, c)
	if diags.HasErrors() {
		// 处理错误
		// logrus.Error(diags.Error())
	}
	return c
}
