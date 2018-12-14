// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"restgateway/api-router/share"
	"strings"
	"net/http"
)

type Hostname struct {
	Action string   `json:"action"`
	Method string   `json:"method"`
	Property string `json:"property"`
	Value string    `json:"value"`
}

func SetHostname(hostname *Hostname) {
	conn, r := share.GetSystemBusPrivateConn()
	if r != nil {
		log.Error("Failed to get systemd bus connection: ", r)
		return
	}
	defer conn.Close()

	if (len(strings.TrimSpace(hostname.Value)) == 0) {
		log.Error("invalid parameter")
		return
	}

	method := "SetStaticHostname"

	switch(hostname.Method) {
	case "pretty":
		method = "PrettyHostname"
		break
	case "transient":
		method = "SetHostname"
		break
	case "static":
		method = "SetStaticHostname"
		break
	}

	h := conn.Object("org.freedesktop.hostname1", "/org/freedesktop/hostname1")
	err := h.Call("org.freedesktop.hostname1." + method, 0, hostname.Value, false)
	if err != nil {
		log.Error("failed to set hostname:")
		return
	}
}

func GetHostname(rw http.ResponseWriter, hostname *Hostname) {
	conn, r := share.GetSystemBusPrivateConn()
	if r != nil {
		log.Error("Failed to get dbus connection: ", r)
		return
	}
	defer conn.Close()

	prop := strings.TrimSpace(hostname.Property)
	if (len(prop) == 0) {
		log.Error("invalid parameter")
		return
	}

	property := "static"

	switch(prop) {
	case "pretty":
		property = "PrettyHostname"
		break
	case "transient":
		property = "Hostname"
		break
	case "static":
		property = "StaticHostname"
		break
	}

	h := conn.Object("org.freedesktop.hostname1", "/org/freedesktop/hostname1")
	p, pe := h.GetProperty("org.freedesktop.hostname1." + property)
	if pe != nil {
		log.Error("Failed to get org.freedesktop.hostname1.%s", property)
		return
	}

	if p.Value() == nil {
		log.Error("empty value received when reading property : ", property)
		return
	}

	v, be := p.Value().(string)
	if !be {
		log.Error("Received unexpected type as value, expected string got :", property , v)
		return
	}

	host := Hostname {Action: "get-hostname", Property: property, Value: v}
	json.NewEncoder(rw).Encode(host)
}
