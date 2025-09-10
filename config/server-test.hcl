xorm {
  datasource = ["root:123456@tcp(127.0.0.1)/top_maas?charset=utf8mb4&parseTime=True&loc=Local"]
  show_sql = true
  driver = "mysql"
}

natsmq  {
  url       = "nats://127.0.0.1:4222"
  user      = ""
  pass      = ""
  topic     = "billing.nodeUsage"
  consumer  = "fee-consumer"
  worker_group      = "fee-worker-group"
  buffer_size       = 1024
  ack_wait_mintues  = 5
}
