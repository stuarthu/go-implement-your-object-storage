package heartbeat

import (
	"math/rand"
)

func ChooseRandomDataServer() string {
	n := len(dataServers)
	if n == 0 {
		return ""
	}
	i := rand.Intn(n)
	for s, _ := range dataServers {
        if i == 0 {
            return s
        }
        i--
	}
	panic("should not pass here")
}
