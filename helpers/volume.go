package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"

	se "github.com/byuoitav/av-api/statusevaluators"
	"github.com/byuoitav/qsc-microservice/qsysremote"
	"github.com/fatih/color"
)

func Mute(address, name string) (se.MuteStatus, error) {
	return setMuteStatus(address, name, true)
}

func UnMute(address, name string) (se.MuteStatus, error) {
	return setMuteStatus(address, name, false)
}

func setMuteStatus(address, name string, value bool) (se.MuteStatus, error) {
	//we generate our set status request, then we ship it out

	req := qsysremote.GetGenericSetStatusRequest()
	req.Params.Name = name
	if value {
		req.Params.Value = 1
	} else {
		req.Params.Value = 0
	}

	resp, err := qsysremote.SendCommand(address, req)
	if err != nil {
		log.Printf(color.HiRedString("Error: %v", err.Error()))
		return se.MuteStatus{}, err
	}

	//we need to unmarshal our response, parse it for the value we care about, then role with it from there
	val := qsysremote.QSCSetStatusResponse{}
	err = json.Unmarshal(resp, &val)
	if err != nil {
		log.Printf(color.HiRedString("Error: %v", err.Error()))
		return se.MuteStatus{}, err
	}

	//otherwise we check to see what the value is set to
	if val.Result.Name != name {
		errmsg := fmt.Sprintf("Invalid response, the name recieved does not match the name sent %v/%v", name, val.Result.Name)
		log.Printf(color.HiRedString(errmsg))
		return se.MuteStatus{}, errors.New(errmsg)
	}

	if val.Result.Value == 1.0 {
		return se.MuteStatus{Muted: true}, nil
	}
	if val.Result.Value == 0.0 {
		return se.MuteStatus{Muted: false}, nil
	}
	errmsg := fmt.Sprintf("[QSC-Communication] Invalid response received: %v", val.Result)
	log.Printf(color.HiRedString(errmsg))
	return se.MuteStatus{}, errors.New(errmsg)
}

func SetVolume(address, name string, level int) (se.Volume, error) {
	req := qsysremote.GetGenericSetStatusRequest()
	req.Params.Name = name

	if level == 0 {
		req.Params.Value = -100
	} else {
		//do the logrithmic magic
		req.Params.Value = math.Log10(float64(level)/100) * 20
	}

	resp, err := qsysremote.SendCommand(address, req)
	if err != nil {
		log.Printf(color.HiRedString("Error: %v", err.Error()))
		return se.Volume{}, err
	}

	//we need to unmarshal our response, parse it for the value we care about, then role with it from there
	val := qsysremote.QSCSetStatusResponse{}
	err = json.Unmarshal(resp, &val)
	if err != nil {
		log.Printf(color.HiRedString("Error: %v", err.Error()))
		return se.Volume{}, err
	}
	if val.Result.Name != name {
		errmsg := fmt.Sprintf("Invalid response, the name recieved does not match the name sent %v/%v", name, val.Result.Name)
		log.Printf(color.HiRedString(errmsg))
		return se.Volume{}, errors.New(errmsg)
	}

	//reverse it
	toReturn := math.Pow(10, (val.Result.Value/20)) * 100

	return se.Volume{int(toReturn)}, nil
}

func GetVolume(address, name string) (se.Volume, error) {
	return se.Volume{}, nil
}
func GetMute(address, name string) (se.MuteStatus, error) {
	return se.MuteStatus{}, nil

}
