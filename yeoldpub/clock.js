function startTime() {
    var today = new Date();
    var h = today.getHours();
    var m = today.getMinutes();
    var s = today.getSeconds();
    m = checkTime(m);
    s = checkTime(s);

    document.getElementById('clock').innerHTML =
	"<span id='clock' class='text-center'>" + h + ":" + m + "</span>" +
	"<span id='seconds' class='small'>  " + s + "  </span>";
    
    var t = setTimeout(startTime, 500);
}

function checkTime(i) {
    if (i < 10) {i = "0" + i};  // add zero in front of numbers < 10
    return i;
}
