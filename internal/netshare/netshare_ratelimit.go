// Copyright (C) 2021-2023 Leonid Maslakov.

// This file is part of Lenpaste.

// Lenpaste is free software: you can redistribute it
// and/or modify it under the terms of the
// GNU Affero Public License as published by the
// Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.

// Lenpaste is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero Public License for more details.

// You should have received a copy of the GNU Affero Public License along with Lenpaste.
// If not, see <https://www.gnu.org/licenses/>.

package netshare

import (
	"net"
	"sync"
	"time"
)

type RateLimitSystem struct {
	per5Min  *RateLimit
	per15Min *RateLimit
	per1Hour *RateLimit
}

func NewRateLimitSystem(per5Min, per15Min, per1Hour uint) *RateLimitSystem {
	return &RateLimitSystem{
		per5Min:  NewRateLimit(5*60, per5Min),
		per15Min: NewRateLimit(15*60, per15Min),
		per1Hour: NewRateLimit(60*60, per1Hour),
	}
}

func (rateSys *RateLimitSystem) CheckAndUse(ip net.IP) error {
	var tmp int64

	tmp = rateSys.per5Min.CheckAndUse(ip)
	if tmp != 0 {
		return ErrTooManyRequestsNew(tmp)
	}

	tmp = rateSys.per15Min.CheckAndUse(ip)
	if tmp != 0 {
		return ErrTooManyRequestsNew(tmp)
	}

	tmp = rateSys.per1Hour.CheckAndUse(ip)
	if tmp != 0 {
		return ErrTooManyRequestsNew(tmp)
	}

	return nil
}

type RateLimit struct {
	sync.RWMutex

	limitPeriod int  // N - Rate limit period (in seconds)
	limitCount  uint // X - Max request count per N seconds period

	list map[string]rateLimitIP // Rate limit bucket
}

type rateLimitIP struct {
	UseTime  int64 // Fist IP use time
	UseCount uint  // Requests count by IP
}

func NewRateLimit(rateLimitPeriod int, limitCount uint) *RateLimit {
	rateLimit := &RateLimit{
		limitPeriod: rateLimitPeriod,
		limitCount:  limitCount,
		list:        make(map[string]rateLimitIP),
	}

	go rateLimit.runWorker()

	return rateLimit
}

func (rateLimit *RateLimit) runWorker() {
	for {
		time.Sleep(time.Duration(rateLimit.limitPeriod) * time.Second)

		timeNow := time.Now().Unix()
		rateLimit.Lock()

		for ipStr, data := range rateLimit.list {
			if data.UseTime+int64(rateLimit.limitPeriod) <= timeNow {
				delete(rateLimit.list, ipStr)
			}
		}

		rateLimit.Unlock()
	}
}

func (rateLimit *RateLimit) CheckAndUse(ip net.IP) int64 {
	// If rate limit not need
	if rateLimit.limitCount == 0 {
		return 0
	}

	// Lock
	rateLimit.Lock()
	defer rateLimit.Unlock()

	ipStr := ip.String()
	timeNow := time.Now().Unix()

	// If last use time out
	if rateLimit.list[ipStr].UseTime+int64(rateLimit.limitPeriod) <= timeNow {
		rateLimit.list[ipStr] = rateLimitIP{
			UseTime:  timeNow,
			UseCount: 1,
		}

		return 0

		// Else
	} else {
		if rateLimit.list[ipStr].UseCount < rateLimit.limitCount {
			tmp := rateLimit.list[ipStr]
			tmp.UseCount = tmp.UseCount + 1
			rateLimit.list[ipStr] = tmp
			return 0
		}
	}

	return rateLimit.list[ipStr].UseTime + int64(rateLimit.limitPeriod) - timeNow
}
