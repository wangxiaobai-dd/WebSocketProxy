package util

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
)

func IsClosedNetworkError(err error) bool {
	if err == nil {
		return false
	}
	var opError *net.OpError
	if errors.As(err, &opError) {
		return true
	}
	//var errno syscall.Errno
	//if errors.As(err, &errno) {
	//	if errors.Is(errno, syscall.EPIPE) || errors.Is(errno, syscall.ECONNRESET) {
	//		return true
	//	}
	//}
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
