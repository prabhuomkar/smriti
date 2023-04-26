package cache

import (
	"math"

	"github.com/bluele/gcache"
)

// Init ...
func Init() (gcache.Cache, error) { //nolint: ireturn
	gc := gcache.New(math.MaxInt).
		LRU().
		Build()
	return gc, nil
}
