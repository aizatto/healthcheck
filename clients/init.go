package clients

import (
	"encoding/gob"
	"math/rand"
	"time"
)

func Init() {
	DotEnv()
	gob.Register(map[string]interface{}{})
	rand.Seed(time.Now().UTC().UnixNano())
}
