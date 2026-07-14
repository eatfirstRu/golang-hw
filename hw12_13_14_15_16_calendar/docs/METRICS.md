# Calendar Service Metrics

## Endpoint

`GET /metrics` — returns Prometheus-compatible metrics.

## HTTP Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `calendar_http_requests_total` | Counter | method, path, status | Total HTTP requests count including errors. Useful for tracking traffic volume and error rates. |
| `calendar_http_request_duration_seconds` | Histogram | method, path | Request latency distribution. Helps identify slow endpoints and performance degradation. |

## Business Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `calendar_events_created_total` | Counter | Total events created. Tracks user activity. |
| `calendar_events_updated_total` | Counter | Total events updated. |
| `calendar_events_deleted_total` | Counter | Total events deleted. |

## Scheduler Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `calendar_notifications_sent_total` | Counter | Notifications successfully sent to Kafka. |
| `calendar_notifications_send_errors_total` | Counter | Failed notification sends. Alerts on Kafka issues. |
| `calendar_old_events_deleted_total` | Counter | Old events cleaned up by scheduler. |
| `calendar_scheduler_tick_duration_seconds` | Histogram | Scheduler tick duration. Detects slow DB queries or Kafka latency. |
| `calendar_scheduler_errors_total` | Counter | Scheduler tick errors. Indicates infrastructure problems. |

## Storer Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `calendar_notifications_saved_total` | Counter | Notifications saved to DB by storer. Verifies end-to-end pipeline health. |

## Usage Examples

**Error rate:**
```promql
rate(calendar_http_requests_total{status=~"5.."}[5m])
  / rate(calendar_http_requests_total[5m])
```

**P99 latency:**
```promql
histogram_quantile(0.99, rate(calendar_http_request_duration_seconds_bucket[5m]))
```

**Notification pipeline health:**
```promql
rate(calendar_notifications_sent_total[5m])
rate(calendar_notifications_saved_total[5m])
```

## Infrastructure

- Prometheus scrapes `calendar:8080/metrics` every 5s
- Prometheus UI available at `http://localhost:9090` after `make up`
