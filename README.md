# old-iam-finder

Find old AWS IAM and send message to slack channel

## Getting Started

### Prerequisites

to run as k8s job, it needs inject environment variables from k8s secret.

```YAML
#secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: iam-finder-secret
  namespace: iam-finder
type: Opaque
stringData:
  AWS_ACCESS_KEY_ID: $(AWS_ACCESS_KEY_ID)
  AWS_SECRET_ACCESS_KEY: $(AWS_SECRET_ACCESS_KEY)
  SLACK_WEBHOOK_URL: $(SLACK_WEBHOOK_URL)
  EXPIRE_HOUR: $(EXPIRE_HOUR)
```

after complete all env variables, apply secrets to k8s

```bash
kubectl -f secret.yaml
```

### Run the `job`

and apply `job.yaml` to k8s, it runs once and sends message to slack channel

```bash
kubectl -f job.yaml
```
