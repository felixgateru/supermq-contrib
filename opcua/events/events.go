// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package events

type createClientEvent struct {
	id          string
	opcuaNodeID string
}

type removeClientEvent struct {
	id string
}

type connectionEvent struct {
	channelIDs []string
	clientIDs  []string
	domainID   string
}

type createChannelEvent struct {
	channelID      string
	domainID       string
	opcuaServerURI string
}

type removeChannelEvent struct {
	channelID string
	domainID  string
}
