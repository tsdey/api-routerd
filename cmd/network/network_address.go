// SPDX-License-Identifier: Apache-2.0

package network

import (
	"encoding/json"
	"net/http"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

type Address struct {
	Action  string `json:"action"`
	Link    string `json:"link"`
	Address string `json:"address"`
	Label   string `json:"label"`
}

func (address *Address) AddAddress() (error){
	link, err := netlink.LinkByName(address.Link)
	if err != nil {
		log.Errorf("Failed to find link %s: %s", err, address.Link)
		return err
	}

	addr, err := netlink.ParseAddr(address.Address)
	if err != nil {
		log.Errorf("Failed to parse address %s: %s", err, address.Address)
		return err
	}

	err = netlink.AddrAdd(link, addr)
	if err != nil {
		log.Errorf("Failed to add Address %s to link %s: %s", err, addr, link)
		return err
	}

	return nil
}

func (address *Address) DelAddress() (error) {
	link, err := netlink.LinkByName(address.Link)
	if err != nil {
		log.Errorf("Failed to find link %s: %s", err, address.Link)
		return err
	}

	addr, err := netlink.ParseAddr(address.Address)
	if err != nil {
		log.Errorf("Failed to parse address %s: %s", err, addr)
		return err
	}

	err = netlink.AddrDel(link, addr)
	if err != nil {
		log.Errorf("Failed to add address %s: %s", err, addr, link)
		return err
	}

	return nil
}

func (address *Address) GetAddress(rw http.ResponseWriter) (error) {
	link, err := netlink.LinkByName(address.Link)
	if err != nil {
		log.Errorf("Failed to get link %s: %s", address.Link, err)
		return err
	}

	addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err != nil {
		log.Errorf("Could not get addresses for link %s: %s", link, err)
		return err
	}

	j, err := json.Marshal(addrs)
	if err != nil {
		log.Errorf("Failed to encode json address for link %s: %s", err, address.Link)
		return err
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(j)

	return nil
}
