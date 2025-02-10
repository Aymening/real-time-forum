const socket = new WebSocket("ws://localhost:8080/ws");
const messageInput = document.getElementById("message-input");
const messagesDiv = document.getElementById("messages");

socket.onopen = () => {
    console.log("Connected to WebSocket server");
};

socket.onmessage = (event) => {
    const message = JSON.parse(event.data);
    const msgElement = document.createElement("div");
    msgElement.textContent = `${message.sender}: ${message.content}`;
    messagesDiv.appendChild(msgElement);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
};

function sendMessage() {
    const message = {
        sender: "User1", // Replace with dynamic user identification
        receiver: "User2", // Replace with dynamic user identification
        content: messageInput.value,
    };

    socket.send(JSON.stringify(message));
    messageInput.value = "";
}

function fetchConversation() {
    fetch("/get-conversation?sender=User1&receiver=User2")
        .then((response) => response.json())
        .then((data) => {
            messagesDiv.innerHTML = "";
            data.forEach((msg) => {
                const msgElement = document.createElement("div");
                msgElement.textContent = `${msg.sender}: ${msg.content}`;
                messagesDiv.appendChild(msgElement);
            });
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        });
}

fetchConversation(); // Load conversation on page load