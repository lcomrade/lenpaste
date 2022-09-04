// Lenpaste should work absolutely fine without this script.
// Although you will lose secondary functionality such as line numbers when you create a new paste.
// Therefore, if you are concerned about privacy, you can disable JavaScript in your browser.

function historyRefreshList() {
	// Get and clean list
	let listElement = document.getElementById("js-history-popup-list");
	listElement.innerHTML = "";

	// Read locale storage
	let historyJSON = localStorage.getItem("history");
	if (historyJSON != null) {
		let history = JSON.parse(historyJSON);

		for (let i = 0; history.length > i; i++) {
			listElement.insertAdjacentHTML("beforeend", "<li><a href='/"+history[i].id+"'>"+history[i].title+"</a></li>")
		}
	}
}

function historyPopUpShow() {
	document.getElementById("js-history-popup").style.visibility = "visible";
	document.getElementById("js-history-popup-background").style.visibility = "visible";
	document.addEventListener("keydown", historyPopUpEscEvent);
	historyRefreshList();
}

function historyPopUpHide() {
	document.getElementById("js-history-popup").style.visibility = "hidden";
	document.getElementById("js-history-popup-background").style.visibility = "hidden";
	document.removeEventListener("keydown", historyPopUpEscEvent);
}

function historyPopUpEscEvent(event) {
	// If ESC pressed
	if (event.keyCode == 27) {
		historyPopUpHide();
	}
}

function historyClear() {
	localStorage.removeItem("history");
	historyRefreshList();
}

document.addEventListener("DOMContentLoaded", () => {
	// Edit CSS
	let newStyleSheet = `
#js-history-button:hover {
	cursor: pointer;
}

#js-history-popup {
	background: #444444;
	padding: 20px;

	position: fixed;
	z-index: 2;
	top: 15%;
	bottom: 15%;
	right: 20%;
	left: 20%;
}

#js-history-popup-background {
	width: 100%;
	height: 100%;
	position: fixed;
	top: 0;
	left: 0;
	z-index: 1;
	background-color: black;
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

#js-history-popup-list-div {
	overflow-x: hidden;
	overflow-y: auto;
}
`;
	let styleSheet = document.createElement("style")
	styleSheet.innerText = newStyleSheet
	document.head.appendChild(styleSheet)

	// Add button to header
	document.getElementsByClassName("header-right")[0].insertAdjacentHTML("afterbegin", "<h4 id='js-history-button' onclick='historyPopUpShow()'>{{ call .Translate `historyJS.History` }}</h4>");

	// Add history pop-up
	document.body.insertAdjacentHTML("afterbegin", "<div style='visibility: hidden;' id='js-history-popup-background'></div>")	
	document.body.insertAdjacentHTML("afterbegin", `<div style='visibility: hidden;' id='js-history-popup'>
<div id='js-history-popup-header'>
	<div><h4 style='margin: 0;'>{{ call .Translate `historyJS.History` }}</h4></div
	><div id='js-history-popup-close' onclick='historyPopUpHide()'>&times;</div>
</div>
<hr/>
<div id='js-history-popup-list-div'><ul id='js-history-popup-list'></ul></div>
<div onclick='historyClear()'>{{ call .Translate `historyJS.ClearHistory` }}</div>
</div>`);

	// If exist "create paste" form path it
	let createPasteForm = document.getElementById("create-paste-form");
	if (createPasteForm != null) {
		createPasteForm.addEventListener("submit", (event) => {
			event.preventDefault();

			let historyJSON = localStorage.getItem("history");
			let history = [];
			if (historyJSON != null) {
				history = JSON.parse(historyJSON);
			}
			history.push({id: "404", title: "Test 124"});
			localStorage.setItem("history", JSON.stringify(history));
			
			return false;
		});
	}
});
