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

package storage

import (
	"errors"
	"net"
	"time"

	"git.lcomrade.su/root/lenpaste/internal/model"
)

type rateLimit struct {
	lock chan struct{}

	// category_name: category_data
	m map[model.RLCategory]rateLimitCategory
}

type rateLimitCategory struct {
	limitPeriod time.Duration
	maxReqs     int

	// ip: ip_data
	m map[string]rateLimitIP
}

type rateLimitIP struct {
	startTime time.Time
	count     int
}

func (db *DB) rlNativeInit() {
	// Init data structure
	db.rlNative = &rateLimit{
		lock: make(chan struct{}, 1),
		m:    make(map[model.RLCategory]rateLimitCategory),
	}

	// Run cleanup
	go func() {
		for {
			db.rlNative.lock <- struct{}{}

			now := time.Now()

			for _, cat := range db.rlNative.m {
				for ip, item := range cat.m {
					limitPeriodEnd := item.startTime.Add(cat.limitPeriod)

					if limitPeriodEnd.Before(now) {
						delete(cat.m, ip)
					}
				}
			}

			<-db.rlNative.lock

			timeSub := time.Since(now)
			if timeSub < time.Second {
				time.Sleep(timeSub)
			}
		}
	}()
}

func (rl *rateLimit) writeCategory(name model.RLCategory, limitPeriod time.Duration, maxReqs int) {
	// Lock and unlock
	rl.lock <- struct{}{}
	defer func() { <-rl.lock }()

	// Add category
	rl.m[name] = rateLimitCategory{
		limitPeriod: limitPeriod,
		maxReqs:     maxReqs,
		m:           make(map[string]rateLimitIP),
	}
}

func (rl *rateLimit) check(category model.RLCategory, ip net.IP) error {
	// Lock and unlock
	rl.lock <- struct{}{}
	defer func() { <-rl.lock }()

	// Get category
	cat, ok := rl.m[category]
	if !ok {
		return errors.New("rate limit native: unknown category \"" + category.String() + "\"")
	}

	// Find target IP
	ipStr := ip.String()

	item, ok := cat.m[ipStr]
	if !ok {
		cat.m[ipStr] = rateLimitIP{
			startTime: time.Now(),
			count:     0,
		}
		item = cat.m[ipStr]
	}

	// Check and update rate limits for target IP
	if item.count >= cat.maxReqs {
		limitPeriodEnd := item.startTime.Add(cat.limitPeriod)
		retryAfter := time.Until(limitPeriodEnd)
		return model.ErrTooManyRequestsNew(int64(retryAfter.Seconds()))
	}

	item.count++
	cat.m[ipStr] = item

	return nil
}
