base {
  xorm {
    datasource = ["agentcp:agentcp@tcp(mysql)/modelunion_llm?charset=utf8mb4&parseTime=True&loc=Local"]
    show_sql = true
    driver = "mysql"
  }
  http {
    address = "0.0.0.0"
    port    = "8080"
    path    = ""
    key     = ""
  }
}
natsmq  {
  url       = "nats://nats-mq:4222"
  user      = "agentcp-mq"
  pass      = "mq09Hgyl871xMTblUiTBOLV3MKDeAy"
  topic     = "modelgate/token"
  consumer  = "modelgate-fee-consumer"
  worker_group      = "fee-worker-group"
  buffer_size       = 1024
  ack_wait_mintues  = 5
}
