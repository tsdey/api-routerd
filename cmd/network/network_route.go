// SPDX-License-Identifier: Apache-2.0

package network

import (
	"api-routerd/cmd/share"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"syscall"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type Route struct {
	Action  string `json:"action"`
	Link    string `json:"link"`
	Gateway string `json:"gateway"`
	OnLink  string `json:"onlink"`
}

func (route *Route) AddDefaultGateWay() (error) {
	link, err := netlink.LinkByName(route.Link)
	if err != nil {
		log.Errorf("Failed to find link %s: %s", err, route.Link)
		return err
	}

	ipAddr, _, err := net.ParseCIDR(route.Gateway)
	if err != nil {
		log.Errorf("Failed to parse default GateWay address %s: %s", route.Gateway, err)
		return err
	}

	onlink := 0
	b, err := share.ParseBool(strings.TrimSpace(route.OnLink))
	if err != nil {
		log.Errorf("Failed to parse GatewayOnlink %s: %s", err, route.OnLink)
	} else {
		if b == true {
			onlink |= syscall.RTNH_F_ONLINK
		}
	}

	// add a gateway route
	rt := &netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: link.Attrs().Index,
		Gw:        ipAddr,
		Flags:     onlink,
	}

	err = netlink.RouteAdd(rt)
	if err != nil {
		log.Errorf("Failed to add default GateWay address %s: %s", route.Gateway, err)
	}

	return nil
}

func (route *Route) DelDefaultGateWay() (error) {
	link, err := netlink.LinkByName(route.Link)
	if err != nil {
		log.Errorf("Failed to delete default gateway %s: %s", link, err)
		return err
	}

	ipAddr, _, err := net.ParseCIDR(route.Gateway)
	if err != nil {
		log.Errorf("Failed to parse default GateWay address %s: %s", route.Gateway, err)
		return err
	}

	// del a gateway route
	r := &netlink.Route{
		Scope:     netlink.SCOPE_UNIVERSE,
		LinkIndex: link.Attrs().Index,
		Gw:        ipAddr,
	}

	err = netlink.RouteDel(r)
	if err != nil {
		log.Errorf("Failed to delete default GateWay address %s: %s", ipAddr, err)
		return err
	}

	return nil
}

func GetRoutes(rw http.ResponseWriter) (error) {
	routes, err := netlink.RouteList(nil, 0)
	if err != nil {
		log.Errorf("Failed to get routes %s", err)
		return err
	}

	j, err := json.Marshal(routes)
	if err != nil {
		log.Errorf("Failed to encode json routeInfo : %s", err)
		return err
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(j)

	return nil
}
