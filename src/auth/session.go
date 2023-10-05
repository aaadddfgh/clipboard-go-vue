package auth

import (
	"clipboard-go-vue/src/lib"
	"fmt"
	"time"
)

var Instance *Sessions

type ValInSession struct {
	lib.AESCrypto
	lastUse time.Time
}

type Sessions struct {
	seMap map[string]*ValInSession
}

func (s *Sessions) init() {
	s.seMap = make(map[string]*ValInSession)

}

func (s *Sessions) put(id string, v *lib.AESCrypto) {
	temp := ValInSession{
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
	vv, err := s.seMap[key]
	return vv.AESCrypto, err

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
