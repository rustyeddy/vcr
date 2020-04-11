var ws = new WebSocket("ws://" + document.location.host + "/ws");
var doingAI = false;

// This is in VideoPlayer.vue
function record() {
    console.log("Record");
    var sendmsg = {
	t: "video",
	l: 1,
	v: "on"
    };
    ws.send(JSON.stringify(sendmsg));
}

// This is in VideoPlayer.vue
function pause() {
    console.log("Pause");
    var sendmsg = {
	t: "video",
	l: 1,
	v: "off"
    };
    ws.send(JSON.stringify(sendmsg));
}

function toggleAI() {
    console.log("doai");
    var cmd = "on"

    if (doingAI) {
	doingAI = false;
	cmd = "off";
    } else {
	doingAI = true;
	cmd = "on";
    }
    var sendmsg = {
	t: "ai",
	l: 0,
	v: cmd,
    };
    ws.send(JSON.stringify(sendmsg));
}

window.addEventListener("load", function(evt) {
    ws.binaryType = 'arraybuffer';

    ws.onopen = function(evt) {
        console.log("OPEN");
	var sendmsg = {
	    message: "hello",
	};
	ws.send(JSON.stringify(sendmsg));
    }
    
    ws.onclose = function(evt) {
        console.log("Websocket CLOSE"); 
        ws = null;
    }
    
    // We assume the incoming message is a JSON string containing a single
    // field 'message' with a string as a value.
    ws.onmessage = function(evt) {
	var obj = JSON.parse(evt.data);
	if (obj == null) {
	    console.log("WS bummer to message");
	    return;
	}

	console.log(obj);
	for (id in obj) {
	    var o = obj[id];
	    var key = o['k'];
	    var val = o['v'];

	    var ele = document.getElementById(key);
	    if (!ele) {
		console.log("Unknown element: " + id);
		continue;
	    }

	    ele.innerHTML = val;
	}
    }
    
    ws.onerror = function(evt) {
        console.log("WS ERROR: " + evt.data);
    }

    return false;
});
