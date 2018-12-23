// SPDX-License-Identifier: Apache-2.0

package networkd

import (
	"api-routerd/cmd/share"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strings"
)

type Address struct {
	Address string
	Peer    string
	Label   string
	Scope   string
}

type Route struct {
	Gateway         string
	GatewayOnlink   string
	Destination     string
	Source          string
	PreferredSource string
	Table           string
}

type Network struct {
	MAC    string
	Name   string
	Driver string

	Addresses interface{}
	Routes    interface{}

	Gateway             string
	DHCP                string
	DNS                 string
	Domains             string
	NTP                 string
	IPv6AcceptRA        string
	LinkLocalAddressing string
	LLDP                string
	EmitLLDP            string

	Bridge  string
	Bond    string
	VRF     string
	VLAN    string
	MACVLAN string
	VXLAN   string
	Tunnel  string
}

func (network *Network) CreateMatchSectionConfig() string {
	conf := "[Match]\n"

	if network.MAC != "" {
		mac := fmt.Sprintf("MACAddress=%s\n", network.MAC)
		conf += mac
	}

	if network.Driver != "" {
		driver := fmt.Sprintf("Driver=%s\n", network.Driver)
		conf += driver
	}

	if network.Name != "" {
		name := fmt.Sprintf("Name=%s\n", network.Name)
		conf += name
	}

	return conf
}

func (network *Network) CreateRouteSectionConfig() string {
	var routeConf string

	switch v := network.Routes.(type) {
	case []interface{}:
		for _, b := range v {
			var preferredSource string
			var gatewayOnLink string
			var destination string
			var gateway string
			var source string
			var table string

			if b.(map[string]interface{})["Gateway"] != nil {
				gateway = strings.TrimSpace(b.(map[string]interface{})["Gateway"].(string))
			}

			if b.(map[string]interface{})["GatewayOnlink"] != nil {
				gatewayOnLink = strings.TrimSpace(b.(map[string]interface{})["GatewayOnlink"].(string))
			}

			if b.(map[string]interface{})["Destination"] != nil {
				destination = strings.TrimSpace(b.(map[string]interface{})["Destination"].(string))
			}

			if b.(map[string]interface{})["Source"] != nil {
				source = strings.TrimSpace(b.(map[string]interface{})["Source"].(string))
			}

			if b.(map[string]interface{})["PreferredSource"] != nil {
				preferredSource = strings.TrimSpace(b.(map[string]interface{})["PreferredSource"].(string))
			}

			if b.(map[string]interface{})["Table"] != nil {
				table = strings.TrimSpace(b.(map[string]interface{})["Table"].(string))
			}

			routeConf += "\n[Route]\n"

			if len(gateway) != 0 {
				ip := net.ParseIP(gateway)
				if ip != nil {
					routeConf += "Gateway=" + gateway + "\n"
				} else {
					log.Error("Failed to parse Gateway: ", gateway)
				}
			}

			if len(gatewayOnLink) != 0 {
				onlink := strings.TrimSpace(gatewayOnLink)
				b, r := share.ParseBool(onlink)
				if r != nil {
					log.Error("Failed to parse GatewayOnlink: ", r, gatewayOnLink)
				} else {
					if b == true {
						routeConf += "GatewayOnlink=yes\n"
					} else {
						routeConf += "GatewayOnlink=no\n"
					}
				}
			}

			if len(destination) != 0 {
				ip := net.ParseIP(destination)
				if ip != nil {
					routeConf += "Destination=" + destination + "\n"
				} else {
					log.Error("Failed to parse Destination: ", destination)
				}
			}

			if len(source) != 0 {
				ip := net.ParseIP(source)
				if ip != nil {
					routeConf += "Source=" + source + "\n"
				} else {
					log.Error("Failed to parse Source: ", source)
				}
			}

			if len(preferredSource) != 0 {
				ip := net.ParseIP(preferredSource)
				if ip != nil {
					routeConf += "PreferredSource=" + preferredSource + "\n"
				} else {
					log.Error("Failed to parse PreferredSource: ", preferredSource)
				}
			}

			if len(table) != 0 {
				routeConf += "Table=" + table + "\n"
			}
		}
		break
	}

	return routeConf
}

