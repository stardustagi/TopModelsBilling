# create consumer
```bash
nats consumer add billing fee-worker-group --deliver all --ack explicit --filter "billing.nodeUsage"
```
