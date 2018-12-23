// SPDX-License-Identifier: Apache-2.0

package systemd

import (
	"encoding/json"
	sd "github.com/coreos/go-systemd/dbus"
	"github.com/godbus/dbus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Unit struct {
	Action string `json:"action"`
	Unit string `json:"unit"`
	Property string `json:"property"`
	Value string `json:"value"`
}

type Property struct {
	Property string `json:"property"`
	Value string `json:"value"`
}

type UnitStatus struct {
	Status string `json:"property"`
	Unit string `json:"unit"`
}

func StartUnit(unit string) {
	conn, r := sd.NewSystemdConnection()
	if r != nil {
		log.Errorf("Failed to get systemd bus connection: %s", r)
		return
	}
	defer conn.Close()

	reschan := make(chan string)
	_, r = conn.StartUnit(unit, "replace", reschan)
	if r != nil {
		log.Errorf("Failed to start unit %s: %s",  unit)
		return
	}

	job := <-reschan
	if job != "done" {
		log.Errorf("Failed job %s: %s", job, unit)
	}
}

func StopUnit(unit string) {
	conn, r := sd.NewSystemdConnection()
	if r != nil {
		log.Errorf("Failed to get systemd bus connection: %s", r)
		return
	}
	defer conn.Close()

	reschan := make(chan string)
	_, r = conn.StopUnit(unit, "fail", reschan)
	if r != nil {
		log.Errorf("Failed to stop unit %s: %s", unit, r)
		return
	}

	job := <-reschan
	if job != "done" {
		log.Errorf("Failed job: %s", job, unit)
		return
	}
}

func RestartUnit(unit string) {
	conn, r := sd.NewSystemdConnection()
	if r != nil {
		log.Errorf("Failed to get systemd bus connection: %s", r)
		return
	}
	defer conn.Close()

	reschan := make(chan string)
	_, r = conn.RestartUnit(unit, "replace", reschan)
	if r != nil {
		log.Errorf("Failed to restart unit %s: %s",  unit, r)
		return
	}

	job := <-reschan
	if job != "done" {
		log.Errorf("Failed job %s: %s", job, unit)
	}
}

func ReloadUnit(unit string) {
	conn, r := sd.NewSystemdConnection()
	if r != nil {
		log.Errorf("Failed to get systemd bus connection: %s", r)
		return
	}
	defer conn.Close()

	r = conn.Reload()
	if r != nil {
		log.Errorf("Failed to reload unit %s: %s",  unit, r)
		return
	}
}

func GetUnitStatus(w http.ResponseWriter, unit string) {
	conn, r := sd.NewSystemdConnection()
	if r != nil {
		log.Errorf("Failed to get systemd bus connection: %s", r)
		return
	}
	defer conn.Close()

	units, r := conn.ListUnitsByNames([]string{unit})
	if r != nil {
		log.Errorf("Failed get unit %s status: %s", unit, r)
		return
	}

	status := UnitStatus{ Status: units[0].ActiveState, Unit: unit }
	json.NewEncoder(w).Encode(status)
}

func GetUnitProperty(w http.ResponseWriter, unit string, property string) {
	conn, r := sd.NewSystemdConnection()
	if r != nil {
		log.Errorf("Failed to get systemd bus connection: %s", r)
		return
	}
	defer conn.Close()

	prop, r := conn.GetUnitProperty(unit, property)
	if r != nil {
		log.Errorf("Failed to get unit %s property %s: %s", unit, property, r)
		return
	}

	unitprop := Property{ Property: prop.Name, Value: prop.Value.Value().(string) }
	json.NewEncoder(w).Encode(unitprop)
}

func SetUnitProperty(w http.ResponseWriter, unit string, property string, value string) {
	conn, r := sd.NewSystemdConnection()
	if r != nil {
		log.Errorf("Failed to get systemd bus connection: %s", r)
		return
	}
	defer conn.Close()

	switch property {
	case "CPUShares":
		n, r := strconv.ParseInt(value, 10, 64)
		if r != nil {
			log.Errorf("Failed to parse CPUShares: ", value, r)
			return
		}

		r = conn.SetUnitProperties(unit, true, sd.Property{"CPUShares", dbus.MakeVariant(uint64(n))})
		if r != nil {
			log.Errorf("Failed to set CPUShares %s: %s", value, r)
			return
		}
	}
}
