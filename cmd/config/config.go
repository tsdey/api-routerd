// SPDX-License-Identifier: Apache-2.0

package config

import (
        "encoding/json"
	"api-routerd/cmd/share" 
        log "github.com/sirupsen/logrus"
        "net/http"
        "strings"
)

const SudoersPath = "/etc/sudoers"
const SSHConfigFile = "/etc/ssh/sshd_config"

type Config struct {
        Action string `json:"path"`
        Path string `json:"path"`
        Property string `json:"property"`
        Value string `json:"value"`
}

//read sudoers file
func GetSudoers(rw http.ResponseWriter, property string) (error) {
        lines, err := share.ReadFullFile(SudoersPath)

        if err != nil {
                log.Fatal("Failed to read %s", SudoersPath)
        }

        if err != nil {
                return err
        }

        line := string("null")

        for _, line = range lines {
                if (strings.Contains(line, "%sudo")) {
                        break
                }
        }

        confproperty := Config {Property: property, Value: line}
        json.NewEncoder(rw).Encode(confproperty)

        return nil
}

//check if sshd_conf exists	
func SSHConfigAvailable(rw http.ResponseWriter, property string) (error) {
	var sshdConfExists = "sshd_config not available"

	if share.PathExists(SSHConfigFile) {
		sshdConfExists = "sshd_config available to read"
	}

        confproperty := Config {Property: property, Value: sshdConfExists}
        json.NewEncoder(rw).Encode(confproperty)

        return nil
}

//read sshd configuration file
func SSHConfFileRead(rw http.ResponseWriter) (error) {
	lines, err := share.ReadFullFile(SSHConfigFile)
	if err != nil {
		log.Fatal("Failed to read: %s", SSHConfigFile)
	}

	//prepare a map of key value pairs from
	//sshd_conf file
	sshdConfMap := make(map[string]string)
	for _, line := range lines {
		fields := strings.Fields(line)

		paramName := fields[0]
		if err != nil {
			continue
		}
		sshdConfMap[paramName] = fields[1]
	}
	
	//encoding the contents of map to
	//JSON format
	j, err := json.Marshal(sshdConfMap)
	if err != nil {
		log.Fatal("there was an issue encoding JSON payload")
	}

	//HHTP Status
	rw.WriteHeader(http.StatusOK)
	rw.Write(j)

	return nil
}
