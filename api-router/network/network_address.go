// SPDX-License-Identifier: Apache-2.0

package network

import (
	"encoding/json"
	"github.com/vishvananda/netlink"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Address struct {
	Action string  `json:"action"`
	Link string    `json:"link"`
	Address string `json:"address"`
	Label string   `json:"label"`
}

func AddAddress(address *Address) {
	link, r := netlink.LinkByName(address.Link)
	if r != nil {
		log.Errorf("Failed to find link %s: %s", r, address.Link)
		return
	}

	addr, r := netlink.ParseAddr(address.Address)
	if r != nil {
		log.Errorf("Failed to parse address %s: %s", r, address.Address)
		return
	}

	r = netlink.AddrAdd(link, addr)
	if r != nil {
		log.Errorf("Failed to add Address %s: %s", r, addr, link)
		return
	}
}

func DelAddress(address *Address) {
	link, r := netlink.LinkByName(address.Link)
	if r != nil {
		log.Errorf("Failed to find link %s: %s", r, address.Link)
		return
	}

	addr, r := netlink.ParseAddr(address.Address)
	if r != nil {
		log.Errorf("Failed to parse address %s: %s", r, addr)
		return
	}

	r = netlink.AddrDel(link, addr)
	if r != nil {
		log.Errorf("Failed to add address %s: %s", r, addr, link)
		return
	}
}

func GetAddress(rw http.ResponseWriter, address *Address) {
	link, r := netlink.LinkByName(address.Link)
	if r != nil {
		log.Errorf("Failed to get link %s: %s", r, address.Link)
		return
	}

	addrs, r := netlink.AddrList(link, netlink.FAMILY_ALL)
	if r != nil {
		log.Errorf("Could not get addresses %s: %s", link, r)
		return
	}

	addresses := make([]Address, len(addrs))
	for i, address := range addrs {
		addresses[i].Address = address.IPNet.String()
		addresses[i].Link = link.Attrs().Name
	}

	json.NewEncoder(rw).Encode(addresses)
}
