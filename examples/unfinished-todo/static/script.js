const loginForm = document.getElementById('login-form');
const registerForm = document.getElementById('register-form');
const loginMessage = document.getElementById('login-message');
const registerMessage = document.getElementById('register-message');

if (loginForm) {
    loginForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const username = loginForm.username.value;
        const password = loginForm.password.value;

        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams({ username, password }),
        });

        const data = await response.json();

        if (data.success) {
            loginMessage.textContent = data.message;
            window.location.href = '/';
        } else {
            loginMessage.textContent = data.message;
        }
    });
}

if (registerForm) {
    registerForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const username = registerForm.username.value;
        const password = registerForm.password.value;

        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams({ username, password }),
        });

        const data = await response.json();

        if (data.success) {
            registerMessage.textContent = data.message;
            window.location.href = '/login';
        } else {
            registerMessage.textContent = data.message;
        }
    });
}