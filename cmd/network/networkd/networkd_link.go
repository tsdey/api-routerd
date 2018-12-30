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

type Link struct {
	ConfFile                   string      `json:ConfFile",omitempty"`
	Match                      interface{} `json:Match",omitempty"`

	Description                string      `json:"Description,omitempty"`
	Alias                      string      `json:"Alias,omitempty"`
	MACAddressPolicy           string      `json:"MACAddressPolicy,omitempty"`
	MACAddress                 string      `json:"MACAddress,omitempty"`
	NamePolicy                 string      `json:"NamePolicy,omitempty"`
	Name                       string      `json:"Name,omitempty"`
	MTUBytes                   string      `json:"MTUBytes,omitempty"`
	BitsPerSecond              string      `json:"BitsPerSecond,omitempty"`
	Duplex                     string      `json:"Duplex,omitempty"`
	AutoNegotiation            string      `json:"AutoNegotiation,omitempty"`
	WakeOnLan                  string      `json:"WakeOnLan,omitempty"`
	Port                       string      `json:"Port,omitempty"`
	TCPSegmentationOffload     string      `json:"TCPSegmentationOffload,omitempty"`
	TCP6SegmentationOffload    string      `json:"TCP6SegmentationOffload,omitempty"`
	GenericSegmentationOffload string      `json:"GenericSegmentationOffload,omitempty"`
	GenericReceiveOffload      string      `json:"GenericReceiveOffload,omitempty"`
	LargeReceiveOffload        string      `json:"LargeReceiveOffload,omitempty"`
	RxChannels                 string      `json:"RxChannels,omitempty"`
	TxChannels                 string      `json:"TxChannels,omitempty"`
	OtherChannels              string      `json:"OtherChannels,omitempty"`
	CombinedChannels           string      `json:"CombinedChannels,omitempty"`
}

func (link *Link) CreateLinkMatchSectionConfig() string {
	conf := "[Match]\n"

	switch v := link.Match.(type) {
	case []interface{}:
		for _, b := range v {
			var mac string
			var driver string
			var name string

			if b.(map[string]interface{})["MAC"] != nil {
				mac = strings.TrimSpace(b.(map[string]interface{})["MAC"].(string))
			}

			if b.(map[string]interface{})["Driver"] != nil {
				driver = strings.TrimSpace(b.(map[string]interface{})["Driver"].(string))
			}

			if b.(map[string]interface{})["Name"] != nil {
				name = strings.TrimSpace(b.(map[string]interface{})["Name"].(string))
			}


			if mac != "" {
				mac := fmt.Sprintf("MACAddress=%s\n", mac)
				conf += mac
			}

			if driver != "" {
				driver := fmt.Sprintf("Driver=%s\n", driver)
				conf += driver
			}

			if name != "" {
				if link.ConfFile == "" {
					link.ConfFile = name
				}

				name := fmt.Sprintf("Name=%s\n", name)
				conf += name
			}
		}
		break
	}

	fmt.Println(conf)
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
