// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package opcua

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"

	"github.com/absmach/supermq-contrib/opcua/db"
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// CreateClient creates clientID:OPC-UA-nodeID route-map
	CreateClient(ctx context.Context, clientID, nodeID string) error

	// UpdateClient updates clientID:OPC-UA-nodeID route-map
	UpdateClient(ctx context.Context, clientID, nodeID string) error

	// RemoveClient removes clientID:OPC-UA-nodeID route-map
	RemoveClient(ctx context.Context, clientID string) error

	// CreateChannel creates channelID:OPC-UA-serverURI route-map
	CreateChannel(ctx context.Context, channelID, domainID, serverURI string) error

	// UpdateChannel updates channelID:OPC-UA-serverURI route-map
	UpdateChannel(ctx context.Context, channelID, domainID, serverURI string) error

	// RemoveChannel removes channelID:OPC-UA-serverURI route-map
	RemoveChannel(ctx context.Context, channelID, domainID string) error

	// Connect creates clientID:channelID route-map
	Connect(ctx context.Context, domainID string, channelIDs, clientIDs []string) error

	// Disconnect removes clientID:channelID route-map
	Disconnect(ctx context.Context, channelIDs, clientIDs []string) error

	// Browse browses available nodes for a given OPC-UA Server URI and NodeID
	Browse(ctx context.Context, serverURI, namespace, identifier, identifierType string) ([]BrowsedNode, error)
}

// Config OPC-UA Server.
type Config struct {
	ServerURI string
	NodeID    string
	Interval  string `env:"SMQ_OPCUA_ADAPTER_INTERVAL_MS"     envDefault:"1000"`
	Policy    string `env:"SMQ_OPCUA_ADAPTER_POLICY"          envDefault:""`
	Mode      string `env:"SMQ_OPCUA_ADAPTER_MODE"            envDefault:""`
	CertFile  string `env:"SMQ_OPCUA_ADAPTER_CERT_FILE"       envDefault:""`
	KeyFile   string `env:"SMQ_OPCUA_ADAPTER_KEY_FILE"        envDefault:""`
}

var (
	_         Service = (*adapterService)(nil)
	guidRegex         = regexp.MustCompile(`^\{?[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}\}?$`)
)

type adapterService struct {
	subscriber Subscriber
	browser    Browser
	clientsRM  RouteMapRepository
	channelsRM RouteMapRepository
	connectRM  RouteMapRepository
	cfg        Config
	logger     *slog.Logger
}

// New instantiates the OPC-UA adapter implementation.
func New(sub Subscriber, brow Browser, clientsRM, channelsRM, connectRM RouteMapRepository, cfg Config, log *slog.Logger) Service {
	return &adapterService{
		subscriber: sub,
		browser:    brow,
		clientsRM:  clientsRM,
		channelsRM: channelsRM,
		connectRM:  connectRM,
		cfg:        cfg,
		logger:     log,
	}
}

func (as *adapterService) CreateClient(ctx context.Context, clientID, nodeID string) error {
	return as.clientsRM.Save(ctx, clientID, nodeID)
}

func (as *adapterService) UpdateClient(ctx context.Context, clientID, nodeID string) error {
	return as.clientsRM.Save(ctx, clientID, nodeID)
}

func (as *adapterService) RemoveClient(ctx context.Context, clientID string) error {
	return as.clientsRM.Remove(ctx, clientID)
}

func (as *adapterService) CreateChannel(ctx context.Context, channelID, domainID, serverURI string) error {
	val := fmt.Sprintf("%s:%s", channelID, domainID)
	return as.channelsRM.Save(ctx, val, serverURI)
}

func (as *adapterService) UpdateChannel(ctx context.Context, channelID, domainID, serverURI string) error {
	val := fmt.Sprintf("%s:%s", channelID, domainID)
	return as.channelsRM.Save(ctx, val, serverURI)
}

func (as *adapterService) RemoveChannel(ctx context.Context, channelID, domainID string) error {
	val := fmt.Sprintf("%s:%s", channelID, domainID)
	return as.channelsRM.Remove(ctx, val)
}

func (as *adapterService) Connect(ctx context.Context, domainID string, channelIDs, clientIDs []string) error {
	for _, chanID := range channelIDs {
		val := fmt.Sprintf("%s:%s", chanID, domainID)
		serverURI, err := as.channelsRM.Get(ctx, val)
		if err != nil {
			return err
		}

		for _, clientID := range clientIDs {
			nodeID, err := as.clientsRM.Get(ctx, clientID)
			if err != nil {
				return err
			}

			cfg := as.cfg
			cfg.NodeID = nodeID
			cfg.ServerURI = serverURI

			c := fmt.Sprintf("%s:%s", chanID, clientID)
			if err := as.connectRM.Save(ctx, c, c); err != nil {
				return err
			}

			go func(cfg Config) {
				if err := as.subscriber.Subscribe(context.Background(), cfg); err != nil {
					as.logger.Warn("subscription failed", slog.Any("error", err))
				}
			}(cfg)

			// Store subscription details
			if err := db.Save(serverURI, nodeID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (as *adapterService) Browse(ctx context.Context, serverURI, namespace, identifier, identifierType string) ([]BrowsedNode, error) {
	idFormat := "s"
	switch identifierType {
	case "string":
		break
	case "numeric":
		if _, err := strconv.Atoi(identifier); err != nil {
			args := []any{
				slog.String("namespace", namespace),
				slog.String("identifier", identifier),
				slog.Any("error", err),
			}
			as.logger.Warn("failed to parse numeric identifier", args...)
			break
		}
		idFormat = "i"
	case "guid":
		if !guidRegex.MatchString(identifier) {
			args := []any{
				slog.String("namespace", namespace),
				slog.String("identifier", identifier),
			}
			as.logger.Warn("GUID identifier has invalid format", args...)
			break
		}
		idFormat = "g"
	case "opaque":
		if _, err := base64.StdEncoding.DecodeString(identifier); err != nil {
			args := []any{
				slog.String("namespace", namespace),
				slog.String("identifier", identifier),
				slog.Any("error", err),
			}
			as.logger.Warn("opaque identifier has invalid base64 format", args...)
			break
		}
		idFormat = "b"
	}
	nodeID := fmt.Sprintf("ns=%s;%s=%s", namespace, idFormat, identifier)
	nodes, err := as.browser.Browse(serverURI, nodeID)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func (as *adapterService) Disconnect(ctx context.Context, channelIDs, clientIDs []string) error {
	for _, channelID := range channelIDs {
		for _, clientID := range clientIDs {
			c := fmt.Sprintf("%s:%s", channelID, clientID)
			if err := as.connectRM.Remove(ctx, c); err != nil {
				return err
			}
		}
	}
	return nil
}
