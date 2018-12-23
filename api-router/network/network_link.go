// SPDX-License-Identifier: Apache-2.0

package network

import (
	"encoding/json"
	"github.com/vishvananda/netlink"
	log "github.com/sirupsen/logrus"
//	"net"
	//"restgateway/cmd/share"
	"strconv"
	"strings"
	"net/http"
)

type Link struct {
	Action string    `json:"action"`
	Link   string    `json:"link"`
	MTU    string    `json:"mtu"`
	Kind   string    `json:"kind"`
	Enslave []string `json:"enslave"`
}

type LinkInfo struct {
	Index        int    `json:"index"`
	MTU          int    `json:"MTU"`
	TxQLen       int    `json:"TxQLen"`
	Name         string `json:"Name"`
	HardwareAddr string `json:"HardwareAdd"`
	OperState    string `json:"LinkOperState"`
}

func (req *Link) LinkSetMasterBridge() (error) {
	bridge, r := netlink.LinkByName(req.Link)
	if r != nil {
		log.Errorf("Failed to find bridge link %s: %s", req.Link, r)
	}

	br, b := bridge.(*netlink.Bridge)
	if !b {
		log.Errorf("Link is not a bridge %s: %s", req.Link, r)
	}

	for _, n := range req.Enslave {
		link, r := netlink.LinkByName(n)
		if r != nil {
			log.Errorf("Failed to find slave link %s: %s", r)
			continue
		}

		r = netlink.LinkSetMaster(link, br)
		if r != nil {
			log.Errorf("Failed to set link %s master device %s: %s", n, req.Link, r)
		}
	}

	return nil
}

func (req *Link) LinkCreateBridge() (error) {
	_, r := netlink.LinkByName(req.Link)
	if r == nil {
		log.Infof("Bridge link %s exists. Using the bridge", req.Link)
	} else {

		bridge := &netlink.Bridge{LinkAttrs: netlink.LinkAttrs{Name: req.Link}}
		r = netlink.LinkAdd(bridge)
		if r != nil {
			log.Errorf("Failed to create bride %s: %s", req.Link, r)
			return r
		}

		log.Debugf("Successfully create bridge link: %s", req.Link)
	}
	return req.LinkSetMasterBridge()
}

func (req *Link)LinkDelete()(error) {
	l, r := netlink.LinkByName(req.Link)
	if r != nil {
		log.Errorf("Failed to find link %s: %s", req.Link, r)
		return r
	}

	r = netlink.LinkDel(l)
	if (r != nil) {
		log.Errorf("Failed to delete link %s up: %s", l, r)
		return r
	}

	return nil
}

func LinkSetUp(link string)(error) {
	l, r := netlink.LinkByName(link)
	if r != nil {
		log.Errorf("Failed to find link %s: %s", link, r)
		return r
	}

	r = netlink.LinkSetUp(l)
	if (r != nil) {
		log.Errorf("Failed to set link %s up: %s", l, r)
		return r
	}

	return nil
}

func LinkSetDown(link string) (error) {
	l, r := netlink.LinkByName(link)
	if r != nil {
		log.Errorf("Failed to find link %s: %s", link, r)
		return r
	}

	r = netlink.LinkSetDown(l)
	if (r != nil) {
		log.Errorf("Failed to set link down %s: %s", l, r)
		return r
	}

	return nil
}

func LinkSetMTU(link string, mtu int) (error) {
	l, r := netlink.LinkByName(link)
	if r != nil {
		log.Errorf("Failed to find link %s: %s", link, r)
		return r
	}

	r = netlink.LinkSetMTU(l, mtu)
	if (r != nil) {
		log.Errorf("Failed to set link %s MTU %d: %s", link, mtu, r)
		return r
	}

	return nil
}

func (req *Link) SetLink() (error){

	link := strings.TrimSpace(req.Link)

	switch req.Action {
	case "set-link-up":
		return LinkSetUp(link)
	case "set-link-down":
		return LinkSetDown(link)
	case "set-link-mtu":

		mtu, r := strconv.ParseInt(strings.TrimSpace(req.MTU), 10, 64)
		if (r != nil) {
			log.Errorf("Failed to parse received link %s MTU %s: %s", req.Link, req.MTU, r)
			return r
		}

		return LinkSetMTU(link, int(mtu))
	}

	return nil
}

func (req *Link) GetLink(rw http.ResponseWriter) {
	link, r := netlink.LinkByName(strings.TrimSpace(req.Link))
	if r != nil {
		log.Errorf("Failed to find link %s: ", req.Link, r)
		return
	}

	linkInfo := LinkInfo {
		Index: link.Attrs().Index,
		MTU: link.Attrs().MTU,
		Name: link.Attrs().Name,
		HardwareAddr: link.Attrs().HardwareAddr.String(),
	}

	json.NewEncoder(rw).Encode(linkInfo)
}
