package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"

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
	log.Printf("got: %v", level)
	req := qsysremote.GetGenericSetStatusRequest()
	req.Params.Name = name

	if level == 0 {
		req.Params.Value = -100
	} else {
		//do the logrithmic magic
		req.Params.Value = VolToDb(level)
	}
	log.Printf("sending: %v", req.Params.Value)

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

	return se.Volume{DbToVolumeLevel(val.Result.Value)}, nil
}

func DbToVolumeLevel(level float64) int {
	return int(math.Pow(10, (level/20)) * 100)
}

func VolToDb(level int) float64 {
	return math.Log10(float64(level)/100) * 20
}

func GetVolume(address, name string) (se.Volume, error) {

	resp, err := GetControlStatus(address, name)
	if err != nil {
		log.Printf(color.HiRedString("There was an error: %v", err.Error()))
		return se.Volume{}, err
	}

	log.Printf(color.HiBlueString("[QSC-Communication] Response received: %+v", resp))

	//get the volume out of the dsp and run it through our equation to reverse it
	for _, res := range resp.Result {
		if res.Name == name {
			return se.Volume{DbToVolumeLevel(res.Value)}, nil
		}
	}

	return se.Volume{}, errors.New("[QSC-Communication] No value returned with the name matching the requested state")
}
func GetMute(address, name string) (se.MuteStatus, error) {
	resp, err := GetControlStatus(address, name)
	if err != nil {
		log.Printf(color.HiRedString("There was an error: %v", err.Error()))
		return se.MuteStatus{}, err
	}

	//get the volume out of the dsp and run it through our equation to reverse it
	for _, res := range resp.Result {
		if res.Name == name {
			if res.Value == 1.0 {
				return se.MuteStatus{Muted: true}, nil
			}
			if res.Value == 0.0 {
				return se.MuteStatus{Muted: false}, nil
			}
		}
	}
	errmsg := "[QSC-Communication] No value returned with the name matching the requested state"
	log.Printf(color.HiRedString(errmsg))
	return se.MuteStatus{}, errors.New(errmsg)
}

func GetControlStatus(address, name string) (qsysremote.QSCGetStatusResponse, error) {
	req := qsysremote.GetGenericGetStatusRequest()
	req.Params = append(req.Params, name)

	toReturn := qsysremote.QSCGetStatusResponse{}

	resp, err := qsysremote.SendCommand(address, req)
	if err != nil {
		log.Printf(color.HiRedString(err.Error()))
		return toReturn, err
	}

	err = json.Unmarshal(resp, &toReturn)
	if err != nil {
		log.Printf(color.HiRedString(err.Error()))
	}

	return toReturn, err
}

func SetControlStatus(address, name, value string) (qsysremote.QSCSetStatusResponse, error) {
	var err error
	req := qsysremote.GetGenericSetStatusRequest()
	val := qsysremote.QSCSetStatusResponse{}

	req.Params.Name = name
	req.Params.Value, err = strconv.ParseFloat(value, 64)
	if err != nil {
		return val, errors.New("Invalid value, must be a float")
	}
	log.Printf("sending: %v:%v to %v", req.Params.Name, req.Params.Value, address)

	resp, err := qsysremote.SendCommand(address, req)
	if err != nil {
		log.Printf(color.HiRedString("Error: %v", err.Error()))
		return val, err
	}

	//we need to unmarshal our response, parse it for the value we care about, then role with it from there
	err = json.Unmarshal(resp, &val)
	if err != nil {
		log.Printf(color.HiRedString("Error: %v", err.Error()))
		return val, err
	}
	if val.Result.Name != name {
		errmsg := fmt.Sprintf("Invalid response, the name recieved does not match the name sent %v/%v", name, val.Result.Name)
		log.Printf(color.HiRedString(errmsg))
		return val, errors.New(errmsg)
	}

	return val, nil
}
