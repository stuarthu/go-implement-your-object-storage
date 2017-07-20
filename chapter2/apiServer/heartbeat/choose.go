package heartbeat

import (
	"math/rand"
)

func ChooseRandomDataServer() string {
	n := len(DataServers)
	if n == 0 {
		return ""
	}
	i := rand.Intn(n)
	for s, _ := range DataServers {
		if i == 0 {
			return s
		}
		i--
	}
	panic("should not pass here")
}
