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
    terminal.innerHTML += `<div>${event.data}</div>`;
};

function sendCommand(command) {
    ws.send(command);
    terminal.innerHTML += `<div>> ${command}</div>`;
}

