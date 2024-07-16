document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('loginForm').addEventListener('submit', async function(event) {
        event.preventDefault(); // 阻止表单的默认提交行为

        // 获取表单数据
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        // 构建JSON对象
        const loginData = {
            username: username,
            password: password
        };

        try {
            // 发送POST请求
            const response = await fetch('/users/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(loginData)
            });

            // 处理响应
            if (response.ok) {
                const data = await response.json();
                // 将token存储到localStorage
                localStorage.setItem('accessToken', data.access_token);
                localStorage.setItem('refreshToken', data.refresh_token);

                // 使用fetchWithAuth请求/auth/account
                const accountResponse = await fetchWithAuth('/account');
                if (accountResponse.ok) {
                    const html = await accountResponse.text();
                    document.documentElement.innerHTML = html;
                } else {
                    throw new Error('Failed to fetch account page');
                }
            } else {
                const errorData = await response.json();
                const errorMessage = errorData.error || 'Login failed';
                displayErrorMessage(errorMessage);
            }
        } catch (error) {
            console.error('Error:', error);
            displayErrorMessage('An unexpected error occurred');
        }
    });

    document.getElementById('registerButton').addEventListener('click', function() {
        window.location.href = '/users/register';
    });

    function displayErrorMessage(message) {
        const errorMessageDiv = document.getElementById('error-message');
        errorMessageDiv.textContent = message;
        errorMessageDiv.style.display = 'block';
    }

    async function fetchWithAuth(url, options = {}) {
        const accessToken = localStorage.getItem('accessToken');
        const headers = {
            'Authorization': `Bearer ${accessToken}`,
            ...options.headers
        };
        const response = await fetch(url, {
            ...options,
            headers
        });
        return response;
    }
});
