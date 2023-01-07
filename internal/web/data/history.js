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

function isLocalStoageSupported() {
	if (typeof localStorage === 'object') {
		try {
			localStorage.setItem("__local_storage_test__", 1);
			localStorage.removeItem("__local_storage_test__");
			return true;
		} catch (e) {
			return false;
		}
	}

	return false;
}

function historyRefreshList() {
	const shortMonth = [{{call .Translate `pasteJS.ShortMonth`}}];

	// Get and clean list
	let listElement = document.getElementById("js-history-popup-list");
	listElement.innerHTML = "";

	// Read locale storage
	let timeNowUnix = Math.floor(Date.now() / 1000)
	let historyJSON = localStorage.getItem("history");
	if (historyJSON != null) {
		let history = JSON.parse(historyJSON);

		for (let i = 0; history.length > i; i++) {
			// Check title
			let title = history[i].title;
			if (title == "") {
				title = "{{ call .Translate `historyJS.Untitled` }}";
			}

			// Convert create date to string
			let date = new Date(history[i].createTime * 1000)

			let dateDayStr = date.getDate();
			if (date.getDate() < 10) {
				dateDayStr = "0" + dateDayStr;
			}
			let dateStr = dateDayStr + " " + shortMonth[date.getMonth()] + ", " + date.getFullYear();

			// Add row
			if (timeNowUnix < history[i].deleteTime || history[i].deleteTime == 0) {
				listElement.insertAdjacentHTML("beforeend", "<li>[" + dateStr + "] <a href='/"+history[i].id+"'>"+title+"</a></li>");
			} else {
				listElement.insertAdjacentHTML("beforeend", "<li><del>[" + dateStr + "] <a class='text-grey' href='/"+history[i].id+"'>"+title+"</a></del></li>");
			}
		}
	}
}

function historyPopUpShow() {
	document.getElementById("js-history-popup").style.visibility = "visible";
	document.getElementById("js-history-popup-background").style.visibility = "visible";
	document.getElementsByTagName("body")[0].style.overflow = "hidden";
	document.addEventListener("keydown", historyPopUpEscEvent);
	historyRefreshList();
}

function historyPopUpHide() {
	document.getElementById("js-history-popup").style.visibility = "hidden";
	document.getElementById("js-history-popup-background").style.visibility = "hidden";
	document.getElementsByTagName("body")[0].style.overflow = null;
	document.removeEventListener("keydown", historyPopUpEscEvent);
}

function historyPopUpEscEvent(event) {
	// If ESC pressed
	if (event.keyCode == 27) {
		historyPopUpHide();
	}
}

function historyEnable() {
	if (document.getElementById("js-history-popup-enable").checked == true) {
		localStorage.removeItem("DisableHistory");
		alert("{{ call .Translate `historyJS.HistoryEnabledAlert` }}");
	} else {
		localStorage.setItem("DisableHistory", true);
		alert("{{ call .Translate `historyJS.HistoryDisabledAlert` }}");
	}
}

function historyClear() {
	if (window.confirm("{{ call .Translate `historyJS.ClearHistoryConfirm` }}")) {
		localStorage.removeItem("history");
		historyRefreshList();
	}
}

