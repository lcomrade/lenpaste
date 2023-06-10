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

type RunnerInfo struct {
	Provider string `json:"provider"`

	Workers int `json:"workers"`

	Environments map[string]RunnerEnvironment `json:"environments"`
}

type RunnerEnvironment struct {
	CpuCores int `json:"cpu_cores"`
	RAM      int `json:"ram"`

	Arch string `json:"arch"`
	OS   string `json:"os"`

	Language string `json:"language"`
	ToolName string `json:"tool_name"`

	Timeout int `json:"timeout"`
}
