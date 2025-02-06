package util

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

func IsClosedNetworkError(err error) bool {
	if err == nil {
		return false
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return true
	}

	var wsErr *websocket.CloseError
	if errors.As(err, &wsErr) && wsErr.Code == websocket.CloseNormalClosure {
		return true

	}

	if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
		return true
	}
	return false
}

func IsClosedServerError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, http.ErrServerClosed)
}

func RespondWithError(w http.ResponseWriter, message string, err error) {
	fmt.Fprintf(w, message)
	if err != nil {
		log.Printf("Error:%v", err)
	}
}
