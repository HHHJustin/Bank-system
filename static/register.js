document.addEventListener('DOMContentLoaded', () => {
    // 处理注册表单提交
    const registerForm = document.getElementById('registerForm');
    registerForm.addEventListener('submit', async function(event) {
        event.preventDefault(); // 阻止表单的默认提交行为

        // 获取表单数据
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        const fullname = document.getElementById('fullname').value;
        const email = document.getElementById('email').value;

        // 构建JSON对象
        const registerData = {
            username: username,
            password: password,
            fullname: fullname,
            email: email
        };

        try {
            // 发送POST请求
            const response = await fetch('/users/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(registerData)
            });

            // 处理响应
            if (response.ok) {
                // 显示成功消息
                document.getElementById('success-message').textContent = 'Registration successful! Redirecting to login page...';
                document.getElementById('success-message').style.display = 'block';
                document.getElementById('error-message').style.display = 'none';

                // 3秒后跳转到登录页面
                setTimeout(() => {
                    window.location.href = '/users/login';
                }, 3000);
            } else {
                const errorData = await response.json();
                const errorMessage = errorData.error || 'Registration failed';
                displayErrorMessage(errorMessage);
            }
        } catch (error) {
            console.error('Error:', error);
            displayErrorMessage('An unexpected error occurred');
        }
    });

    document.getElementById('returnButton').addEventListener('click', function() {
        window.location.href = '/users/login';
    });

    function displayErrorMessage(message) {
        const errorMessageDiv = document.getElementById('error-message');
        errorMessageDiv.textContent = message;
        errorMessageDiv.style.display = 'block';
        document.getElementById('success-message').style.display = 'none';
    }
});
