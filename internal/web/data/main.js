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
	var editor = document.getElementById("editor");

	editor.addEventListener("keydown", (e) => {
		// If TAB pressed
		if (e.keyCode === 9) {
			e.preventDefault();

			let startOrig = editor.selectionStart;
			let endOrig = editor.selectionEnd;

			editor.value = editor.value.substring(0, startOrig) + "\t" + editor.value.substring(endOrig);

			editor.selectionStart = editor.selectionEnd = startOrig + 1;
		}
	});

	// Add HTML and CSS code for line numbers support
	editor.insertAdjacentHTML("beforebegin", "<textarea id='editorLines' wrap='off' tabindex=-1 readonly>1</textarea>");
	var editorLines = document.getElementById("editorLines");
	editorLines.rows = editor.rows;
	
	var styleSheet = document.createElement("style");
	styleSheet.innerText = `
	#editor {
		margin-left: 60px;
		resize: none;
		
		width: calc(100% - 60px);
		min-width: calc(100% - 60px);
		max-width: calc(100% - 60px);
	}

	#editorLines {
		display: flex;
		user-select: none;

		text-align: right;
		position: absolute;
		resize: none;
		overflow-y: hidden;
		width: 60px;
		max-width: 60px;
		min-width: 60px;
	}

	#editor:focus-visible, #editorLines:focus-visible {
		outline: 0;
	}
`;
	document.head.appendChild(styleSheet);

	editorLines.addEventListener("focus", () => {
		editor.focus();
	});

	// Add JS code for line numbers
	editor.addEventListener("scroll", () => {
		editorLines.scrollTop = editor.scrollTop;
		editorLines.scrollLeft = editor.scrollLeft;
	});

	var lineCountCache = 0;
	editor.addEventListener("input", () => {
		let lineCount = editor.value.split("\n").length;

		if (lineCountCache != lineCount) {
			editorLines.value = "";
			
			for (var i = 0; i < lineCount; i++) {
				editorLines.value = editorLines.value + (i + 1) + "\n";
			}
			
			lineCountCache = lineCount;
		}
	});

	// Add symbol counter
	document.getElementById("symbolCounterContainer").innerHTML = "<span id='symbolCounter' class='text-grey'></span>";
	var symbolCounter = document.getElementById("symbolCounter");

	function updateSymbolCounter() {
		symbolCounter.textContent = editor.value.length;
		
		if (editor.maxLength !== -1) {
			symbolCounter.textContent = symbolCounter.textContent + "/" + editor.maxLength;
		} else {
			symbolCounter.textContent = symbolCounter.textContent + "/&infin;";
		}
	}

	editor.addEventListener("input", updateSymbolCounter);
	updateSymbolCounter();
});
