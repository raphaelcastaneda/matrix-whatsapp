// mautrix-whatsapp - A Matrix-WhatsApp puppeting bridge.
// Copyright (C) 2018 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package whatsapp_ext

import (
	"encoding/json"
	"strings"

	"github.com/Rhymen/go-whatsapp"
)

type PresenceType string

const (
	PresenceUnavailable PresenceType = "unavailable"
	PresenceAvailable   PresenceType = "available"
	PresenceComposing   PresenceType = "composing"
)

type Presence struct {
	JID       string       `json:"id"`
	SenderJID string       `json:"participant"`
	Status    PresenceType `json:"type"`
	Timestamp int64        `json:"t"`
	Deny      bool         `json:"deny"`
}

type PresenceHandler interface {
	whatsapp.Handler
	HandlePresence(Presence)
}

func (ext *ExtendedConn) handleMessagePresence(message []byte) {
	var event Presence
	err := json.Unmarshal(message, &event)
	if err != nil {
		ext.jsonParseError(err)
		return
	}
	event.JID = strings.Replace(event.JID, OldUserSuffix, NewUserSuffix, 1)
	if len(event.SenderJID) == 0 {
		event.SenderJID = event.JID
	} else {
		event.SenderJID = strings.Replace(event.SenderJID, OldUserSuffix, NewUserSuffix, 1)
	}
	for _, handler := range ext.handlers {
		presenceHandler, ok := handler.(PresenceHandler)
		if !ok {
			continue
		}
		go presenceHandler.HandlePresence(event)
	}
}