func (network *Network) CreateAddressSectionConfig() string {
	var addressConf string

	switch v := network.Addresses.(type) {
	case []interface{}:
		for _, b := range v {
			var address string
			var peer string
			var scope string
			var label string

			if b.(map[string]interface{})["Address"] != nil {
				address = strings.TrimSpace(b.(map[string]interface{})["Address"].(string))
			}

			if b.(map[string]interface{})["Peer"] != nil {
				peer = strings.TrimSpace(b.(map[string]interface{})["Peer"].(string))
			}

			if b.(map[string]interface{})["Scope"] != nil {
				scope = strings.TrimSpace(b.(map[string]interface{})["Scope"].(string))
			}

			if b.(map[string]interface{})["Label"] != nil {
				label = strings.TrimSpace(b.(map[string]interface{})["Label"].(string))
			}

			if len(address) != 0 {
				addressConf += "\n[Address]\nAddress="

				ip := net.ParseIP(address)
				if ip != nil {
					addressConf += address + "\n"
				} else {
					log.Error("Failed to parse address: ", address)
				}

				ip = net.ParseIP(peer)
				if ip != nil {
					addressConf += "Peer=" + peer + "\n"
				} else {
					log.Error("Failed to parse peer address: ", peer)
				}

				if len(scope) != 0 {
					addressConf += "Scope=" + scope + "\n"
				}

				if len(label) != 0 {
					addressConf += "Label=" + label + "\n"
				}
			}
		}
		break
	}

	return addressConf
}

func (network *Network) CreateNetworkSectionConfig() string {
	conf := "[Network]\n"

	if network.DHCP != "" {
		dhcpConf := "DHCP="

		dhcp := strings.TrimSpace(network.DHCP)
		b, r := share.ParseBool(dhcp)
		if r != nil {
			switch dhcp {
			case "ipv4", "ipv6":
				dhcpConf += dhcp
				break
			default:
				log.Error("Failed to parse DHCP: ", r, network.DHCP)
			}
		} else {
			if b == true {
				dhcpConf += "yes"
			} else {
				dhcpConf += "no"
			}
		}
		conf += dhcpConf + "\n"
	}

	if network.Gateway != "" {
		gatewayConf := "Gateway="

		gw := strings.TrimSpace(network.Gateway)
		ip := net.ParseIP(gw)
		if ip != nil {
			gatewayConf += gw
			conf += gatewayConf + "\n"
		} else {
			log.Error("Failed to parse Gateway Address: ", network.Gateway)
		}
	}

	if network.DNS != "" {
		conf += "DNS=" + network.DNS
	}

	if network.Domains != "" {
		conf += "Domains=" + network.Domains + "\n"
	}

	if network.NTP != "" {
		conf += "NTP=" + network.NTP + "\n"
	}

	if network.IPv6AcceptRA != "" {
		IPv6AcceptRAConf := "IPv6AcceptRA="

		IPv6RA := strings.TrimSpace(network.IPv6AcceptRA)
		b, r := share.ParseBool(IPv6RA)
		if r != nil {
			log.Error("Failed to parse IPv6AcceptRA: ", r, network.IPv6AcceptRA)
		} else {
			if b == true {
				IPv6AcceptRAConf += "yes"
			} else {
				IPv6AcceptRAConf += "no"
			}
		}
		conf += IPv6AcceptRAConf + "\n"
	}

	if network.LinkLocalAddressing != "" {
		LinkLocalAddressingConf := "LinkLocalAddressing="

		IPv6RA := strings.TrimSpace(network.LinkLocalAddressing)
		b, r := share.ParseBool(IPv6RA)
		if r != nil {
			log.Error("Failed to parse LinkLocalAddressing: ", r, network.LinkLocalAddressing)
		} else {
			if b == true {
				LinkLocalAddressingConf += "yes"
			} else {
				LinkLocalAddressingConf += "no"
			}
		}
		conf += LinkLocalAddressingConf + "\n"
	}

	if network.LLDP != "" {
		LLDPConf := "LLDP="

		LLDP := strings.TrimSpace(network.LLDP)
		b, r := share.ParseBool(LLDP)
		if r != nil {
			log.Error("Failed to parse LLDP: ", r, network.LLDP)
		} else {
			if b == true {
				LLDPConf += "yes"
			} else {
				LLDPConf += "no"
			}
		}
		conf += LLDPConf + "\n"
	}

	if network.EmitLLDP != "" {
		EmitLLDPConf := "EmitLLDP="

		EmitLLDP := strings.TrimSpace(network.EmitLLDP)
		b, r := share.ParseBool(EmitLLDP)
		if r != nil {
			log.Error("Failed to parse EmitLLDP: ", r, network.EmitLLDP)
		} else {
			if b == true {
				EmitLLDPConf += "yes"
			} else {
				EmitLLDPConf += "no"
			}
		}
		conf += EmitLLDPConf + "\n"
	}

	if network.NTP != "" {
		conf += "NTP=" + network.NTP + "\n"
	}

	if network.Bridge != "" {
		conf += "Bridge=" + network.Bridge + "\n"
	}

	if network.Bond != "" {
		conf += "Bond=" + network.Bond + "\n"
	}

	if network.VLAN != "" {
		conf += "VLAN=" + network.VLAN + "\n"
	}

	return conf
}

