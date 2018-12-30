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
)

type Link struct {
	MAC                        string `json:"MAC"`
	MatchName                  string `json:"MatchName"`
	Driver                     string `json:"Driver"`

	Description                string `json:"Description"`
	Alias                      string `json:"Alias"`
	MACAddressPolicy           string `json:"MACAddressPolicy"`
	MACAddress                 string `json:"MACAddress"`
	NamePolicy                 string `json:"NamePolicy"`
	Name                       string `json:"Name"`
	MTUBytes                   string `json:"MTUBytes"`
	BitsPerSecond              string `json:"BitsPerSecond"`
	Duplex                     string `json:"Duplex"`
	AutoNegotiation            string `json:"AutoNegotiation"`
	WakeOnLan                  string `json:"WakeOnLan"`
	Port                       string `json:"Port"`
	TCPSegmentationOffload     string `json:"TCPSegmentationOffload"`
	TCP6SegmentationOffload    string `json:"TCP6SegmentationOffload"`
	GenericSegmentationOffload string `json:"GenericSegmentationOffload"`
	GenericReceiveOffload      string `json:"GenericReceiveOffload"`
	LargeReceiveOffload        string `json:"LargeReceiveOffload"`
	RxChannels                 string `json:"RxChannels"`
	TxChannels                 string `json:"TxChannels"`
	OtherChannels              string `json:"OtherChannels"`
	CombinedChannels           string `json:"CombinedChannels"`
}

func (link *Link) CreateLinkMatchSectionConfig() string {
	conf := "[Match]\n"

	if link.MAC != "" {
		mac := fmt.Sprintf("MACAddress=%s\n", link.MAC)
		conf += mac
	}

	if link.Driver != "" {
		driver := fmt.Sprintf("Driver=%s\n", link.Driver)
		conf += driver
	}

	if link.MatchName != "" {
		name := fmt.Sprintf("Name=%s\n", link.MatchName)
		conf += name
	}

	return conf
}

func (link *Link) CreateLinkSectionConfig() string {
	conf := "[Link]\n"

	if link.Description != "" {
		conf += "Description=" + link.Description + "\n"
	}

	if link.Alias != "" {
		conf += "Alias=" + link.Alias + "\n"
	}

	if link.MACAddressPolicy != "" {
		conf += "MACAddressPolicy=" + link.MACAddressPolicy + "\n"
	}

	if link.MACAddress != "" {
		conf += "MACAddress=" + link.MACAddress + "\n"
	}

	if link.NamePolicy != "" {
		conf += "NamePolicy=" + link.NamePolicy + "\n"
	}

	if link.Name != "" {
		conf += "Name=" + link.Name + "\n"
	}

	if link.MTUBytes != "" {
		conf += "MTUBytes=" + link.MTUBytes + "\n"
	}

	if link.BitsPerSecond != "" {
		conf += "BitsPerSecond=" + link.BitsPerSecond + "\n"
	}

	if link.Duplex != "" {
		conf += "Duplex=" + link.Duplex + "\n"
	}

	if link.AutoNegotiation != "" {
		conf += "AutoNegotiation=" + link.AutoNegotiation + "\n"
	}

	if link.WakeOnLan != "" {
		conf += "WakeOnLan=" + link.WakeOnLan + "\n"
	}

	if link.Port != "" {
		conf += "Port=" + link.Port + "\n"
	}

	if link.TCPSegmentationOffload != "" {
		conf += "TCPSegmentationOffload=" + link.TCPSegmentationOffload + "\n"
	}

	if link.TCP6SegmentationOffload != "" {
		conf += "TCP6SegmentationOffload=" + link.TCP6SegmentationOffload + "\n"
	}

	if link.GenericSegmentationOffload != "" {
		conf += "GenericSegmentationOffload=" + link.GenericSegmentationOffload + "\n"
	}

	if link.GenericReceiveOffload != "" {
		conf += "GenericReceiveOffload=" + link.GenericReceiveOffload + "\n"
	}

	if link.LargeReceiveOffload != "" {
		conf += "LargeReceiveOffload=" + link.LargeReceiveOffload + "\n"
	}

	if link.RxChannels != "" {
		conf += "RxChannels=" + link.RxChannels + "\n"
	}

	if link.TxChannels != "" {
		conf += "TxChannels=" + link.TxChannels + "\n"
	}

	if link.OtherChannels != "" {
		conf += "OtherChannels=" + link.OtherChannels + "\n"
	}

	if link.CombinedChannels != "" {
		conf += "CombinedChannels=" + link.CombinedChannels + "\n"
	}

	return conf
}

func LinkParseJsonFromHttpReq(req *http.Request) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("Failed to parse HTTP request: %s ", err)
		return err
	}

	link := new(Link)
	json.Unmarshal([]byte(body), &link)

	matchConfig := link.CreateLinkMatchSectionConfig()
	linkConfig := link.CreateLinkSectionConfig()

	config := []string{matchConfig, linkConfig}

	fmt.Println(config)

	unitName := fmt.Sprintf("00-%s.link", link.Name)
	unitPath := filepath.Join(NetworkdUnitPath, unitName)

	return share.WriteFullFile(unitPath, config)
}

func ConfigureLinkFile(rw http.ResponseWriter, req *http.Request) {
	LinkParseJsonFromHttpReq(req)
}
