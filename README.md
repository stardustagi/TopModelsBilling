# create consumer
```bash
nats consumer add billing stat-worker-group --deliver all --ack explicit --filter "billing.userConsume"
```
