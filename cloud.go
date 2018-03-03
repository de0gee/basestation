package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/cihub/seelog"
	cloud "github.com/de0gee/de0gee-cloud/src"
)

func uploadToServer(payloadBytes []byte, endpoint string) (target cloud.ServerResponse, err error) {
	log.Debugf("%s", payloadBytes)

	body := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest("POST", CloudServer+"/"+endpoint, body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&target)
	log.Debugf("response from /%s: %+v", endpoint, target)
	if err == nil {
		if !target.Success {
			err = fmt.Errorf(target.Message)
		}
	}
	return
}
