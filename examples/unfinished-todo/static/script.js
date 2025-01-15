const loginForm = document.getElementById('login-form');
const registerForm = document.getElementById('register-form');
const loginMessage = document.getElementById('login-message');
const registerMessage = document.getElementById('register-message');

if (loginForm) {
    loginForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const username = loginForm.username.value;
        const password = loginForm.password.value;

        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: new URLSearchParams({ username, password }),
            });

            if (response.ok) {
                const data = await response.json();

                loginMessage.textContent = data.message;

                if (data.success) {
                    setTimeout(() => {
                        window.location.href = '/';
                    }, 500);
                }
            } else {
                const errorData = await response.json();
                loginMessage.textContent = errorData.message || 'An error occurred during login.';
            }
        } catch (error) {
            console.error('Fetch error:', error);
            loginMessage.textContent = 'A network error occurred.';
        }
    });
}

if (registerForm) {
    registerForm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const username = registerForm.username.value;
        const password = registerForm.password.value;

        try {
            const response = await fetch('/api/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: new URLSearchParams({ username, password }),
            });

            if (response.ok) {
                const data = await response.json();
                registerMessage.textContent = data.message;

                if (data.success) {
                    setTimeout(() => {
                        window.location.href = '/login';
                    }, 500);
                }
            } else {
                const errorData = await response.json();
                registerMessage.textContent = errorData.message || 'An error occurred during registration.';
            }
        } catch (error) {
            console.error('Fetch error:', error);
            registerMessage.textContent = 'A network error occurred.';
        }
    });
}