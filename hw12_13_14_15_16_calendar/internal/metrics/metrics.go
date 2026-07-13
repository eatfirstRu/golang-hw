package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "calendar",
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "calendar",
			Name:      "http_request_duration_seconds",
			Help:      "Duration of HTTP requests in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	EventsCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "calendar",
			Name:      "events_created_total",
			Help:      "Total number of events created.",
		},
	)

	EventsUpdated = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "calendar",
			Name:      "events_updated_total",
			Help:      "Total number of events updated.",
		},
	)

	EventsDeleted = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "calendar",
			Name:      "events_deleted_total",
			Help:      "Total number of events deleted.",
		},
	)

	NotificationsSent = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "calendar",
			Name:      "notifications_sent_total",
			Help:      "Total number of notifications sent to Kafka.",
		},
	)

	NotificationsSendErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "calendar",
			Name:      "notifications_send_errors_total",
			Help:      "Total number of failed notification sends.",
		},
	)

	NotificationsSaved = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "calendar",
			Name:      "notifications_saved_total",
			Help:      "Total number of notifications saved to DB.",
		},
	)

	OldEventsDeleted = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "calendar",
			Name:      "old_events_deleted_total",
			Help:      "Total number of old events cleaned up.",
		},
	)

	SchedulerTickDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "calendar",
			Name:      "scheduler_tick_duration_seconds",
			Help:      "Duration of scheduler tick in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
	)

	SchedulerErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "calendar",
			Name:      "scheduler_errors_total",
			Help:      "Total number of scheduler tick errors.",
		},
	)
)
