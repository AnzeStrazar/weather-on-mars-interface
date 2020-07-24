package cache

import (
	"sync"
	"weather-on-mars-interface/types"
)

type Cache struct {
	SolCache map[string]types.Sol
	SolMutex sync.RWMutex
}
