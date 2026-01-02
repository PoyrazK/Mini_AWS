package platform

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	WSConnectionsActive = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mini_aws_ws_connections_active",
		Help: "The total number of active WebSocket connections",
	})

	// Auto-Scaling metrics
	AutoScalingEvaluations = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mini_aws_autoscaling_evaluations_total",
		Help: "Total number of auto-scaling evaluation cycles",
	})
	AutoScalingScaleOutEvents = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mini_aws_autoscaling_scale_out_total",
		Help: "Total number of scale-out events",
	})
	AutoScalingScaleInEvents = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mini_aws_autoscaling_scale_in_total",
		Help: "Total number of scale-in events",
	})
	AutoScalingCurrentInstances = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mini_aws_autoscaling_current_instances",
		Help: "Current instance count per scaling group",
	}, []string{"scaling_group_id"})

	// Load Balancer metrics
	LBRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mini_aws_lb_requests_total",
		Help: "Total requests proxied by load balancers",
	}, []string{"lb_id"})
)
