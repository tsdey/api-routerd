// SPDX-License-Identifier: Apache-2.0

package network

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"net/http"
	"strconv"
	"strings"
)

type Link struct {
	Action  string   `json:"action"`
	Link    string   `json:"link"`
	MTU     string   `json:"mtu"`
	Kind    string   `json:"kind"`
	Enslave []string `json:"enslave"`
}

func (req *Link) LinkSetMasterBridge() (error) {
	bridge, err := netlink.LinkByName(req.Link)
	if err != nil {
		log.Errorf("Failed to find bridge link %s: %s", req.Link, err)
		return err
	}

	br, b := bridge.(*netlink.Bridge)
	if !b {
		log.Errorf("Link is not a bridge: %s", req.Link)
		return errors.New("Link is not a bridge")
	}

	for _, n := range req.Enslave {
		link, err := netlink.LinkByName(n)
		if err != nil {
			log.Errorf("Failed to find slave link %s: %s", n, err)
			continue
		}

		err = netlink.LinkSetMaster(link, br)
		if err != nil {
			log.Errorf("Failed to set link %s master device %s: %s", n, req.Link, err)
		}
	}

	return nil
}

func (req *Link) LinkCreateBridge() (error) {
	_, err := netlink.LinkByName(req.Link)
	if err == nil {
		log.Infof("Bridge link %s exists. Using the bridge", req.Link)
	} else {

		bridge := &netlink.Bridge{LinkAttrs: netlink.LinkAttrs{Name: req.Link}}
		err = netlink.LinkAdd(bridge)
		if err != nil {
			log.Errorf("Failed to create bride %s: %s", req.Link, err)
			return err
		}

		log.Debugf("Successfully create bridge link: %s", req.Link)
	}
	return req.LinkSetMasterBridge()
}

func (req *Link) LinkDelete() (error) {
	l, err := netlink.LinkByName(req.Link)
	if err != nil {
		log.Errorf("Failed to find link %s: %s", req.Link, err)
		return err
	}

	err = netlink.LinkDel(l)
	if err != nil {
		log.Errorf("Failed to delete link %s up: %s", l, err)
		return err
	}

	return nil
}

func LinkSetUp(link string) (error) {
	l, err := netlink.LinkByName(link)
	if err != nil {
		log.Errorf("Failed to find link %s: %s", link, err)
		return err
	}

	err = netlink.LinkSetUp(l)
	if err != nil {
		log.Errorf("Failed to set link %s up: %s", l, err)
		return err
	}

	return nil
}

func LinkSetDown(link string) (error) {
	l, err := netlink.LinkByName(link)
	if err != nil {
		log.Errorf("Failed to find link %s: %s", link, err)
		return err
	}

	err = netlink.LinkSetDown(l)
	if err != nil {
		log.Errorf("Failed to set link down %s: %s", l, err)
		return err
	}

	return nil
}

func LinkSetMTU(link string, mtu int) (error) {
	l, err := netlink.LinkByName(link)
	if err != nil {
		log.Errorf("Failed to find link %s: %s", link, err)
		return err
	}

	err = netlink.LinkSetMTU(l, mtu)
	if err != nil {
		log.Errorf("Failed to set link %s MTU %d: %s", link, mtu, err)
		return err
	}

	return nil
}

func (req *Link) SetLink() (error) {
	link := strings.TrimSpace(req.Link)

	switch req.Action {
	case "set-link-up":
		return LinkSetUp(link)
	case "set-link-down":
		return LinkSetDown(link)
	case "set-link-mtu":

		mtu, err := strconv.ParseInt(strings.TrimSpace(req.MTU), 10, 64)
		if err != nil {
			log.Errorf("Failed to parse received link %s MTU %s: %s", req.Link, req.MTU, err)
			return err
		}

		return LinkSetMTU(link, int(mtu))
	}

	return nil
}

func (req *Link) GetLink(rw http.ResponseWriter) (error) {
	l := strings.TrimSpace(req.Link)
	if l != "" {
		link, err := netlink.LinkByName(l)
		if err != nil {
			log.Errorf("Failed to find link %s: %s", req.Link, err)
			return err
		}

		j, err := json.Marshal(link)
		if err != nil {
			log.Errorf("Failed to encode json linkInfo for link %s: %s", req.Link, err)
			return err
		}


		rw.WriteHeader(http.StatusOK)
		rw.Write(j)

	} else	{

		links, err := netlink.LinkList()
		if err != nil {
			return err
		}

		j, err := json.Marshal(links)
		if err != nil {
			log.Errorf("Failed to encode json linkInfo for link %s: %s", req.Link, err)
			return err
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(j)
	}

	return nil
}
