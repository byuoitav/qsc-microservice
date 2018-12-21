package helpers

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/byuoitav/common/log"
	"github.com/byuoitav/common/nerr"
	"github.com/byuoitav/common/structs"
	"github.com/byuoitav/qsc-microservice/qsysremote"
	"github.com/fatih/color"
)

// GetDetails is all the juicy details about the QSC that everyone is DYING to know about
func GetDetails(address string) (structs.HardwareInfo, *nerr.E) {

	// toReturn is the struct of Hardware info
	var details structs.HardwareInfo

	// get the hostname
	addr, e := net.LookupAddr(address)
	if e != nil {
		details.Hostname = address
	} else {
		details.Hostname = strings.Trim(addr[0], ".")
	}

	resp, err := GetStatus(address)
	if err != nil {
		return details, nerr.Translate(err).Addf("There was an error getting the status")
	}

	log.L.Infof("response: %v", resp)
	details.ModelName = resp.Result.Platform
	details.PowerStatus = resp.Result.State

	details.NetworkInfo = structs.NetworkInfo{
		IPAddress: address,
	}

	return details, nil
}

// GetStatus will be getting responses for us I hope...
func GetStatus(address string) (qsysremote.QSCStatusGetResponse, error) {
	req := qsysremote.GetGenericStatusGetRequest()

	log.L.Infof("In GetStatus...")
	toReturn := qsysremote.QSCStatusGetResponse{}

	resp, err := qsysremote.SendCommand(address, req)
	if err != nil {
		log.L.Info(color.HiRedString(err.Error()))
		return toReturn, err
	}

	err = json.Unmarshal(resp, &toReturn)
	if err != nil {
		log.L.Infof(color.HiRedString(err.Error()))
	}

	return toReturn, err
}
