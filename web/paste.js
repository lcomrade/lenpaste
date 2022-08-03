// Lenpaste should work absolutely fine without this script.
// Although you will lose secondary functionality such as line numbers when you create a new paste.
// Therefore, if you are concerned about privacy, you can disable JavaScript in your browser.

document.addEventListener("DOMContentLoaded", () => {
	function dateToString(date) {
		let dateStr = "";

		switch (date.getDay()) {
			case 0: dateStr = "Sun"; break;
			case 1: dateStr = "Mon"; break;
			case 2: dateStr = "Tue"; break;
			case 3: dateStr = "Wed"; break;
			case 4: dateStr = "Thu"; break;
			case 5: dateStr = "Fri"; break;
			case 6: dateStr = "Sat"; break;
		}

		dateStr = dateStr + ", " + date.getDate() + " ";

		switch (date.getMonth()) {
			case 0: dateStr = dateStr + "Jan"; break;
			case 1: dateStr = dateStr + "Feb"; break;
			case 2: dateStr = dateStr + "Mar"; break;
			case 3: dateStr = dateStr + "Apr"; break;
			case 4: dateStr = dateStr + "May"; break;
			case 5: dateStr = dateStr + "Jun"; break;
			case 6: dateStr = dateStr + "Jul"; break;
			case 7: dateStr = dateStr + "Aug"; break;
			case 8: dateStr = dateStr + "Sep"; break;
			case 9: dateStr = dateStr + "Oct"; break;
			case 10: dateStr = dateStr + "Nov"; break;
			case 11: dateStr = dateStr + "Dec"; break;
		}

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
