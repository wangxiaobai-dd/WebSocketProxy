package proxyserver

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"ZTWssProxy/pkg/util"
)

type Token struct {
	loginTempID uint32
	accid       uint32
	zoneID      uint32
	gatePort    uint32
	overdueTime int64
}

func (t *Token) info() string {
	return fmt.Sprintf("loginTempID:%d,accid:%d,zoneID:%d,gatePort:%d", t.loginTempID, t.accid, t.zoneID, t.gatePort)
}

type TokenManager struct {
	tokens sync.Map
}

func (tm *TokenManager) createTokenWithRequest(r *http.Request) *Token {
	t := &Token{}
	t.loginTempID, _ = util.StrToUint32(r.FormValue("loginTempID"))
	t.accid, _ = util.StrToUint32(r.FormValue("accid"))
	t.zoneID, _ = util.StrToUint32(r.FormValue("zoneID"))
	t.gatePort, _ = util.StrToUint32(r.FormValue("port"))
	t.overdueTime = time.Now().Unix() + 10
	return t
}

func (tm *TokenManager) Add(t *Token) {
	tm.tokens.Store(t.loginTempID, t)
}

func (tm *TokenManager) Get(loginTempID uint64) (Token, bool) {
	val, ok := tm.tokens.Load(loginTempID)
	if ok {
		return val.(Token), true
	}
	return Token{}, false
}

func (tm *TokenManager) Delete(loginTempID uint64) {
	tm.tokens.Delete(loginTempID)
}
