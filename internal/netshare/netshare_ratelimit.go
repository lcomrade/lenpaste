// Copyright (C) 2021-2022 Leonid Maslakov.

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

const (
	RateLimitPeriod = 5 * 60
)

type RateLimitIP struct {
	UseTime  int64
	UseCount int
}

type RateLimitList struct {
	sync.RWMutex
	m map[string]RateLimitIP
}

type RateLimit struct {
	ReqPer5Minute int

	List RateLimitList
}

func NewRateLimit(reqPer5Minute int) *RateLimit {
	return &RateLimit{
		ReqPer5Minute: reqPer5Minute,
		List: RateLimitList{
			m: make(map[string]RateLimitIP),
		},
	}
}

func (rateLimit *RateLimit) CheckAndUse(ip net.IP) bool {
	// If rate limit not need
	if rateLimit.ReqPer5Minute <= 0 {
		return true
	}

	// Lock
	rateLimit.List.Lock()
	defer rateLimit.List.Unlock()

	ipStr := ip.String()

	// If last use time out
	if rateLimit.List.m[ipStr].UseTime+RateLimitPeriod < time.Now().Unix() {
		rateLimit.List.m[ipStr] = RateLimitIP{
			UseTime:  time.Now().Unix(),
			UseCount: 1,
		}

		return true

		// Else
	} else {
		if rateLimit.List.m[ipStr].UseCount < rateLimit.ReqPer5Minute {
			tmp := rateLimit.List.m[ipStr]
			tmp.UseCount = tmp.UseCount + 1
			rateLimit.List.m[ipStr] = tmp
			return true
		}
	}

	return false
}
