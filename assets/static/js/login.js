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
            console.log(response.status);

            if (response.ok) {
                response.json()
                window.location.href = "/";
            } else if (response.status === 409) {
                return response.json().then(data => {
                    throw new Error(data.message);
                });
            } else {
                return response.json().then(data => {
                    alert(data.message);
                });
            }
        })
        // .then(data => {
        //     console.log('login valid', data);
        //     window.location.href = "/";
        // })

        .catch(error => {
           // console.error('Failed to fetch page: ', error)
            document.getElementById('emailError').innerText = error.message;
        })
})