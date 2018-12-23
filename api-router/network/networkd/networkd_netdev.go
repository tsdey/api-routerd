// SPDX-License-Identifier: Apache-2.0

package networkd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"restgateway/api-router/share"
	"strings"

	log "github.com/sirupsen/logrus"
)

// NetDev
type NetDev struct {
	Description string
	MACAddress  string
	MTUBytes    string
	Name        string
	Kind        string

	// Bond
	Mode               string
	TransmitHashPolicy string

	// Vlan
	VlanId string

	//Bridge
	HelloTimeSec    string
	ForwardDelaySec string
	AgeingTimeSec   string
}

func (netdev *NetDev) CreateNetDevSectionConfig() string {
	conf := "[NetDev]\n"

	if netdev.Name != "" {
		conf += "Name=" + strings.TrimSpace(netdev.Name) + "\n"
	}

	if netdev.Description != "" {
		conf += "Description=" + strings.TrimSpace(netdev.Description) + "\n"
	}

	if netdev.Kind != "" {
		conf += "Kind=" + strings.TrimSpace(netdev.Kind) + "\n"
	}

	if netdev.MACAddress != "" {
		conf += "MACAddress=" + strings.TrimSpace(netdev.MACAddress) + "\n"
	}

	if netdev.MTUBytes != "" {
		conf += "MTUBytes=" + strings.TrimSpace(netdev.MTUBytes) + "\n"
	}

	if netdev.Kind == "bond" {
		conf += "\n[Bond]\n"

		if netdev.Mode != "" {
			conf += "Mode=" + strings.TrimSpace(netdev.Mode) + "\n"
		}

		if netdev.TransmitHashPolicy != "" {
			conf += "TransmitHashPolicy=" + strings.TrimSpace(netdev.TransmitHashPolicy) + "\n"
		}
	}

	if netdev.Kind == "vlan" {
		conf += "\n[VLAN]\n"

		if netdev.VlanId != "" {
			conf += "Id=" + strings.TrimSpace(netdev.VlanId) + "\n"
		}
	}

	if netdev.Kind == "bridge" {
		conf += "\n[Bridge]\n"

		if netdev.HelloTimeSec != "" {
			conf += "HelloTimeSec=" + strings.TrimSpace(netdev.HelloTimeSec) + "\n"
		}

		if netdev.ForwardDelaySec != "" {
			conf += "ForwardDelaySec=" + strings.TrimSpace(netdev.ForwardDelaySec) + "\n"
		}

		if netdev.AgeingTimeSec != "" {
			conf += "AgeingTimeSec=" + strings.TrimSpace(netdev.AgeingTimeSec) + "\n"
		}
	}

	return conf
}

func NetdevdParseJsonFromHttpReq(req *http.Request) error {
	var configs map[string]interface{}

	body, r := ioutil.ReadAll(req.Body)
	if r != nil {
		log.Error("Failed to parse HTTP request: ", r)
		return r
	}

	json.Unmarshal([]byte(body), &configs)

	var netdev NetDev
	for key, value := range configs {
		switch key {
		case "MACAddress":
			netdev.MACAddress = value.(string)
			break
		case "Name":
			netdev.Name = value.(string)
			break
		case "MTUBytes":
			netdev.MTUBytes = value.(string)
			break
		case "Kind":
			netdev.Kind = value.(string)
			break
		case "Description":
			netdev.Description = value.(string)
			break
		case "Mode":
			netdev.Mode = value.(string)
			break
		case "TransmitHashPolicy":
			netdev.TransmitHashPolicy = value.(string)
			break
		case "VlanId":
			netdev.VlanId = value.(string)
			break
		case "HelloTimeSec":
			netdev.HelloTimeSec = value.(string)
			break
		case "ForwardDelaySec":
			netdev.ForwardDelaySec = value.(string)
			break
		case "AgeingTimeSec":
			netdev.AgeingTimeSec = value.(string)
			break
		}
	}

	netdevConfig := netdev.CreateNetDevSectionConfig()
	config := []string{netdevConfig}

	fmt.Println(config)

	unitName := fmt.Sprintf("25-%s.netdev", netdev.Name)
	unitPath := filepath.Join(NetworkdUnitPath, unitName)

	share.WriteFullFile(unitPath, config)

	return nil

}

func ConfigureNetDevFile(rw http.ResponseWriter, req *http.Request) {
	NetdevdParseJsonFromHttpReq(req)
}
