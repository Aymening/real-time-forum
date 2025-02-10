const socket = new WebSocket("ws://localhost:8080/ws");
const messageInput = document.getElementById("message-input");
const messagesDiv = document.getElementById("messages");
let currentSender = "User1";
let currentReceiver = "User2";

// Handle WebSocket connection
socket.onopen = () => {
    console.log("Connected to WebSocket server");
};

socket.onmessage = (event) => {
    const message = JSON.parse(event.data);

    // Check if the message belongs to the current conversation
    if (
        (message.sender === currentSender && message.receiver === currentReceiver) ||
        (message.sender === currentReceiver && message.receiver === currentSender)
    ) {
        addMessage(message.sender, message.content);
    }
};

function sendMessage() {
    const message = {
        sender: currentSender,
        receiver: currentReceiver,
        content: messageInput.value,
    };

    socket.send(JSON.stringify(message));
    messageInput.value = "";
}

function fetchConversation(sender, receiver) {
    fetch(`/get-conversation?sender=${sender}&receiver=${receiver}`)
        .then((response) => response.json())
        .then((data) => {
            messagesDiv.innerHTML = "";
            data.forEach((msg) => {
                addMessage(msg.sender, msg.content);
            });
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        });
}

function selectConversation(sender, receiver) {
    currentSender = sender;
    currentReceiver = receiver;
    document.getElementById("current-conversation").textContent = `${sender} & ${receiver}`;
    fetchConversation(sender, receiver);
}

function addMessage(sender, content) {
    const msgElement = document.createElement("div");
    msgElement.className = `message ${sender === currentSender ? "self" : "other"}`;

    const contentElement = document.createElement("div");
    contentElement.className = "content";
    contentElement.textContent = content;

    msgElement.appendChild(contentElement);
    messagesDiv.appendChild(msgElement);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

// Load initial conversation
fetchConversation(currentSender, currentReceiver);