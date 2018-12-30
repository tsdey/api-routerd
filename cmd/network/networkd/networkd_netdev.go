// SPDX-License-Identifier: Apache-2.0

package networkd

import (
	"api-routerd/cmd/share"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

// NetDev
type NetDev struct {
	Description         string `json:"Description,omitempty"`
	MACAddress          string `json:"MACAddress,omitempty"`
	MTUBytes            string `json:"MTUBytes,omitempty"`
	Name                string `json:"Name,omitempty"`
	Kind                string `json:"Kind,omitempty"`

	// Bond
	Mode                string `json:"Mode,omitempty"`
	TransmitHashPolicy  string `json:"TransmitHashPolicy,omitempty"`

	// Vlan
	VlanId              string `json:"VlanId,omitempty"`

	//Bridge
	HelloTimeSec        string `json:"HelloTimeSec,omitempty"`
	ForwardDelaySec     string `json:"ForwardDelaySec,omitempty"`
	AgeingTimeSec       string `json:"AgeingTimeSec,omitempty"`

	//Tunnel
	Local               string `json:"Local,omitempty"`
	Remote              string `json:"Remote,omitempty"`
	TTL                 string `json:"TTL,omitempty"`
	DiscoverPathMTU     string `json:"DiscoverPathMTU,omitempty"`
	IPv6FlowLabel       string `json:"IPv6FlowLabel,omitempty"`
	EncapsulationLimit  string `json:"EncapsulationLimit,omitempty"`
	Key                 string `json:"Key,omitempty"`
	Independent         string `json:"Independent,omitempty"`

	//VxLan
	Id                  string `json:"Id,omitempty"`
	TOS                 string `json:"TOS,omitempty"`
	MacLearning         string `json:"MacLearning,omitempty"`
	DestinationPort     string `json:"DestinationPort,omitempty"`
	PortRange           string `json:"PortRange,omitempty"`
	FlowLabel           string `json:"FlowLabel,omitempty"`

	//Veth
	Peer               string `json:"Peer,omitempty"`
	PeerMACAddress     string `json:"PeerMACAddress,omitempty"`
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

	switch netdev.Kind {
	case "bond":

		conf += "\n[Bond]\n"

		if netdev.Mode != "" {
			conf += "Mode=" + strings.TrimSpace(netdev.Mode) + "\n"
		}

		if netdev.TransmitHashPolicy != "" {
			conf += "TransmitHashPolicy=" + strings.TrimSpace(netdev.TransmitHashPolicy) + "\n"
		}

		break

	case "vlan":

		conf += "\n[VLAN]\n"

		if netdev.VlanId != "" {
			conf += "Id=" + strings.TrimSpace(netdev.VlanId) + "\n"
		}

		break

	case "bridge":

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

		break
	case "tunnel":

		conf += "\n[Tunnel]\n"

		if netdev.Local != "" {
			conf += "Local=" + strings.TrimSpace(netdev.Local) + "\n"
		}

		if netdev.Remote != "" {
			conf += "Remote=" + strings.TrimSpace(netdev.Remote) + "\n"
		}

		if netdev.TTL != "" {
			conf += "TTL=" + strings.TrimSpace(netdev.TTL) + "\n"
		}

		if netdev.DiscoverPathMTU != "" {
			conf += "DiscoverPathMTU=" + strings.TrimSpace(netdev.DiscoverPathMTU) + "\n"
		}

		if netdev.IPv6FlowLabel != "" {
			conf += "IPv6FlowLabel=" + strings.TrimSpace(netdev.IPv6FlowLabel) + "\n"
		}

		if netdev.EncapsulationLimit != "" {
			conf += "EncapsulationLimit=" + strings.TrimSpace(netdev.EncapsulationLimit) + "\n"
		}

		if netdev.Key != "" {
			conf += "Key=" + strings.TrimSpace(netdev.Key) + "\n"
		}

		if netdev.Independent != "" {
			conf += "Independent=" + strings.TrimSpace(netdev.Independent) + "\n"
		}

		break
	case "veth":

		conf += "\n[Peer]\n"

		if netdev.Peer != "" {
			conf += "Name=" + strings.TrimSpace(netdev.Peer) + "\n"
		}

		if netdev.PeerMACAddress != "" {
			conf += "MACAddress=" + strings.TrimSpace(netdev.PeerMACAddress) + "\n"
		}

		break

	case "macvlan":

		conf += "\n[MACVLAN]\n"

		if netdev.Peer != "" {
			conf += "Mode=" + strings.TrimSpace(netdev.Mode) + "\n"
		}

		break
	case "macvtap":

		conf += "\n[MACVTAP\n"

		if netdev.Peer != "" {
			conf += "Mode=" + strings.TrimSpace(netdev.Mode) + "\n"
		}

		break

	case "ipvlan":

		conf += "\n[IPVLAN]\n"

		if netdev.Peer != "" {
			conf += "Mode=" + strings.TrimSpace(netdev.Mode) + "\n"
		}

	case "vxlan":

		conf += "\n[VXLAN]\n"

		if netdev.Id != "" {
			conf += "Id=" + strings.TrimSpace(netdev.Id) + "\n"
		}

		if netdev.Local != "" {
			conf += "Local=" + strings.TrimSpace(netdev.Local) + "\n"
		}

		if netdev.Remote != "" {
			conf += "Remote=" + strings.TrimSpace(netdev.Remote) + "\n"
		}

		if netdev.TOS != "" {
			conf += "TOS=" + strings.TrimSpace(netdev.TOS) + "\n"
		}

		if netdev.TTL != "" {
			conf += "TTL=" + strings.TrimSpace(netdev.TTL) + "\n"
		}

		if netdev.MacLearning != "" {
			conf += "MacLearning=" + strings.TrimSpace(netdev.MacLearning) + "\n"
		}

		if netdev.DestinationPort != "" {
			conf += "DestinationPort=" + strings.TrimSpace(netdev.DestinationPort) + "\n"
		}

		if netdev.PortRange != "" {
			conf += "PortRange=" + strings.TrimSpace(netdev.PortRange) + "\n"
		}

		if netdev.FlowLabel != "" {
			conf += "FlowLabel=" + strings.TrimSpace(netdev.FlowLabel) + "\n"
		}
	}

	return conf
}

func NetdevdParseJsonFromHttpReq(req *http.Request) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Failed to parse HTTP request: ", err)
		return err
	}

	netdev := new(NetDev)
	json.Unmarshal([]byte(body), &netdev)

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
