package floodcontrol

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type checkFlood struct {
	pastTime time.Duration
	counter  int
	redis    *redis.Client
	mutex    sync.Mutex
}

func (cF *checkFlood) Check(ctx context.Context, userID int64) (bool, error) {
	cF.mutex.Lock()
	defer cF.mutex.Unlock()

	now := time.Now()
	userid := strconv.FormatInt(userID, 10)

	pTime, err := cF.redis.Get(ctx, userid).Int64()
	if err != nil {
		return false, err
	}
	timeCheck := now.Sub(time.Unix(pTime, 0))

	counter, err := cF.redis.Incr(ctx, userid).Result()
	if err != nil {
		return false, err
	}

	if cF.pastTime <= timeCheck {
		return int(counter) <= cF.counter, nil
	}

	return true, nil
}

func CheckFloodConstructor(pastTime time.Duration, counter int, redis *redis.Client) *checkFlood {
	return &checkFlood{
		pastTime: pastTime,
		counter:  counter,
		redis:    redis,
		mutex:    sync.Mutex{},
	}
}
