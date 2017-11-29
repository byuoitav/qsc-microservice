package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/byuoitav/qsc-microservice/helpers"
	"github.com/byuoitav/qsc-microservice/qsysremote"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func Mute(context echo.Context) error {
	address := context.Param("address")
	name := context.Param("name")
	name = name + "Mute"

	status, err := helpers.Mute(address, name)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}
	return context.JSON(http.StatusOK, status)
}

func UnMute(context echo.Context) error {
	address := context.Param("address")
	name := context.Param("name")
	name = name + "Mute"

	status, err := helpers.UnMute(address, name)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}
	return context.JSON(http.StatusOK, status)

}

func SetVolume(context echo.Context) error {
	address := context.Param("address")
	name := context.Param("name")
	levelstr := context.Param("level")
	level, err := strconv.Atoi(levelstr)
	name = name + "Gain"

	if err != nil {
		errmsg := fmt.Sprintf("%v is not a valid parameter for level. Must be a valid float", levelstr)
		log.Printf(color.HiRedString(errmsg))
		return context.JSON(http.StatusBadRequest, errors.New(errmsg))
	}

	status, err := helpers.SetVolume(address, name, level)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}
	return context.JSON(http.StatusOK, status)
}

func GetVolume(context echo.Context) error {
	address := context.Param("address")
	name := context.Param("name")
	name = name + "Gain"

	status, err := helpers.GetVolume(address, name)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}
	return context.JSON(http.StatusOK, status)
}

func GetMute(context echo.Context) error {
	address := context.Param("address")
	name := context.Param("name")
	name = name + "Mute"

	status, err := helpers.GetMute(address, name)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}
	return context.JSON(http.StatusOK, status)
}

func Test(context echo.Context) error {
	toBind := qsysremote.QSCStatusReport{}
	err := context.Bind(&toBind)
	if err != nil {
		log.Printf(color.HiRedString(err.Error()))
		return context.JSON(http.StatusBadRequest, err.Error())
	}
	return context.JSON(http.StatusOK, toBind)
}
