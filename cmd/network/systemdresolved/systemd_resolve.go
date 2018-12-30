package systemdresolved

import (
	"api-routerd/cmd/share"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	ResolvedConfPath = "/etc/systemd/resolved.conf"
)

type DnsConfig struct {
	DNS          []string `json:"dns"`
	FallbackDNS  []string `json:"fallback_dns"`
}

func (d *DnsConfig) WriteResolveConfig() (error) {
	f, err := os.OpenFile(ResolvedConfPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	conf := "[Resolve]\n"

	dnsConf := "DNS="
	for _, s := range d.DNS {
		dnsConf += s + " "
	}
	conf += dnsConf + "\n"

	fallbackDns := "FallbackDNS="
	for _, s := range d.FallbackDNS {
		fallbackDns +=  s + " "
	}
	conf += fallbackDns + "\n"

	fmt.Fprintln(w, conf)
	w.Flush()

	return nil
}

func ReadResolveConf() (*DnsConfig, error) {
	cfg, err := ini.Load(ResolvedConfPath)
	if err != nil {
		fmt.Errorf("Fail to read file %s: %v", err)
		return nil, err
	}

	conf := new(DnsConfig)

	dns := cfg.Section("Resolve").Key("DNS").String()
	fallbackDns := cfg.Section("Resolve").Key("FallbackDNS").String()

	conf.DNS = strings.Fields(dns)
	conf.FallbackDNS = strings.Fields(fallbackDns)

	return conf, nil
}

func GetResolveConf(rw http.ResponseWriter) (error) {
	conf, err := ReadResolveConf()
	if err != nil {
		return err
	}

	j, err := json.Marshal(conf)
	if err != nil {
		log.Errorf("Failed to encode json for resolv %s", ResolvedConfPath, err)
		return err
	}

	rw.Write(j)

	return nil
}

func UpdateResolveConf(rw http.ResponseWriter, r *http.Request) (error) {
	dns := DnsConfig{
		DNS: []string{""},
		FallbackDNS: []string{""},
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

	conf, err := ReadResolveConf()
	if err != nil {
		return err
	}

	// update DNS
	for _, s := range dns.DNS {
		if share.StringContains(conf.DNS, s) {
			continue
		}

		conf.DNS = append(conf.DNS, s)
	}

	// update fallback
	for _, s := range dns.FallbackDNS {
		if share.StringContains(conf.FallbackDNS, s) {
			continue
		}
		conf.FallbackDNS = append(conf.FallbackDNS, s)
	}

	err = conf.WriteResolveConfig()
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

func DeleteResolveConf(rw http.ResponseWriter, r *http.Request) (error) {
	dns := DnsConfig{
		DNS: []string{""},
		FallbackDNS: []string{""},
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

	conf, err := ReadResolveConf()
	if err != nil {
		return err
	}

	// update DNS
	for _, s := range dns.DNS {
		if !share.StringContains(conf.DNS, s) {
			continue
		}

		conf.DNS, _ = share.StringDeleteSlice(conf.DNS, s)
	}

	// update Fallback
	for _, s := range dns.FallbackDNS {
		if !share.StringContains(conf.FallbackDNS, s) {
			continue
		}

		conf.FallbackDNS, _ = share.StringDeleteSlice(conf.FallbackDNS, s)
	}

	err = conf.WriteResolveConfig()
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
