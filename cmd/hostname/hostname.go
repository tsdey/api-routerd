// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"api-routerd/cmd/share"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"net/http"
)

type Hostname struct {
	Property string `json:"property"`
	Value string    `json:"value"`
}

func (hostname *Hostname) SetHostname() (error) {
	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		log.Error("Failed to get systemd bus connection: ", err)
		return err
	}
	defer conn.Close()

	if (len(strings.TrimSpace(hostname.Value)) == 0) {
		log.Error("Empty hostname received")
		return errors.New("Empty hostname received")
	}

	property := "SetStaticHostname"

	switch(hostname.Property) {
	case "pretty":
		property = "PrettyHostname"
		break
	case "transient":
		property = "SetHostname"
		break
	case "static":
		property = "SetStaticHostname"
		break
	}

	h := conn.Object("org.freedesktop.hostname1", "/org/freedesktop/hostname1")
	errDbus := h.Call("org.freedesktop.hostname1." + property, 0, hostname.Value, false).Err
	if errDbus != nil {
		log.Errorf("Failed to set hostname: ", errDbus)
		return errors.New("Failed to set hostname")
	}

	return nil
}

func GetHostname(rw http.ResponseWriter, property string) (error) {
	conn, err := share.GetSystemBusPrivateConn()
	if err != nil {
		log.Error("Failed to get dbus connection: ", err)
		return err
	}
	defer conn.Close()

	switch(property) {
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
	p, perr := h.GetProperty("org.freedesktop.hostname1." + property)
	if perr != nil {
		log.Error("Failed to get org.freedesktop.hostname1.%s", property)
		return perr
	}

	if p.Value() == nil {
		log.Error("Empty value received when reading property : ", property)
		return errors.New("Invalid Value")
	}

	v, be := p.Value().(string)
	if !be {
		log.Error("Received unexpected type as value, expected string got :", property , v)
		return errors.New("Invalid value")
	}

	host := Hostname {Property: property, Value: v}

	b, err := json.Marshal(host)
	if err != nil {
		return err
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(b)

	return nil
}
