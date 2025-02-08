// Connect WebSocket
let socket = new WebSocket("ws://localhost:8080/ws");

// User login
let username = prompt("Enter your username:");
socket.onopen = function () {
    console.log("✅ Connected to WebSocket");
    socket.send(JSON.stringify(username)); // Send username when connecting
};

socket.onmessage = function (event) {
    let msg = JSON.parse(event.data);
    console.log("📩 Received Message:", msg);

    // Show message in chat
    displayMessage(msg.sender, msg.message);
};

// Send message
function sendMessage() {
    let receiver = document.getElementById("receiver").value.trim();
    let message = document.getElementById("message").value.trim();

    if (receiver === "" || message === "") {
        alert("⚠️ Please fill in all fields!");
        return;
    }

    let msg = { sender: username, receiver: receiver, message: message };
    socket.send(JSON.stringify(msg));

    // Show message instantly
    displayMessage(username, message);

    document.getElementById("message").value = ""; // Clear input
}

// Display messages in the chat UI
function displayMessage(user, message) {
    let chatBox = document.getElementById("chat-box");
    let newMessage = document.createElement("div");
    newMessage.classList.add("message");
    newMessage.innerHTML = `<strong>${user}:</strong> ${message}`;
    chatBox.appendChild(newMessage);

    // Auto-scroll
    chatBox.scrollTop = chatBox.scrollHeight;
}
