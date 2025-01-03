package proxyserver

import (
	"fmt"
	"net/http"
	"time"

	"ZTWssProxy/pkg/util"
)

type token struct {
	loginTempID uint32
	accid       uint32
	zoneID      uint32
	gatePort    uint32
	overdueTime int64
}

func (t *token) info() string {
	return fmt.Sprintf("loginTempID:%d,accid:%d,zoneID:%d,gatePort:%d", t.loginTempID, t.accid, t.zoneID, t.gatePort)
}

func createTokenWithRequest(r *http.Request) *token {
	t := &token{}
	t.loginTempID, _ = util.StrToUint32(r.FormValue("loginTempID"))
	t.accid, _ = util.StrToUint32(r.FormValue("accid"))
	t.zoneID, _ = util.StrToUint32(r.FormValue("zoneID"))
	t.gatePort, _ = util.StrToUint32(r.FormValue("port"))
	t.overdueTime = time.Now().Unix() + 10
	return t
}
