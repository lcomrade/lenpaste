// Lenpaste should work absolutely fine without this script.
// Although you will lose secondary functionality such as line numbers when you create a new paste.
// Therefore, if you are concerned about privacy, you can disable JavaScript in your browser.

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
