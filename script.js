const terminal = document.getElementById('terminal');
const input = document.getElementById('input');

input.addEventListener('keydown', function(event) {
    if (event.key === 'Enter') {
        sendCommand(input.value);
        input.value = '';
    }
});

const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = function(event) {
    console.log("Connected to WebSocket");
};

ws.onmessage = function(event) {
    const safeText = escapeHTML(event.data);
    terminal.innerHTML += `<div>${safeText}</div>`;
};

ws.onerror = function(event) {
    console.error("WebSocket error observed:", event);
};

ws.onclose = function(event) {
    console.log("WebSocket connection closed:", event);
};

function sendCommand(command) {
    ws.send(command);
    const safeCommand = escapeHTML(command);
    terminal.innerHTML += `<div>> ${safeCommand}</div>`;
}

function escapeHTML(text) {
    return text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
}
