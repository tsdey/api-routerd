package resolv

import (
	"api-routerd/cmd/share"
	"bufio"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"fmt"
)

const ResolvConfPath = "/etc/resolv.conf"

type DnsConfig struct {
	Servers []string `json:"servers"`
	Search  []string `json:"search"`
}

func (conf *DnsConfig) WriteResolvConfig() (err error) {
	f, err := os.OpenFile(ResolvConfPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for _, server := range conf.Servers {
		line := "nameserver " + server
		fmt.Fprintln(w, line)
	}
	for _, s := range conf.Search {
		line := "search " + s
		fmt.Fprintln(w, line)
	}

	w.Flush()

	return nil
}

func ReadResolvConf() (*DnsConfig, error) {
	lines, err := share.ReadFullFile(ResolvConfPath)
	if err != nil {
		log.Fatal("Failed to read: %s", ResolvConfPath)
		return nil, errors.New("Failed to read resolv")
	}

	conf := new(DnsConfig)

	for _, line := range lines {
		fields := strings.Fields(line)

		switch fields[0] {
		case "nameserver":
			conf.Servers = append(conf.Servers, fields[1])
			break
		case "search":

			for i, search := range strings.Fields(line) {
				if i == 0 {
					continue
				}

				conf.Search = append(conf.Search, search)
			}
		}
	}

	// Don't return nil in json
	if (len(conf.Servers) == 0) {
		conf.Servers = []string{""}
	}

	if (len(conf.Search) == 0) {
		conf.Search = []string{""}
	}

	return conf, nil
}

func GetResolvConf(rw http.ResponseWriter) (error) {
	conf, err := ReadResolvConf()
	if err != nil {
		return err
	}

	j, err := json.Marshal(conf)
	if err != nil {
		log.Errorf("Failed to encode json for resolv %s", err)
		return err
	}

	rw.Write(j)

	return nil
}

func UpdateResolvConf(rw http.ResponseWriter, r *http.Request) (error) {
	dns := DnsConfig{
		Servers: []string{""},
		Search: []string{""},
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to parse HTTP request: ", err)
		return err
	}

	err = json.Unmarshal([]byte(body), &dns)
	if err != nil {
		log.Error("Failed to Decode HTTP request to json: ", err)
		return err
	}

	conf, err := ReadResolvConf()
	if err != nil {
		return err
	}

	// update nameserver
	for _, s := range dns.Servers {
		if share.StringContains(conf.Servers, s) {
			continue
		}

		conf.Servers = append(conf.Servers, s)
	}

	// update domains
	for _, s := range dns.Search {
		if share.StringContains(conf.Search, s) {
			continue
		}
		conf.Search = append(conf.Search, s)
	}

	err = conf.WriteResolvConfig()
	if err != nil {
		log.Errorf("Failed Write to resolv conf: %s", err)
		return err
	}

	j, err := json.Marshal(conf)
	if err != nil {
		log.Errorf("Failed to encode json for resolv %s", err)
		return err
	}

	rw.Write(j)

	return nil
}

func DeleteResolvConf(rw http.ResponseWriter, r *http.Request) (error) {
	dns := DnsConfig{
		Servers: []string{""},
		Search: []string{""},
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to parse HTTP request: ", err)
		return err
	}

	err = json.Unmarshal([]byte(body), &dns)
	if err != nil {
		log.Error("Failed to Decode HTTP request to json: ", err)
		return err
	}

	conf, err := ReadResolvConf()
	if err != nil {
		return err
	}

	// update nameserver
	for _, s := range dns.Servers {
		if !share.StringContains(conf.Servers, s) {
			continue
		}

		conf.Servers, _ = share.StringDeleteSlice(conf.Servers, s)
	}

	// update domains
	for _, s := range dns.Search {
		if !share.StringContains(conf.Search, s) {
			continue
		}

		conf.Search, _ = share.StringDeleteSlice(conf.Search, s)
	}

	err = conf.WriteResolvConfig()
	if err != nil {
		log.Errorf("Failed Write to resolv conf: %s", err)
		return err
	}

	j, err := json.Marshal(conf)
	if err != nil {
		log.Errorf("Failed to encode json for resolv %s", err)
		return err
	}

	rw.Write(j)

	return nil
}
