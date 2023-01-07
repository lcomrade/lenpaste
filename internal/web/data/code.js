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

function copyToClipboard(text) {
	let tmp   = document.createElement("textarea");
	let focus = document.activeElement;

	tmp.value = text;

	document.body.appendChild(tmp);
	tmp.select();
	document.execCommand("copy");
	document.body.removeChild(tmp);
	focus.focus();
}

function copyButton(element) {
	let result = "";

	let strings = element.parentNode.getElementsByTagName("code")[0].textContent.split("\n");
	let stringsLen = strings.length;
	let cutLen = stringsLen.toString().length;
	for (let i = 0; stringsLen > i; i++) {
		if (i != 0) {
			result = result + "\n"
		}

		result = result + strings[i].slice(cutLen);
	}

	result = result.trim() + "\n";
	copyToClipboard(result);
}


document.addEventListener("DOMContentLoaded", () => {
	// Edit CSS
	let newStyleSheet = `
pre {
	position: relative;
	overflow: auto;
}
	
pre button {
	visibility: hidden;
}

pre:hover > button {
	visibility: visible;
}
`;
	let styleSheet = document.createElement("style")
	styleSheet.innerText = newStyleSheet
	document.head.appendChild(styleSheet)

	// Edit pre tags
	let preElements = document.getElementsByTagName("pre");

	for (var i = 0; preElements.length > i; i++) {
		preElements[i].insertAdjacentHTML("beforeend", "<button class='button-green' style='position: absolute; top: 16px; right: 16px; margin: 0; animation: fadeout .2s both;' onclick='copyButton(this)'>{{call .Translate `codeJS.Paste`}}</button>");
	}
});
