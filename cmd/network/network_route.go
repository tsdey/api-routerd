// SPDX-License-Identifier: Apache-2.0

package network

import (
	"api-routerd/cmd/share"
	"net"
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

func AddDefaultGateWay(route *Route) {
	link, r := netlink.LinkByName(route.Link)
	if r != nil {
		log.Errorf("Failed to find link %s: %s", r, route.Link)
		return
	}

	ipAddr, _, r := net.ParseCIDR(route.Gateway)
	if r != nil {
		log.Errorf("Failed to parse default GateWay address %s: %s", route.Gateway, r)
		return
	}

	onlink := 0
	b, r := share.ParseBool(strings.TrimSpace(route.OnLink))
	if r != nil {
		log.Errorf("Failed to parse GatewayOnlink %s: %s", r, route.OnLink)
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

	r = netlink.RouteAdd(rt)
	if r != nil {
		log.Errorf("Failed to add default GateWay address %s: %s", route.Gateway, r)
	}
}

func DelDefaultGateWay(route *Route) {
	link, err := netlink.LinkByName(route.Link)
	if err != nil {
		log.Errorf("Failed to delete default gateway %s: %s", link, err)
		return
	}

	ipAddr, _, err := net.ParseCIDR(route.Gateway)
	if err != nil {
		log.Errorf("Failed to parse default GateWay address %s: %s", route.Gateway, err)
		return
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
	}
}
