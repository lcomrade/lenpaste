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

package runner

import (
	"errors"
	"net/url"

	"git.lcomrade.su/root/lenpaste/internal/config"
	"git.lcomrade.su/root/lenpaste/internal/logger"
	"git.lcomrade.su/root/lenpaste/internal/model"
)

type runner struct {
	workers_busy int
	info         model.RunnerInfo
}

type RunCoordinator struct {
	cfg *config.Config

	lock    chan struct{}
	runners map[string]runner
}

func New(cfg *config.Config, log *logger.Logger) (*RunCoordinator, error) {
	coord := &RunCoordinator{
		cfg:     cfg,
		lock:    make(chan struct{}, 1),
		runners: map[string]runner{},
	}

	for _, part := range cfg.CodeRun.Runners {
		// Check runner URL
		{
			u, err := url.Parse(part.BaseURL)
			if err != nil {
				return nil, errors.New("runner: bad URL \"" + part.BaseURL + "\"")
			}

			goodURL := u.Scheme + "://" + u.Host + u.Path
			if part.BaseURL != goodURL {
				return nil, errors.New("runner: base URL contains unnecessary elements, should be \"" + goodURL + "\", not \"" + part.BaseURL + "\"")
			}
		}

		// If runner already added
		_, exists := coord.runners[part.BaseURL]
		if exists {
			return nil, errors.New("runner: duplicate runner URL \"" + part.BaseURL + "\"")
		}

		// Get information about runner
		//
		// TODO: Add asynchronous receiving of information about Runners.
		info, err := getInfo(part.BaseURL, model.RunnerTimeout, part.SharedSecret)
		if err != nil {
			if part.Required {
				return nil, errors.New("failed to connect to the required Runner \"" + part.BaseURL + "\"")

			} else {
				log.Warning("Failed to connect to Runner \"" + part.BaseURL + "\"")
				continue
			}
		}

		coord.runners[part.BaseURL] = runner{
			workers_busy: 0,
			info:         info,
		}
	}

	return coord, nil
}
