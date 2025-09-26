xorm {
  datasource = ["topai:tjLjIVjsVbtqDQ4SEkUm@tcp(top-maas-prod-db.mysql.database.azure.com:3306)/top_maas?charset=utf8mb4&parseTime=true&loc=Local"]
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
