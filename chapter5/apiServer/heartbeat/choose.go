package heartbeat

import (
	"math/rand"
)

func ChooseRandomDataServers(n int, exclude map[int]string) (ds []string) {
	candidates := make([]string, 0)
	reverseExcludeMap := make(map[string]int)
	for id, addr := range exclude {
		reverseExcludeMap[addr] = id
	}
	for s, _ := range dataServers {
		_, excluded := reverseExcludeMap[s]
		if !excluded {
			candidates = append(candidates, s)
		}
	}
	if len(candidates) < n {
		return
	}
	p := rand.Perm(n)
	for i := range p {
		ds = append(ds, candidates[i])
	}
	return
}
