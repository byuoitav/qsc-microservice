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
	qsc := group.Group("/QSC/:address")

	qsc.GET("/block/:block/volume", func(c echo.Context) error {
		addr := c.Param("address")
		block := c.Param("block")
		dsp := h.CreateDSP(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("getting volumes")

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		vols, err := dsp.Volumes(ctx, []string{block})
		if err != nil {
			l.Printf("unable to get volumes: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Got volumes: %+v", vols)

		vol, ok := vols[block]
		if !ok {
			l.Printf("invalid block %q requested", block)
			return c.String(http.StatusBadRequest, "invalid block")
		}

		return c.JSON(http.StatusOK, status.Volume{
			Volume: vol,
		})
	})

	qsc.GET("/block/:block/mute", func(c echo.Context) error {
		addr := c.Param("address")
		block := c.Param("block")
		dsp := h.CreateDSP(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("getting mutes")

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		mutes, err := dsp.Mutes(ctx, []string{block})
		if err != nil {
			l.Printf("unable to get mutes: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Got mutes: %+v", mutes)

		mute, ok := mutes[block]
		if !ok {
			l.Printf("invalid block %q requested", block)
			return c.String(http.StatusBadRequest, "invalid block")
		}

		return c.JSON(http.StatusOK, status.Mute{
			Muted: mute,
		})
	})

	qsc.GET("/block/:block/volume/:volume", func(c echo.Context) error {
		addr := c.Param("address")
		block := c.Param("block")
		dsp := h.CreateDSP(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		vol, err := strconv.Atoi(c.Param("volume"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		l.Printf("setting volume on %q to %d", block, vol)

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		err = dsp.SetVolume(ctx, block, vol)
		if err != nil {
			l.Printf("unable to set volume: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Set volume")
		return c.JSON(http.StatusOK, status.Volume{
			Volume: vol,
		})
	})

	qsc.GET("/block/:block/muted/:mute", func(c echo.Context) error {
		addr := c.Param("address")
		block := c.Param("block")
		dsp := h.CreateDSP(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		mute, err := strconv.ParseBool(c.Param("mute"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		l.Printf("setting mute on %q to %v", block, mute)

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		err = dsp.SetMute(ctx, block, mute)
		if err != nil {
			l.Printf("unable to set volume: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Set mute")
		return c.JSON(http.StatusOK, status.Mute{
			Muted: mute,
		})
	})
}
