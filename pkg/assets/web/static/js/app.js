// Main application functionality

document.addEventListener('DOMContentLoaded', function() {
    // Check authentication status and show appropriate page
    if (isAuthenticated()) {
        showAppPage();
    } else {
        showLoginPage();
    }
    
    // Login form submission
    const loginForm = document.getElementById('login-form');
    loginForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        const loginError = document.getElementById('login-error');
        
        if (!username || !password) {
            loginError.textContent = 'Username and password are required';
            loginError.classList.remove('hidden');
            return;
        }
        
        try {
            loginError.classList.add('hidden');
            await login(username, password);
            showAppPage();
        } catch (error) {
            loginError.textContent = error.message || 'Login failed. Please check your credentials.';
            loginError.classList.remove('hidden');
        }
    });
    
    // Logout link
    const logoutLink = document.getElementById('logout-link');
    logoutLink.addEventListener('click', function(e) {
        e.preventDefault();
        logout();
    });
    
    // User menu toggle
    const userMenuButton = document.getElementById('user-menu-button');
    const userMenu = document.getElementById('user-menu');
    
    userMenuButton.addEventListener('click', function() {
        userMenu.classList.toggle('hidden');
    });
    
    // Close user menu when clicking outside
    document.addEventListener('click', function(e) {
        if (!userMenuButton.contains(e.target) && !userMenu.contains(e.target)) {
            userMenu.classList.add('hidden');
        }
    });
    
    // Change password modal
    const changePasswordLink = document.getElementById('change-password-link');
    const changePasswordModal = document.getElementById('change-password-modal');
    const closeModalButton = document.getElementById('close-modal-button');
    const cancelChangePassword = document.getElementById('cancel-change-password');
    
    changePasswordLink.addEventListener('click', function(e) {
        e.preventDefault();
        userMenu.classList.add('hidden');
        changePasswordModal.classList.remove('hidden');
    });
    
    function closeModal() {
        changePasswordModal.classList.add('hidden');
        document.getElementById('change-password-error').classList.add('hidden');
        document.getElementById('change-password-success').classList.add('hidden');
        document.getElementById('change-password-form').reset();
    }
    
    closeModalButton.addEventListener('click', closeModal);
    cancelChangePassword.addEventListener('click', closeModal);
    
    // Change password form submission
    const changePasswordForm = document.getElementById('change-password-form');
    changePasswordForm.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const currentPassword = document.getElementById('current-password').value;
        const newPassword = document.getElementById('new-password').value;
        const confirmPassword = document.getElementById('confirm-password').value;
        const errorElement = document.getElementById('change-password-error');
        const successElement = document.getElementById('change-password-success');
        
        // Reset messages
        errorElement.classList.add('hidden');
        successElement.classList.add('hidden');
        
        // Validate inputs
        if (!currentPassword || !newPassword || !confirmPassword) {
            errorElement.textContent = 'All fields are required';
            errorElement.classList.remove('hidden');
            return;
        }
        
        if (newPassword !== confirmPassword) {
            errorElement.textContent = 'New password and confirmation do not match';
            errorElement.classList.remove('hidden');
            return;
        }
        
        try {
            await changePassword(currentPassword, newPassword);
            successElement.textContent = 'Password changed successfully';
            successElement.classList.remove('hidden');
            
            // Reset form
            changePasswordForm.reset();
            
            // Close modal after 2 seconds
            setTimeout(closeModal, 2000);
        } catch (error) {
            errorElement.textContent = error.message || 'Failed to change password';
            errorElement.classList.remove('hidden');
        }
    });
    
    // Execute command
    const executeButton = document.getElementById('execute-button');
    const commandInput = document.getElementById('command');
    const commandTypeSelect = document.getElementById('command-type');
    const responseElement = document.getElementById('response');
    
    executeButton.addEventListener('click', async function() {
        const command = commandInput.value.trim();
        const commandType = commandTypeSelect.value;
        
        if (!command) {
            responseElement.textContent = 'Please enter a command';
            return;
        }
        
        responseElement.textContent = 'Processing...';
        
        try {
            const token = getAuthToken();
            
            if (!token) {
                throw new Error('Not authenticated');
            }
            
            const payload = {
                command: command
            };
            
            if (commandType) {
                payload.type = commandType;
            }
            
            const response = await fetch('/api/v1/execute', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(payload)
            });
            
            if (!response.ok) {
                if (response.status === 401) {
                    // Token might be expired, try to refresh
                    await refreshToken();
                    // Retry the request
                    return executeButton.click();
                }
                
                const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
                throw new Error(errorData.error || response.statusText);
            }
            
            const data = await response.json();
            
            if (data.success) {
                responseElement.textContent = data.output;
            } else {
                responseElement.textContent = `Error: ${data.error || 'Unknown error'}`;
            }
        } catch (error) {
            responseElement.textContent = `Error: ${error.message}`;
            console.error('Execute error:', error);
            
            // If authentication error, redirect to login
            if (error.message === 'Not authenticated') {
                logout();
            }
        }
    });
    
    // Allow pressing Enter in command input to execute
    commandInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            executeButton.click();
        }
    });
});
