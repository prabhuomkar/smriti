package cache

import (
	"math"

	"github.com/bluele/gcache"
)

// nolint:ireturn
// Init ...
func Init() (gcache.Cache, error) {
	gc := gcache.New(math.MaxInt).
		LRU().
		Build()
	return gc, nil
}
