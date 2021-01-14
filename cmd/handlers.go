package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/byuoitav/common/status"
	qsc "github.com/byuoitav/qsc"
	"github.com/labstack/echo"
)

type Handlers struct {
	CreateDSP func(string) *qsc.DSP
}

func (h *Handlers) RegisterRoutes(group *echo.Group) {
	qsc := group.Group("")

	qsc.GET("/:address/:name/volume/level", func(c echo.Context) error {
		addr := c.Param("address")
		name := c.Param("name")
		name += "Gain"
		dsp := h.CreateDSP(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("getting volumes")

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		vols, err := dsp.Volumes(ctx, []string{name})
		if err != nil {
			l.Printf("unable to get volumes: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Got volumes: %+v", vols)

		vol, ok := vols[name]
		if !ok {
			l.Printf("invalid name %q requested", name)
			return c.String(http.StatusBadRequest, "invalid name")
		}

		return c.JSON(http.StatusOK, status.Volume{
			Volume: vol,
		})
	})

	qsc.GET("/:address/:name/mute/status", func(c echo.Context) error {
		addr := c.Param("address")
		name := c.Param("name")
		name += "Mute"
		dsp := h.CreateDSP(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("getting mutes")

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		mutes, err := dsp.Mutes(ctx, []string{name})
		if err != nil {
			l.Printf("unable to get mutes: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Got mutes: %+v", mutes)

		mute, ok := mutes[name]
		if !ok {
			l.Printf("invalid name %q requested", name)
			return c.String(http.StatusBadRequest, "invalid name")
		}

		return c.JSON(http.StatusOK, status.Mute{
			Muted: mute,
		})
	})

	qsc.GET("/:address/:name/volume/set/:volume", func(c echo.Context) error {
		addr := c.Param("address")
		name := c.Param("name")
		name += "Gain"
		dsp := h.CreateDSP(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		vol, err := strconv.Atoi(c.Param("volume"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		l.Printf("setting volume on %q to %d", name, vol)

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		err = dsp.SetVolume(ctx, name, vol)
		if err != nil {
			l.Printf("unable to set volume: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Set volume")
		return c.JSON(http.StatusOK, status.Volume{
			Volume: vol,
		})
	})

	qsc.GET("/:address/:name/volume/mute", func(c echo.Context) error {
		addr := c.Param("address")
		name := c.Param("name")
		name += "Mute"
		dsp := h.CreateDSP(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("setting mute on %q to true", name)

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		err := dsp.SetMute(ctx, name, true)
		if err != nil {
			l.Printf("unable to mute: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Set mute")
		return c.JSON(http.StatusOK, status.Mute{
			Muted: true,
		})
	})

	qsc.GET("/:address/:name/volume/unmute", func(c echo.Context) error {
		addr := c.Param("address")
		name := c.Param("name")
		name += "Mute"
		dsp := h.CreateDSP(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("setting mute on %q to false", name)

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		err := dsp.SetMute(ctx, name, false)
		if err != nil {
			l.Printf("unable to unmute: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Set mute")
		return c.JSON(http.StatusOK, status.Mute{
			Muted: false,
		})
	})
}
