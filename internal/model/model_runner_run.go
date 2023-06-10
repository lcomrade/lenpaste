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

package model

type RunnerRequest struct {
	EnvironmentID string               `json:"environment_id"`
	Build         RunnerRequestBuild   `json:"build"`
	Execute       RunnerRequestExecute `json:"execute"`
}

type RunnerRequestBuild struct {
	File string `json:"file"`
}

type RunnerRequestExecute struct {
	CustomTimeout int               `json:"custom_timeout"`
	Env           map[string]string `json:"env"`
	RunArgs       []string          `json:"run_args"`
	Stdin         string            `json:"stdin"`
}

type RunnerResponse struct {
	TimeoutExit bool   `json:"timeout_exit"`
	Output      string `json:"output"`
	ExitCode    int    `json:"exit_code"`
}
