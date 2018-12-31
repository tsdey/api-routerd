// SPDX-License-Identifier: Apache-2.0

package config

import (
        "encoding/json"
        "github.com/gorilla/mux"
        "net/http"
)

//function to read sudoers rule
func ReadConfigSudoers(rw http.ResponseWriter, req *http.Request) {
        config := new(Config)
        _= json.NewDecoder(req.Body).Decode(&config);
        GetSudoers(rw, config.Property)
}

//function check if sshd configuration exists
func ReadSSHConfigAvailable(rw http.ResponseWriter, req *http.Request) {
        config := new(Config)
        _= json.NewDecoder(req.Body).Decode(&config);
        SSHConfigAvailable(rw, config.Property)
}

//function to read sshd configuration file
func ReadSSHConfigRead(rw http.ResponseWriter, req *http.Request) {
        config := new(Config)
        _= json.NewDecoder(req.Body).Decode(&config);
        SSHConfFileRead(rw)
}

//register router to read configuration files
//for various services
func RegisterRouterConfig(router *mux.Router) {
        s := router.PathPrefix("/config").Subrouter()
        s.HandleFunc("/sudoers", ReadConfigSudoers)
        s.HandleFunc("/sshd", ReadSSHConfigAvailable)
        s.HandleFunc("/sshd/sshdconf", ReadSSHConfigRead)
}

