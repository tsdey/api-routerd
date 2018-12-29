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
	Action string   `json:"action"`
	Unit string     `json:"unit"`
	UnitType string `json:"unit_type"`
	Property string `json:"property"`
	Value string    `json:"value"`
}

type Property struct {
	Property string `json:"property"`
	Value string `json:"value"`
}

type UnitStatus struct {
	Status string `json:"property"`
	Unit string `json:"unit"`
}

func SystemState(w http.ResponseWriter) (error) {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	p, err := conn.SystemState()
	if err != nil {
		log.Errorf("Failed to get system state: %s",  err)
		return err
	}

	state := p.Value.Value().(string)
	prop := Property{ Property: p.Name, Value: state}

	j, err := json.Marshal(prop)
	if err != nil {
		log.Errorf("Failed to encode prop: %s", err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)

	return nil
}

func (u *Unit) StartUnit() (error) {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	reschan := make(chan string)
	_, err = conn.StartUnit(u.Unit, "replace", reschan)
	if err != nil {
		log.Errorf("Failed to start unit %s: %s", u.Unit)
		return err
	}

	return nil
}

func (u *Unit) StopUnit() (error) {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	reschan := make(chan string)
	_, err = conn.StopUnit(u.Unit, "fail", reschan)
	if err != nil {
		log.Errorf("Failed to stop unit %s: %s", u.Unit, err)
		return err
	}

	return nil
}

func (u *Unit) RestartUnit() (error) {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	reschan := make(chan string)
	_, err = conn.RestartUnit(u.Unit, "replace", reschan)
	if err != nil {
		log.Errorf("Failed to restart unit %s: %s", u.Unit, err)
		return err
	}

	return nil
}

func (u *Unit) ReloadUnit() (error) {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	err = conn.Reload()
	if err != nil {
		log.Errorf("Failed to reload unit %s: %s", u.Unit, err)
		return err
	}

	return nil
}

func (u *Unit) GetUnitStatus(w http.ResponseWriter) (error) {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	units, err := conn.ListUnitsByNames([]string{u.Unit})
	if err != nil {
		log.Errorf("Failed get unit %s status: %s", u.Unit, err)
		return err
	}

	status := UnitStatus{ Status: units[0].ActiveState, Unit: u.Unit }
	json.NewEncoder(w).Encode(status)

	return nil
}

func (u *Unit) GetUnitProperty(w http.ResponseWriter) (error) {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	p, err := conn.GetServiceProperty(u.Unit, u.Property)

	switch u.Property {
	case "CPUShares":
		cpu := strconv.FormatUint(p.Value.Value().(uint64), 10)
		prop := Property{ Property: p.Name, Value: cpu}

		j, err := json.Marshal(prop)
		if err != nil {
			log.Errorf("Failed to encode prop: %s", err)
			return err
		}

		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}

	return nil
}

func (u *Unit) SetUnitProperty(w http.ResponseWriter) (error) {
	conn, err := sd.NewSystemdConnection()
	if err != nil {
		log.Errorf("Failed to get systemd bus connection: %s", err)
		return err
	}
	defer conn.Close()

	switch u.Property {
	case "CPUShares":
		n, err := strconv.ParseInt(u.Value, 10, 64)
		if err != nil {
			log.Errorf("Failed to parse CPUShares: ", u.Value, err)
			return err
		}

		err = conn.SetUnitProperties(u.Unit, true, sd.Property{"CPUShares", dbus.MakeVariant(uint64(n))})
		if err != nil {
			log.Errorf("Failed to set CPUShares %s: %s", u.Value, err)
			return err
		}
	}

	return nil
}
