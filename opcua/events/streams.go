// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package events

import (
	"context"
	"errors"

	"github.com/absmach/supermq-contrib/opcua"
	"github.com/absmach/supermq/pkg/events"
)

const (
	keyType      = "opcua"
	keyNodeID    = "node_id"
	keyServerURI = "server_uri"

	clientPrefix = "client."
	clientCreate = clientPrefix + "create"
	clientUpdate = clientPrefix + "update"
	clientRemove = clientPrefix + "remove"

	channelPrefix     = "channel."
	channelCreate     = channelPrefix + "create"
	channelUpdate     = channelPrefix + "update"
	channelRemove     = channelPrefix + "remove"
	channelConnect    = channelPrefix + "connect"
	channelDisconnect = channelPrefix + "disconnect"
)

var (
	errMetadataType = errors.New("metadatada is not of type opcua")

	errMetadataFormat = errors.New("malformed metadata")

	errMetadataServerURI = errors.New("ServerURI not found in channel metadatada")

	errMetadataNodeID = errors.New("NodeID not found in client metadatada")
)

type eventHandler struct {
	svc opcua.Service
}

// NewEventHandler returns new event store handler.
func NewEventHandler(svc opcua.Service) events.EventHandler {
	return &eventHandler{
		svc: svc,
	}
}

func (es *eventHandler) Handle(ctx context.Context, event events.Event) error {
	msg, err := event.Encode()
	if err != nil {
		return err
	}

	switch msg["operation"] {
	case clientCreate:
		cte, e := decodeCreateClientEvent(msg)
		if e != nil {
			err = e
			break
		}
		err = es.svc.CreateClient(ctx, cte.id, cte.opcuaNodeID)
	case clientUpdate:
		ute, e := decodeCreateClientEvent(msg)
		if e != nil {
			err = e
			break
		}
		err = es.svc.CreateClient(ctx, ute.id, ute.opcuaNodeID)
	case clientRemove:
		rte := decodeRemoveClientEvent(msg)
		err = es.svc.RemoveClient(ctx, rte.id)
	case channelCreate:
		cce, e := decodeCreateChannelEvent(msg)
		if e != nil {
			err = e
			break
		}
		err = es.svc.CreateChannel(ctx, cce.channelID, cce.domainID, cce.opcuaServerURI)
	case channelUpdate:
		uce, e := decodeCreateChannelEvent(msg)
		if e != nil {
			err = e
			break
		}
		err = es.svc.CreateChannel(ctx, uce.channelID, uce.domainID, uce.opcuaServerURI)
	case channelRemove:
		rce := decodeRemoveChannelEvent(msg)
		err = es.svc.RemoveChannel(ctx, rce.channelID, rce.domainID)
	case channelConnect:
		rce := decodeConnectionEvent(msg)
		err = es.svc.Connect(ctx, rce.domainID, rce.channelIDs, rce.clientIDs)
	case channelDisconnect:
		rce := decodeConnectionEvent(msg)
		err = es.svc.Disconnect(ctx, rce.channelIDs, rce.clientIDs)
	}
	if err != nil && err != errMetadataType {
		return err
	}

	return nil
}

func decodeCreateClientEvent(event map[string]interface{}) (createClientEvent, error) {
	metadata := events.Read(event, "metadata", map[string]interface{}{})

	cte := createClientEvent{
		id: events.Read(event, "id", ""),
	}

	metadataOpcua, ok := metadata[keyType]
	if !ok {
		return createClientEvent{}, errMetadataType
	}

	metadataVal, ok := metadataOpcua.(map[string]interface{})
	if !ok {
		return createClientEvent{}, errMetadataFormat
	}

	val, ok := metadataVal[keyNodeID].(string)
	if !ok || val == "" {
		return createClientEvent{}, errMetadataNodeID
	}

	cte.opcuaNodeID = val
	return cte, nil
}

func decodeRemoveClientEvent(event map[string]interface{}) removeClientEvent {
	return removeClientEvent{
		id: events.Read(event, "id", ""),
	}
}

func decodeCreateChannelEvent(event map[string]interface{}) (createChannelEvent, error) {
	metadata := events.Read(event, "metadata", map[string]interface{}{})

	cce := createChannelEvent{
		channelID: events.Read(event, "id", ""),
		domainID:  events.Read(event, "domain", ""),
	}

	metadataOpcua, ok := metadata[keyType]
	if !ok {
		return createChannelEvent{}, errMetadataType
	}

	metadataVal, ok := metadataOpcua.(map[string]interface{})
	if !ok {
		return createChannelEvent{}, errMetadataFormat
	}

	val, ok := metadataVal[keyServerURI].(string)
	if !ok || val == "" {
		return createChannelEvent{}, errMetadataServerURI
	}

	cce.opcuaServerURI = val
	return cce, nil
}

func decodeRemoveChannelEvent(event map[string]interface{}) removeChannelEvent {
	return removeChannelEvent{
		channelID: events.Read(event, "id", ""),
		domainID:  events.Read(event, "domain", ""),
	}
}

func decodeConnectionEvent(event map[string]interface{}) connectionEvent {
	return connectionEvent{
		channelIDs: events.ReadStringSlice(event, "channel_ids"),
		clientIDs:  events.ReadStringSlice(event, "client_ids"),
		domainID:   events.Read(event, "domain", ""),
	}
}
