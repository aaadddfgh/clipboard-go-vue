package auth

import (
	"clipboard-go-vue/src/lib"
	"fmt"
	"time"
)

var Instance *Sessions

type ValueInSession struct {
	lib.AESCrypto
	lastUse time.Time
}

// Store every connection's AES key for secure
type Sessions struct {
	seMap map[string]*ValueInSession
}

func (s *Sessions) init() {
	s.seMap = make(map[string]*ValueInSession)

}

func (s *Sessions) put(id string, v *lib.AESCrypto) {
	temp := ValueInSession{
		*v,
		time.Now(),
	}

	s.seMap[id] = &temp
}

func (s *Sessions) Gen(v *lib.AESCrypto) string {
	id := fmt.Sprint(time.Now().UnixNano())
	s.put(id, v)
	return id
}

func (s *Sessions) Exsist(key string) bool {

	defer func() {
		for k, vis := range s.seMap {
			if vis != nil && vis.lastUse.Add(time.Hour*6).After(time.Now()) {
			} else {
				delete(s.seMap, k)
			}
		}
	}()

	return s.seMap[key] != nil && s.seMap[key].lastUse.Add(time.Hour*6).After(time.Now())

}

func (s *Sessions) Get(key string) (aes lib.AESCrypto, exists bool) {
	vv, exists := s.seMap[key]
	//更新最后使用时间
	if exists {
		vv.lastUse = time.Now()
	}
	return vv.AESCrypto, exists

}

func GetInstance() *Sessions {
	if Instance != nil {
		return Instance
	} else {
		Instance = new(Sessions)
		Instance.init()
		return Instance
	}
}
