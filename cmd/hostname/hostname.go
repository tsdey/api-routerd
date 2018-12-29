// SPDX-License-Identifier: Apache-2.0

package hostname

import (
	"api-routerd/cmd/share"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"fmt"
)

var HostNameInfo = map[string]string{
	"Hostname"                  : "",
	"StaticHostname"            : "",
	"PrettyHostname"            : "",
	"IconName"                  : "",
	"Chassis"                   : "",
	"Deployment"                : "",
	"Location"                  : "",
	"KernelName"                : "",
	"KernelRelease"             : "",
	"KernelVersion"             : "",
	"OperatingSystemPrettyName" : "",
	"OperatingSystemCPEName"    : "",
	"HomeURL"                   : "",
}

var HostMethodInfo = map[string]string{
	"SetHostname"       : "",
	"SetStaticHostname" : "",
	"SetPrettyHostname" : "",
	"SetIconName"       : "",
	"SetChassis"        : "",
	"SetDeployment"     : "",
	"SetLocation"       : "",
}

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

	_, k := HostMethodInfo[hostname.Property]
	if !k {
		return fmt.Errorf("Failed to set hostname property: %s not found", k)
	}

	h := conn.Object("org.freedesktop.hostname1", "/org/freedesktop/hostname1")
	r := h.Call("org.freedesktop.hostname1." + hostname.Property, 0, hostname.Value, false).Err
	if r != nil {
		log.Errorf("Failed to set hostname: ", r)
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


	h := conn.Object("org.freedesktop.hostname1", "/org/freedesktop/hostname1")
	for k, _ := range HostNameInfo {
		p, perr := h.GetProperty("org.freedesktop.hostname1." + k)
		if perr != nil {
			log.Errorf("Failed to get org.freedesktop.hostname1.%s", k)
			continue
		}

		hv, b := p.Value().(string)
		if !b {
			log.Error("Received unexpected type as value, expected string got :", property , hv)
			continue
		}

		HostNameInfo[k] = hv
	}

	if property == "" {
		b, err := json.Marshal(HostNameInfo)
		if err != nil {
			return err
		}

		rw.Write(b)
	} else {
		host := Hostname {Property: property, Value: HostNameInfo[property]}
		b, err := json.Marshal(host)
		if err != nil {
			return err
		}

		rw.Write(b)
	}

	return nil
}
