package proxyserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"ZTWssProxy/configs"
)

type Token struct {
	LoginTempID uint32 `json:"loginTempID"`
	AccID       uint32 `json:"accid"`
	ZoneID      uint32 `json:"zoneID"`
	GateIp      string `json:"GateIp"`
	GatePort    uint32 `json:"gatePort"`
	expireTime  time.Time
}

func (t *Token) check() error {
	if t.LoginTempID == 0 || t.AccID == 0 || t.ZoneID == 0 || len(t.GateIp) == 0 || t.GatePort == 0 {
		return errors.New("Invalid token: " + t.info())
	}
	return nil
}

func (t *Token) info() string {
	return fmt.Sprintf("loginTempID:%d,accid:%d,zoneID:%d,gateIp:%s,gatePort:%d", t.LoginTempID, t.AccID, t.ZoneID, t.GateIp, t.GatePort)
}

func (t *Token) isExpired() bool {
	return time.Now().After(t.expireTime)
}

type TokenManager struct {
	tokens sync.Map
}

func (tm *TokenManager) createTokenWithRequest(r *http.Request) (*Token, error) {
	t := &Token{}
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, err
	}
	defer r.Body.Close()
	t.expireTime = time.Now().Add(configs.TokenExpireTime * time.Second)
	return t, nil
}

func (tm *TokenManager) add(t *Token) {
	tm.tokens.Store(t.LoginTempID, t)
	log.Printf("Add token:%v", t.info())
}

func (tm *TokenManager) get(loginTempID uint32) (*Token, bool) {
	val, ok := tm.tokens.Load(loginTempID)
	if ok {
		return val.(*Token), true
	}
	return nil, false
}

func (tm *TokenManager) delete(loginTempID uint32) {
	tm.tokens.Delete(loginTempID)
}

func (tm *TokenManager) cleanExpiredTokens() {
	tm.tokens.Range(func(key, value interface{}) bool {
		t, ok := value.(*Token)
		if ok && t.isExpired() {
			tm.tokens.Delete(key)
			log.Printf("Token expired:%v", t.LoginTempID)
		}
		return true
	})
}
