package qsysremote

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/fatih/color"
)

const TIMEOUT_IN_SECONDS = 2.0

func SendCommand(address string, request interface{}) ([]byte, error) {
	log.Printf(color.HiBlueString("[QSC-Communication] Sending a requst to %v", address))
	toSend, err := json.Marshal(request)
	if err != nil {
		errmsg := fmt.Sprintf("[QSC-Communication] Invalid request, could not marshal: %v", err.Error())
		log.Printf(color.HiRedString(errmsg))
		return []byte{}, errors.New(errmsg)
	}

	conn, err := getConnection(address, "1710")
	if err != nil {
		return []byte{}, err
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)

	conn.SetReadDeadline(time.Now().Add(time.Duration(TIMEOUT_IN_SECONDS) * time.Second))
	msg, err := reader.ReadBytes('\x00')
	if err != nil {
		log.Printf(color.HiRedString("Error: %v", err.Error()))
		return []byte{}, err
	}

	//we can validate that the message is what we think it should be
	report := QSCStatusReport{}
	msg = bytes.Trim(msg, "\x00")

	err = json.Unmarshal(msg, &report)
	if err != nil {
		errmsg := fmt.Sprintf("[QSC-Communication] bad state recieved from device %v on connection: %s. Error: %v", address, msg, err.Error())
		log.Printf(color.HiRedString(errmsg))
		return []byte{}, errors.New(errmsg)
	}

	//now we write our command
	conn.Write(toSend)
	conn.Write([]byte{0x00})

	conn.SetReadDeadline(time.Now().Add(time.Duration(TIMEOUT_IN_SECONDS) * time.Second))
	msg, err = reader.ReadBytes('\x00')
	if err != nil {
		log.Printf(color.HiRedString("Error: %v", err.Error()))
		return []byte{}, err
	}
	msg = bytes.Trim(msg, "\x00")
	log.Printf(color.HiBlueString("[QSC-Communication] Done with request to %v.", address))
	return msg, nil
}

func getConnection(address string, port string) (*net.TCPConn, error) {
	addr, err := net.ResolveTCPAddr("tcp", address+":"+port)
	if err != nil {
		log.Printf(color.HiRedString("Error: %v", err.Error()))
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Printf(color.HiRedString("Error: %v", err.Error()))
	}
	return conn, err
}
