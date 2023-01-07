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

document.addEventListener("DOMContentLoaded", () => {
	const shortWeekDay = [{{call .Translate `pasteJS.ShortWeekDay`}}];
	const shortMonth = [{{call .Translate `pasteJS.ShortMonth`}}];

	function dateToString(date) {
		let dateStr = shortWeekDay[date.getDay()] + ", " + date.getDate() + " " + shortMonth[date.getMonth()];
		dateStr = dateStr + " " + date.getFullYear();
		dateStr = dateStr + " " + date.getHours() + ":" + date.getMinutes() + ":" + date.getSeconds();
		
		let tz = date.getTimezoneOffset() / 60 * -1;
		if (tz >= 0) {
			dateStr = dateStr + " +" + tz;
			
		} else {
			dateStr = dateStr + " " + tz;
		}
		
		return dateStr;
	}

	let createTime = document.getElementById("createTime");
	createTime.textContent = dateToString(new Date(createTime.textContent));


	let deleteTime = document.getElementById("deleteTime");
	if (deleteTime != null) {
		deleteTime.textContent = dateToString(new Date(deleteTime.textContent));
	}
});
