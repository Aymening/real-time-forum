const form = document.getElementById('registerForm');

form.addEventListener('submit', async (event) => {
    event.preventDefault();

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const name = document.getElementById('name').value;

    document.getElementById('emailError').innerText = '';
    document.getElementById('passwordError').innerText = '';
    document.getElementById('nameError').innerText = '';

    let isValid = true;

    const nameRegex = /^[A-Za-z0-9._-]{2,15}$/;
    if (!nameRegex.test(name)) {
        document.getElementById('nameError').innerText = 'Invalid name';
        isValid = false;
    }

    const emailRegex = /^(\w+@\w+\.[a-zA-Z]+)$/;
    if (!emailRegex.test(email)) {
        document.getElementById('emailError').innerText = 'Invalid email';
        isValid = false;
    }

    const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.+[@$!%*?&])[a-zA-Z\d@$!%*#?&]{8,}$/;
    if (!passwordRegex.test(password)) {
        document.getElementById('passwordError').innerText = 'your password should include at least one uppercase letter, one lowercase letter, a number, and a special character 😅';
        isValid = false;
    }

    if (isValid) {
        fetch("/register", {
            method: 'POST',
            headers: { 'Content-Type': "application/json" },
            body: JSON.stringify({ 'name': name, 'email': email, 'password': password })
        })
            .then(response => {
                console.log(response.status);

                if (response.ok) {
                    response.json();
                    console.log("1")
                    window.location.href = "/login";
                } else if (response.status === 409) {
                    console.log("2")
                    return response.json().then(data => {
                        throw new Error(data.message);
                    });
                } else {
                    return response.json().then(data => {
                        console.log("3")
                        alert(data.message);
                    });
                }
            })

            // .then(data => {
            //     console.log('Inscription valid', data);
            //     window.location.href = "/login";
            // })

            .catch(error => {
                console.error('Failed to fetch page: ', error);
                document.getElementById('emailError').innerText = error.message;
            })
    }
});
