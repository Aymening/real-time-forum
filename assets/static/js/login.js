import { showAlert } from "/static/js/alert.js";

const form = document.getElementById('loginForm');
form.addEventListener('submit', (event) => {
    event.preventDefault();
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    document.getElementById('emailError').innerText = '';
    document.getElementById('passwordError').innerText = '';

    fetch("/login", {
        method: 'POST',
        headers: { 'Content-Type': "application/json" },
        body: JSON.stringify({ 'email': email, 'password': password })
    })
        .then(response => {
            if (response.ok) {
                response.json()
                window.location.href = "/";
            } else if (response.status === 409) {
                return response.json().then(data => {
                    throw new Error(data.message);
                });
            } else {
                
                return response.json().then(data => {
                    showAlert(data.message);
                });
            }
        })
     
        .catch(error => {
            document.getElementById('Error').innerText = error.message;
        })
})