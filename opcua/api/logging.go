// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

//go:build !test

package api

import (
	"context"
	"log/slog"
	"time"

	"github.com/absmach/supermq-contrib/opcua"
)

var _ opcua.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger *slog.Logger
	svc    opcua.Service
}

// LoggingMiddleware adds logging facilities to the core service.
func LoggingMiddleware(svc opcua.Service, logger *slog.Logger) opcua.Service {
	return &loggingMiddleware{
		logger: logger,
		svc:    svc,
	}
}

func (lm loggingMiddleware) CreateClient(ctx context.Context, clientID, opcuaNodeID string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("client_id", clientID),
			slog.String("node_id", opcuaNodeID),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Create client route-map failed", args...)
			return
		}
		lm.logger.Info("Create client route-map completed successfully", args...)
	}(time.Now())

	return lm.svc.CreateClient(ctx, clientID, opcuaNodeID)
}

func (lm loggingMiddleware) UpdateClient(ctx context.Context, clientID, opcuaNodeID string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("client_id", clientID),
			slog.String("node_id", opcuaNodeID),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update client route-map failed", args...)
			return
		}
		lm.logger.Info("Update client route-map completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdateClient(ctx, clientID, opcuaNodeID)
}

func (lm loggingMiddleware) RemoveClient(ctx context.Context, clientID string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("client_id", clientID),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Remove client route-map failed", args...)
			return
		}
		lm.logger.Info("Remove client route-map completed successfully", args...)
	}(time.Now())

	return lm.svc.RemoveClient(ctx, clientID)
}

func (lm loggingMiddleware) CreateChannel(ctx context.Context, channelID, domainID, opcuaServerURI string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("channel_id", channelID),
			slog.String("domain_id", domainID),
			slog.String("server_uri", opcuaServerURI),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Create channel route-map failed", args...)
			return
		}
		lm.logger.Info("Create channel route-map completed successfully", args...)
	}(time.Now())

	return lm.svc.CreateChannel(ctx, channelID, domainID, opcuaServerURI)
}

func (lm loggingMiddleware) UpdateChannel(ctx context.Context, channelID, domainID, opcuaServerURI string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("channel_id", channelID),
			slog.String("domain_id", domainID),
			slog.String("server_uri", opcuaServerURI),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Update channel route-map failed", args...)
			return
		}
		lm.logger.Info("Update channel route-map completed successfully", args...)
	}(time.Now())

	return lm.svc.UpdateChannel(ctx, channelID, domainID, opcuaServerURI)
}

func (lm loggingMiddleware) RemoveChannel(ctx context.Context, channelID, domainID string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("channel_id", channelID),
			slog.String("domain_id", domainID),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Remove channel route-map failed", args...)
			return
		}
		lm.logger.Info("Remove channel route-map completed successfully", args...)
	}(time.Now())

	return lm.svc.RemoveChannel(ctx, channelID, domainID)
}

func (lm loggingMiddleware) Connect(ctx context.Context, domainID string, channelIDs, clientIDs []string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("domain_id", domainID),
			slog.Any("channel_ids", channelIDs),
			slog.Any("client_ids", clientIDs),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Connect client to channel failed", args...)
			return
		}
		lm.logger.Info("Connect client to channel completed successfully", args...)
	}(time.Now())

	return lm.svc.Connect(ctx, domainID, channelIDs, clientIDs)
}

func (lm loggingMiddleware) Disconnect(ctx context.Context, channelIDs, clientIDs []string) (err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.Any("channel_ids", channelIDs),
			slog.Any("client_ids", clientIDs),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Disconnect client from channel failed", args...)
			return
		}
		lm.logger.Info("Disconnect client from channel completed successfully", args...)
	}(time.Now())

	return lm.svc.Disconnect(ctx, channelIDs, clientIDs)
}

func (lm loggingMiddleware) Browse(ctx context.Context, serverURI, namespace, identifier, identifierType string) (nodes []opcua.BrowsedNode, err error) {
	defer func(begin time.Time) {
		args := []any{
			slog.String("duration", time.Since(begin).String()),
			slog.String("server_uri", serverURI),
			slog.String("namespace", namespace),
			slog.String("identifier", identifier),
			slog.String("identifier_type", identifierType),
		}
		if err != nil {
			args = append(args, slog.Any("error", err))
			lm.logger.Warn("Browse available nodes failed", args...)
			return
		}
		lm.logger.Info("Browse available nodes completed successfully", args...)
	}(time.Now())

	return lm.svc.Browse(ctx, serverURI, namespace, identifier, identifierType)
}