document.addEventListener("DOMContentLoaded", () => {
	// Edit CSS
	let newStyleSheet = `
#js-history-button:hover {
	cursor: pointer;
}

#js-history-popup {
	background: {{ call .Theme `color.Article` }};
	padding: 20px;

	position: fixed;
	z-index: 2;
	top: 15%;
	bottom: 15%;
	right: 20%;
	left: 20%;

	overflow: hidden;
}

@media all and (max-device-width: 640px), all and (orientation: portrait) {
	#js-history-popup {
		top: 5%;
		bottom: 5%;
		right: 0;
		left: 0;
	}
}

#js-history-popup-background {
	width: 100%;
	height: 100%;
	position: fixed;
	top: 0;
	left: 0;
	z-index: 1;
	background-color: {{ call .Theme `color.BackgroundBlack` }};
	opacity: 0.5; 
}

#js-history-popup-header div {
	width: 50%;
	display: inline-block;
}

#js-history-popup-close {
	text-align: right;
}

#js-history-popup-close:hover {
	cursor: pointer;
}

#js-history-popup-clear {
	margin-right: 15px;
}

#js-history-popup-clear:hover {
	cursor: pointer;
}

#js-history-popup-list-div {
	font-family: monospace;

	margin-top: 15px;
	margin-bottom: 15px;

	overflow: auto;
	width: 100%;
	height: 100%;
}

#js-history-popup-list {
	margin: 0;
	margin-bottom: 65px;
}
`;
	let styleSheet = document.createElement("style")
	styleSheet.innerText = newStyleSheet
	document.head.appendChild(styleSheet)

	// Add button to header
	document.getElementsByClassName("header-right")[0].insertAdjacentHTML("afterbegin", "<h4 id='js-history-button' onclick='historyPopUpShow()'>{{ call .Translate `historyJS.History` }}</h4>");

	// Add history pop-up background
	document.body.insertAdjacentHTML("afterbegin", "<div style='visibility: hidden;' id='js-history-popup-background' onclick='historyPopUpHide()'></div>")	

	// If local storage is not supported
	if (isLocalStoageSupported() == false) {
	document.body.insertAdjacentHTML("afterbegin", `<div style='visibility: hidden;' id='js-history-popup'>
<div id='js-history-popup-header'>
	<div><h4 style='margin: 0;'>{{ call .Translate `historyJS.History` }}</h4></div
	><div id='js-history-popup-close' onclick='historyPopUpHide()'>&times;</div>
</div>
<hr/>
<p>{{ call .Translate `historyJS.LocalStorageNotSupported1` }}</p>
<p>{{ call .Translate `historyJS.LocalStorageNotSupported2` }}</p>
`);
		return;
	}

	// Add history pop-up
	document.body.insertAdjacentHTML("afterbegin", `<div style='visibility: hidden;' id='js-history-popup'>
<div id='js-history-popup-header'>
	<div><h4 style='margin: 0;'>{{ call .Translate `historyJS.History` }}</h4></div
	><div id='js-history-popup-close' onclick='historyPopUpHide()'>&times;</div>
</div>
<hr/>
<div>
	<label class='checkbox'><input id='js-history-popup-enable' onchange = 'historyEnable()' type='checkbox'></input>{{ call .Translate `historyJS.EnableHistory` }}</label
	><span id='js-history-popup-clear' class='text-red' onclick='historyClear()'>{{ call .Translate `historyJS.ClearHistory` }}</span>
</div>
<div id='js-history-popup-list-div'><ul id='js-history-popup-list'></ul></div>`);

	// Set "Remember history" checkbox state
	document.getElementById("js-history-popup-enable").checked = !localStorage.getItem("DisableHistory");

	// If exist "create paste" form path it
	let createPasteForm = document.getElementById("create-paste-form");
	if (createPasteForm != null) {
		createPasteForm.addEventListener("submit", (event) => {
			event.preventDefault();

			// Get form data
			let data = "";
			let title = "";
			
			Array.from(createPasteForm.elements)
				.filter((item) => !!item.name)
				.map((element) => {
					let { name, value, type } = element;

					if (type == "checkbox") {
						if (element.checked) {
							value = "true";
						} else {
							value = "false";
						}
					}

					if (name == "title") {
						title = value;
					}

					data = data + "&" + name + "=" + encodeURIComponent(value);
				})
			data = data.slice(1);

			// Send request
			var xhr = new XMLHttpRequest();
			xhr.responseType = "json";
			xhr.open("POST", "/api/v1/new", true);
			xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

			xhr.onload = () => {
				// Check HTTP code
				if (xhr.status != 200) {
					switch (xhr.status) {
						case 400: alert("{{ call .Translate `error.400` | call .Translate `historyJS.Error` 400 }}"); break;
						case 401: alert("{{ call .Translate `error.401` | call .Translate `historyJS.Error` 401 }}"); break;
						case 404: alert("{{ call .Translate `error.404` | call .Translate `historyJS.Error` 404 }}"); break;
						case 405: alert("{{ call .Translate `error.405` | call .Translate `historyJS.Error` 405 }}"); break;
						case 413: alert("{{ call .Translate `error.413` | call .Translate `historyJS.Error` 413 }}"); break;
						case 429: alert("{{ call .Translate `error.429` | call .Translate `historyJS.Error` 429 }}"); break;
						case 500: alert("{{ call .Translate `error.500` | call .Translate `historyJS.Error` 500 }}"); break;
						default: alert("{{ call .Translate `historyJS.ErrorUnknown` `"+xhr.status+"` }}"); break;	
					}
					return;
				}

				// Save to history
				if (localStorage.getItem("DisableHistory") != "true") {
					let historyJSON = localStorage.getItem("history");
					let history = [];
					if (historyJSON != null) {
						history = JSON.parse(historyJSON);
					}
					
					history.splice(0, 0, {id: xhr.response.id, createTime: xhr.response.createTime, deleteTime: xhr.response.deleteTime, title: title});
					localStorage.setItem("history", JSON.stringify(history));	
				}

				// Redirect
				window.location = window.location + xhr.response.id;
			};
			
			xhr.send(data);
			
			return false;
		});
	}
});
