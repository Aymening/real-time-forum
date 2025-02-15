const socket = new WebSocket("ws://localhost:8228/ws");
// console.log(socket.OPEN)
const messageInput = document.getElementById("message-input");
const messagesDiv = document.getElementById("messages");

// Handle WebSocket connection
socket.onopen = () => {
    console.log("Connected to WebSocket server");
};


socket.onmessage = (event) => {
    const message = JSON.parse(event.data);

    // console.log(message.sender)
    addMessage(message, JSON.parse(messageInput.dataset.user))
    // Check if the message belongs to the current conversation
    // if (
    //     (message.sender === currentSender && message.receiver === currentReceiver) ||
    //     (message.sender === currentReceiver && message.receiver === currentSender)
    // ) {
    //     // addMessage(message.sender, message.content);
    // }
};

function sendMessage() {

    const user = JSON.parse(messageInput.dataset.user)
    // console.log(user)

    const message = {

        receiver: user.id,
        // receiver: currentReceiver,
        content: messageInput.value,
    };
    console.log(messageInput.value);


    socket.send(JSON.stringify(message));
    messageInput.value = "";
}

// fetch users


function fetchUsers() {
    fetch('/api/contacts')
        .then(response => response.json())
        .then(data => {
            let usersArray = Object.values(data);

            if (!Array.isArray(usersArray)) {
                console.error("API did not return an array:", usersArray);
                return;
            }

            let conversationList = document.getElementById("conversation-list");
            conversationList.innerHTML = "";

            usersArray.forEach(user => {
                let li = document.createElement("li");
                li.textContent = user.name;
                li.classList.add("contact"); // Add class for styling & selection

                // Attach event listener to each contact
                li.addEventListener("click", function () {
                    selectConversation(user);
                });

                conversationList.appendChild(li);
            });
        })
        .catch(error => console.error("Error fetching users:", error));
}

window.onload = function () {
    fetchUsers();
    // document.getElementById("sendButton").classList.add("hidden")
    // document.getElementById("message-input").classList.add("hidden")
};

function fetchConversation(receiver) {
    fetch(`/api/chat?receiver=${receiver.id}`)
        .then((response) => response.json())
        .then((data) => {
            // console.log('hello');

            const msgInput = document.getElementById("message-input")
            let receiver = JSON.parse(msgInput.dataset.user)

            messagesDiv.innerHTML = "";
            data.forEach((msg) => {
                addMessage(msg, receiver);
            });
            messagesDiv.scrollTop = messagesDiv.scrollHeight;

        }).catch(console.log);

}

function selectConversation(receiver) {
    // console.log(receiver);

    // currentSender = sender;
    // currentReceiver = receiver;
    document.getElementById("current-conversation").textContent = `${receiver.name}`;
    document.getElementById("message-input").setAttribute('data-user', JSON.stringify(receiver))
    document.getElementById("sendButton").classList.remove("hidden");
    document.getElementById("message-input").classList.remove("hidden");
    // console.log(msg.sender);

    fetchConversation(receiver);
}

function addMessage(msg, currentReceiver) {

    const msgElement = document.createElement("div");
    // console.log(msg, currentReceiver)

    msgElement.className = `message ${msg.sender !== currentReceiver.id ? "self" : "other"}`;


    const contentElement = document.createElement("div");
    contentElement.className = "content";
    contentElement.textContent = msg.content;

    msgElement.appendChild(contentElement);
    messagesDiv.appendChild(msgElement);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
    // console.log('this is the sender',msg.sender);
    // console.log('this is the recever',currentReceiver);


}




let currentReceiver = document.getElementById("message-input").dataset.user

// Load initial conversation
fetchConversation(JSON.parse(currentReceiver));


