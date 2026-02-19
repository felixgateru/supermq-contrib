// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"time"

	"github.com/absmach/supermq-contrib/opcua"
	"github.com/go-kit/kit/metrics"
)

var _ opcua.Service = (*metricsMiddleware)(nil)

type metricsMiddleware struct {
	counter metrics.Counter
	latency metrics.Histogram
	svc     opcua.Service
}

// MetricsMiddleware instruments core service by tracking request count and latency.
func MetricsMiddleware(svc opcua.Service, counter metrics.Counter, latency metrics.Histogram) opcua.Service {
	return &metricsMiddleware{
		counter: counter,
		latency: latency,
		svc:     svc,
	}
}

func (mm *metricsMiddleware) CreateClient(ctx context.Context, clientID, opcuaNodeID string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "create_client").Add(1)
		mm.latency.With("method", "create_client").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.CreateClient(ctx, clientID, opcuaNodeID)
}

func (mm *metricsMiddleware) UpdateClient(ctx context.Context, clientID, opcuaNodeID string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "update_client").Add(1)
		mm.latency.With("method", "update_client").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.UpdateClient(ctx, clientID, opcuaNodeID)
}

func (mm *metricsMiddleware) RemoveClient(ctx context.Context, clientID string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "remove_client").Add(1)
		mm.latency.With("method", "remove_client").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.RemoveClient(ctx, clientID)
}

func (mm *metricsMiddleware) CreateChannel(ctx context.Context, channelID, domainID, opcuaServerURI string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "create_channel").Add(1)
		mm.latency.With("method", "create_channel").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.CreateChannel(ctx, channelID, domainID, opcuaServerURI)
}

func (mm *metricsMiddleware) UpdateChannel(ctx context.Context, channelID, domainID, opcuaServerURI string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "update_channel").Add(1)
		mm.latency.With("method", "update_channel").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.UpdateChannel(ctx, channelID, domainID, opcuaServerURI)
}

func (mm *metricsMiddleware) RemoveChannel(ctx context.Context, channelID, domainID string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "remove_channel").Add(1)
		mm.latency.With("method", "remove_channel").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.RemoveChannel(ctx, channelID, domainID)
}

func (mm *metricsMiddleware) Connect(ctx context.Context, domainID string, channelIDs, clientIDs []string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "connect").Add(1)
		mm.latency.With("method", "connect").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Connect(ctx, domainID, channelIDs, clientIDs)
}

func (mm *metricsMiddleware) Disconnect(ctx context.Context, channelIDs, clientIDs []string) error {
	defer func(begin time.Time) {
		mm.counter.With("method", "disconnect").Add(1)
		mm.latency.With("method", "disconnect").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Disconnect(ctx, channelIDs, clientIDs)
}

func (mm *metricsMiddleware) Browse(ctx context.Context, serverURI, namespace, identifier, identifierType string) ([]opcua.BrowsedNode, error) {
	defer func(begin time.Time) {
		mm.counter.With("method", "browse").Add(1)
		mm.latency.With("method", "browse").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mm.svc.Browse(ctx, serverURI, namespace, identifier, identifierType)
}