func NetworkdParseJsonFromHttpReq(req *http.Request) error {
	var configs map[string]interface{}

	body, r := ioutil.ReadAll(req.Body)
	if r != nil {
		log.Errorf("Failed to parse HTTP request: %s ", r)
		return r
	}

	json.Unmarshal([]byte(body), &configs)

	var network Network
	for key, value := range configs {
		switch key {
		case "MAC":
			network.MAC = value.(string)
			break
		case "Name":
			network.Name = value.(string)
			break
		case "Driver":
			network.Driver = value.(string)
			break
		case "Addresses":
			network.Addresses = value
		case "Routes":
			network.Routes = value
		case "Gateway":
			network.Driver = value.(string)
			break
		case "DHCP":
			network.DHCP = value.(string)
			break
		case "Domains":
			network.Domains = value.(string)
			break
		case "DNS":
			network.DNS = value.(string)
			break
		case "NTP":
			network.NTP = value.(string)
			break
		case "IPv6AcceptRA":
			network.IPv6AcceptRA = value.(string)
			break
		case "LinkLocalAddressing":
			network.LinkLocalAddressing = value.(string)
			break
		case "LLDP":
			network.LLDP = value.(string)
			break
		case "EmitLLDP":
			network.EmitLLDP = value.(string)
			break
		case "Bridge":
			network.Bridge = value.(string)
			break
		case "Bond":
			network.Bond = value.(string)
			break
		case "VLAN":
			network.VLAN = value.(string)
			break
		}
	}

	matchConfig := network.CreateMatchSectionConfig()
	networkConfig := network.CreateNetworkSectionConfig()
	addressConfig := network.CreateAddressSectionConfig()
	routeConfig := network.CreateRouteSectionConfig()

	config := []string{matchConfig, networkConfig, addressConfig, routeConfig}

	fmt.Println(config)

	unitName := fmt.Sprintf("25-%s.network", network.Name)
	unitPath := filepath.Join(NetworkdUnitPath, unitName)

	share.WriteFullFile(unitPath, config)

	return nil
}

func ConfigureNetworkFile(rw http.ResponseWriter, req *http.Request) {
	NetworkdParseJsonFromHttpReq(req)
}
