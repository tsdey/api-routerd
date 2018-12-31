// SPDX-License-Identifier: Apache-2.0

package system

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	JournalConfPath = "/etc/systemd/journald.conf"
)

var JournalConfig = map[string]string{
	"Storage"              : "",
	"Compress"             : "",
	"Seal"                 : "",
	"SplitMode"            : "",
	"SyncIntervalSec"      : "",
	"RateLimitIntervalSec" : "",
	"RateLimitBurst"       : "",
	"SystemMaxUse"         : "",
	"SystemKeepFree"       : "",
	"SystemMaxFileSize"    : "",
	"SystemMaxFiles"       : "",
	"RuntimeMaxUse"        : "",
	"RuntimeKeepFree"      : "",
	"RuntimeMaxFileSize"   : "",
	"RuntimeMaxFiles"      : "",
	"MaxRetentionSec"      : "",
	"MaxFileSec"           : "",
	"ForwardToSyslog"      : "",
	"ForwardToKMsg"        : "",
	"ForwardToConsole"     : "",
	"ForwardToWall"        : "",
	"TTYPath"              : "",
	"MaxLevelStore"        : "",
	"MaxLevelSyslog"       : "",
	"MaxLevelKMsg"         : "",
	"MaxLevelConsole"      : "",
	"MaxLevelWall"         : "",
	"LineMax"              : "",
	"ReadKMsg"             : "",
}

func WriteJournalConfig() (error) {
	f, err := os.OpenFile(JournalConfPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	conf :="[Journal]\n"
	for k, v := range JournalConfig {
		if v != "" {
			conf += k + "=" + v
		} else {
			conf += "#" + k + "="
		}
		conf += "\n"
	}

	fmt.Fprintln(w, conf)
	w.Flush()

	return nil
}

func ReadJournalConf() (error) {
	cfg, err := ini.Load(JournalConfPath)
	if err != nil {
		fmt.Errorf("Fail to read file %s: %v", err)
		return err
	}

	for k, _ := range JournalConfig {
		JournalConfig[k] = cfg.Section("Journal").Key(k).String()
	}

	return nil
}

func GetJournalConf(rw http.ResponseWriter) (error) {
	err := ReadJournalConf()
	if err != nil {
		return err
	}

	j, err := json.Marshal(JournalConfig)
	if err != nil {
		log.Errorf("Failed to encode json for Journal %s", JournalConfPath, err)
		return err
	}

	rw.Write(j)

	return nil
}

func UpdateJournalConf(rw http.ResponseWriter, r *http.Request) (error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("Failed to parse HTTP request: ", err)
		return err
	}

	conf := make(map[string]string)
	err = json.Unmarshal([]byte(body), &conf)
	if err != nil {
		log.Error("Failed to Decode HTTP request to json: ", err)
		return err
	}

	err = ReadJournalConf()
	if err != nil {
		return err
	}

	for k, v := range conf {
		_, ok := JournalConfig[k]
		if ok {
			JournalConfig[k] = v
		}
	}

	err = WriteJournalConfig()
	if err != nil {
		log.Errorf("Failed Write to journal conf: %s", err)
		return err
	}

	j, err := json.Marshal(JournalConfig)
	if err != nil {
		log.Errorf("Failed to encode json for system conf %s", err)
		return err
	}

	rw.Write(j)

	return nil
}
