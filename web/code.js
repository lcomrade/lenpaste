// Lenpaste should work absolutely fine without this script.
// Although you will lose secondary functionality such as line numbers when you create a new paste.
// Therefore, if you are concerned about privacy, you can disable JavaScript in your browser.

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
	const pasteTxt = "{{call .Translate `codeJS.Paste`}}";

	let preElements = document.getElementsByTagName("pre");
	for (var i = 0; preElements.length > i; i++) {
		preElements[i].insertAdjacentHTML("beforeend", "<button class='button-green' style='position: absolute; top: 16px; right: 16px; margin: 0; visibility: hidden;' onclick='copyButton(this)'>"+pasteTxt+"</button>");
		preElements[i].style.position = "relative";
		preElements[i].style.overflow = "auto";

		preElements[i].addEventListener("mouseover", (event) => {
			event.currentTarget.getElementsByTagName("button")[0].style.visibility = "visible";
		});

		preElements[i].addEventListener("mouseout", (event) => {
			event.currentTarget.getElementsByTagName("button")[0].style.visibility = "hidden";
		});
	}
});
