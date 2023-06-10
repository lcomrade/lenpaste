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
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"git.lcomrade.su/root/lenpaste/internal/model"
)

var client = http.Client{
	Timeout: 60 * time.Second,
}

func getInfo(baseURL string, timeout int, secret string) (model.RunnerInfo, error) {
	// Prepare response
	p, err := url.JoinPath(baseURL, "/v1/info")
	if err != nil {
		return model.RunnerInfo{}, errors.New("get runner info: " + err.Error())
	}

	req, err := http.NewRequest(http.MethodGet, p, nil)
	if err != nil {
		return model.RunnerInfo{}, errors.New("get runner info: " + err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+secret)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return model.RunnerInfo{}, errors.New("get runner info: " + err.Error())
	}
	defer resp.Body.Close()

	// Parse response
	var out model.RunnerInfo
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return model.RunnerInfo{}, errors.New("get runner info: " + err.Error())
	}

	return out, nil
}

func postRun(baseURL string, timeout int, secret string, runReq model.RunnerRequest) (model.RunnerResponse, error) {
	// Prepare response
	p, err := url.JoinPath(baseURL, "/v1/run")
	if err != nil {
		return model.RunnerResponse{}, errors.New("run code in runner: " + err.Error())
	}

	body, err := json.Marshal(&runReq)
	if err != nil {
		return model.RunnerResponse{}, errors.New("run code in runner: " + err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, p, bytes.NewBuffer(body))
	if err != nil {
		return model.RunnerResponse{}, errors.New("run code in runner: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+secret)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return model.RunnerResponse{}, errors.New("run code in runner: " + err.Error())
	}
	defer resp.Body.Close()

	// Parse response
	var out model.RunnerResponse
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return model.RunnerResponse{}, errors.New("run code in runner: " + err.Error())
	}

	return out, nil
}
