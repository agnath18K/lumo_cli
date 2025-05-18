// Lumo Connect web interface functionality

// Constants
const DefaultChunkSize = 5 * 1024 * 1024; // 5MB

document.addEventListener('DOMContentLoaded', function() {
    // Show connect page by default since authentication is handled by redirect script
    document.getElementById('connect-page').classList.remove('hidden');

    // Update username display
    const username = getUsername();
    document.getElementById('username-display').textContent = username;

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

    // Logout link
    const logoutLink = document.getElementById('logout-link');
    logoutLink.addEventListener('click', function(e) {
        e.preventDefault();
        // Clear authentication data
        localStorage.removeItem('lumo_token');
        localStorage.removeItem('lumo_refresh_token');
        localStorage.removeItem('lumo_username');
        localStorage.removeItem('lumo_token_expiry');
        // Redirect to main page
        window.location.href = '../';
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

    closeModalButton.addEventListener('click', function() {
        changePasswordModal.classList.add('hidden');
    });

    cancelChangePassword.addEventListener('click', function() {
        changePasswordModal.classList.add('hidden');
    });

    // Change password form submission
    const changePasswordForm = document.getElementById('change-password-form');
    changePasswordForm.addEventListener('submit', async function(e) {
        e.preventDefault();

        const currentPassword = document.getElementById('current-password').value;
        const newPassword = document.getElementById('new-password').value;
        const confirmPassword = document.getElementById('confirm-password').value;
        const changePasswordError = document.getElementById('change-password-error');

        if (!currentPassword || !newPassword || !confirmPassword) {
            changePasswordError.textContent = 'All fields are required';
            changePasswordError.classList.remove('hidden');
            return;
        }

        if (newPassword !== confirmPassword) {
            changePasswordError.textContent = 'New passwords do not match';
            changePasswordError.classList.remove('hidden');
            return;
        }

        try {
            changePasswordError.classList.add('hidden');

            const token = getAuthToken();

            if (!token) {
                throw new Error('Not authenticated');
            }

            const response = await fetch('/api/v1/auth/change-password', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    current_password: currentPassword,
                    new_password: newPassword
                })
            });

            if (!response.ok) {
                if (response.status === 401) {
                    // Token might be expired, try to refresh
                    await refreshToken();
                    // Retry the request
                    return changePasswordForm.dispatchEvent(new Event('submit'));
                }

                const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
                throw new Error(errorData.error || response.statusText);
            }

            // Reset form and close modal
            changePasswordForm.reset();
            changePasswordModal.classList.add('hidden');

            // Show success message
            alert('Password changed successfully');
        } catch (error) {
            changePasswordError.textContent = error.message;
            changePasswordError.classList.remove('hidden');
            console.error('Change password error:', error);

            // If authentication error, redirect to main page
            if (error.message === 'Not authenticated') {
                handleAuthError();
            }
        }
    });

    // Connect functionality
    let activeConnection = null;
    let fileToSend = null;
    let fileTransferHistory = [];

    // Start server button
    const startServerButton = document.getElementById('start-server-button');
    startServerButton.addEventListener('click', async function() {
        const portInput = document.getElementById('server-port');
        const pathInput = document.getElementById('server-path');

        // Parse port as integer
        const port = parseInt(portInput.value, 10) || 8080;
        const path = pathInput.value || '';

        try {
            // Get token if available, but don't require it
            const token = getAuthToken() || '';

            const response = await fetch('/api/v1/connect/start-server', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    port: port,
                    path: path
                })
            });

            if (!response.ok) {
                if (response.status === 401) {
                    // Token might be expired, try to refresh
                    await refreshToken();
                    // Retry the request
                    return startServerButton.click();
                }

                const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
                throw new Error(errorData.error || response.statusText);
            }

            const data = await response.json();

            if (data && data.success) {
                // Update active connection
                activeConnection = {
                    mode: 'server',
                    ip: data.ip || 'localhost',
                    port: data.port || 8080,
                    status: 'active'
                };

                // Show active connection container
                updateConnectionInfo();
                document.getElementById('active-connection-container').classList.remove('hidden');
                document.getElementById('file-history-container').classList.remove('hidden');

                // Start WebSocket connection for real-time updates
                startWebSocketConnection();
            } else {
                throw new Error((data && data.error) || 'Failed to start server');
            }
        } catch (error) {
            alert(`Error: ${error.message}`);
            console.error('Start server error:', error);

            // No need to redirect for authentication errors since we don't require authentication
        }
    });

    // Discover button
    const discoverButton = document.getElementById('discover-button');
    discoverButton.addEventListener('click', async function() {
        try {
            // Get token if available, but don't require it
            const token = getAuthToken() || '';

            // Show loading state
            discoverButton.textContent = 'Discovering...';
            discoverButton.disabled = true;

            const response = await fetch('/api/v1/connect/discover', {
                method: 'GET',
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                if (response.status === 401) {
                    // Token might be expired, try to refresh
                    await refreshToken();
                    // Retry the request
                    discoverButton.textContent = 'Discover Devices';
                    discoverButton.disabled = false;
                    return discoverButton.click();
                }

                const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
                throw new Error(errorData.error || response.statusText);
            }

            const data = await response.json();

            // Reset button state
            discoverButton.textContent = 'Discover Devices';
            discoverButton.disabled = false;

            if (data && data.success) {
                // Make sure devices is an array
                const devices = Array.isArray(data.devices) ? data.devices : [];
                // Display discovered devices
                displayDiscoveredDevices(devices);
            } else {
                // Handle case where devices array is missing or success is false
                displayDiscoveredDevices([]);
                console.log("No devices found or discovery failed");
            }
        } catch (error) {
            // Reset button state
            discoverButton.textContent = 'Discover Devices';
            discoverButton.disabled = false;

            alert(`Error: ${error.message}`);
            console.error('Discover error:', error);

            // No need to redirect for authentication errors since we don't require authentication
        }
    });

    // Connect button
    const connectButton = document.getElementById('connect-button');
    connectButton.addEventListener('click', async function() {
        const ip = document.getElementById('client-ip').value;
        const port = document.getElementById('client-port').value;

        if (!ip) {
            alert('Please enter an IP address');
            return;
        }

        try {
            // Get token if available, but don't require it
            const token = getAuthToken() || '';

            // Show loading state
            connectButton.textContent = 'Connecting...';
            connectButton.disabled = true;

            const response = await fetch('/api/v1/connect/connect-to-peer', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    ip: ip,
                    port: port
                })
            });

            // Reset button state
            connectButton.textContent = 'Connect';
            connectButton.disabled = false;

            if (!response.ok) {
                if (response.status === 401) {
                    // Token might be expired, try to refresh
                    await refreshToken();
                    // Retry the request
                    return connectButton.click();
                }

                const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
                throw new Error(errorData.error || response.statusText);
            }

            const data = await response.json();

            if (data && data.success) {
                // Update active connection
                activeConnection = {
                    mode: 'client',
                    ip: ip,
                    port: port,
                    status: 'active'
                };

                // Show active connection container
                updateConnectionInfo();
                document.getElementById('active-connection-container').classList.remove('hidden');
                document.getElementById('file-history-container').classList.remove('hidden');

                // Start WebSocket connection for real-time updates
                startWebSocketConnection();
            } else {
                throw new Error((data && data.error) || 'Failed to connect to peer');
            }
        } catch (error) {
            alert(`Error: ${error.message}`);
            console.error('Connect error:', error);

            // No need to redirect for authentication errors since we don't require authentication
        }
    });

    // File upload handling
    const fileUploadInput = document.getElementById('file-upload');
    const fileUploadArea = fileUploadInput.parentElement.parentElement.parentElement;
    const sendFileButton = document.getElementById('send-file-button');

    fileUploadInput.addEventListener('change', function(e) {
        if (e.target.files.length > 0) {
            fileToSend = e.target.files[0];
            sendFileButton.disabled = false;
            fileUploadArea.querySelector('p').textContent = `Selected: ${fileToSend.name} (${formatFileSize(fileToSend.size)})`;
        } else {
            fileToSend = null;
            sendFileButton.disabled = true;
            fileUploadArea.querySelector('p').textContent = 'or drag and drop';
        }
    });

    // Drag and drop handling
    fileUploadArea.addEventListener('dragover', function(e) {
        e.preventDefault();
        fileUploadArea.classList.add('drag-over');
    });

    fileUploadArea.addEventListener('dragleave', function() {
        fileUploadArea.classList.remove('drag-over');
    });

    fileUploadArea.addEventListener('drop', function(e) {
        e.preventDefault();
        fileUploadArea.classList.remove('drag-over');

        if (e.dataTransfer.files.length > 0) {
            fileToSend = e.dataTransfer.files[0];
            sendFileButton.disabled = false;
            fileUploadArea.querySelector('p').textContent = `Selected: ${fileToSend.name} (${formatFileSize(fileToSend.size)})`;
        }
    });

    // Send file button
    sendFileButton.addEventListener('click', async function() {
        if (!fileToSend) {
            alert('Please select a file to send');
            return;
        }

        if (!activeConnection) {
            alert('No active connection');
            return;
        }

        try {
            // Get token if available, but don't require it
            const token = getAuthToken() || '';

            // Show loading state
            sendFileButton.textContent = 'Sending...';
            sendFileButton.disabled = true;

            let response;
            let data;

            // Check file size to determine transfer method
            if (fileToSend.size > 10 * 1024 * 1024) { // 10MB threshold
                // Use chunked transfer for large files
                await sendLargeFileWithChunkedTransfer(fileToSend);

                // For chunked transfer, we don't need to process response here
                // as it's handled in the sendLargeFileWithChunkedTransfer function

                // Reset button state
                sendFileButton.textContent = 'Send File';
                sendFileButton.disabled = false;

                // Reset file upload
                fileToSend = null;
                fileUploadInput.value = '';
                fileUploadArea.querySelector('p').textContent = 'or drag and drop';
                sendFileButton.disabled = true;

                return;
            } else {
                // Use regular transfer for small files
                // Create FormData
                const formData = new FormData();
                formData.append('file', fileToSend);

                response = await fetch('/api/v1/connect/send-file', {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`
                    },
                    body: formData
                });
            }

            // Reset button state
            sendFileButton.textContent = 'Send File';
            sendFileButton.disabled = false;

            if (!response || !response.ok) {
                if (response && response.status === 401) {
                    // Token might be expired, try to refresh
                    await refreshToken();
                    // Retry the request
                    return sendFileButton.click();
                }

                const errorData = await (response ? response.json().catch(() => ({ error: 'Unknown error' })) : { error: 'Unknown error' });
                throw new Error(errorData.error || (response ? response.statusText : 'Failed to send file'));
            }

            data = await response.json();

            if (data && data.success) {
                // Add to file transfer history
                addFileTransferHistoryEntry({
                    filename: fileToSend.name,
                    size: fileToSend.size,
                    direction: 'sent',
                    status: 'success',
                    time: new Date()
                });

                // Reset file upload
                fileToSend = null;
                fileUploadInput.value = '';
                fileUploadArea.querySelector('p').textContent = 'or drag and drop';
                sendFileButton.disabled = true;
            } else {
                throw new Error((data && data.error) || 'Failed to send file');
            }
        } catch (error) {
            alert(`Error: ${error.message}`);
            console.error('Send file error:', error);

            // No need to redirect for authentication errors since we don't require authentication
        }
    });

    // Disconnect button
    const disconnectButton = document.getElementById('disconnect-button');
    disconnectButton.addEventListener('click', async function() {
        if (!activeConnection) {
            return;
        }

        try {
            // Get token if available, but don't require it
            const token = getAuthToken() || '';

            const response = await fetch('/api/v1/connect/disconnect', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) {
                if (response.status === 401) {
                    // Token might be expired, try to refresh
                    await refreshToken();
                    // Retry the request
                    return disconnectButton.click();
                }

                const errorData = await response.json().catch(() => ({ error: 'Unknown error' }));
                throw new Error(errorData.error || response.statusText);
            }

            // Reset connection state
            activeConnection = null;
            document.getElementById('active-connection-container').classList.add('hidden');

            // Close WebSocket connection
            if (window.connectWebSocket) {
                window.connectWebSocket.close();
                window.connectWebSocket = null;
            }
        } catch (error) {
            alert(`Error: ${error.message}`);
            console.error('Disconnect error:', error);

            // No need to redirect for authentication errors since we don't require authentication
        }
    });

    // Helper functions

    function displayDiscoveredDevices(devices) {
        const container = document.getElementById('discovered-devices-container');
        const devicesContainer = document.getElementById('discovered-devices');

        // Clear existing devices
        devicesContainer.innerHTML = '';

        if (devices.length === 0) {
            container.classList.add('hidden');
            alert('No devices found');
            return;
        }

        // Add device cards
        devices.forEach(device => {
            const card = document.createElement('div');
            card.className = 'device-card';
            card.innerHTML = `
                <h3 class="text-md font-medium text-gray-900">${device.name}</h3>
                <p class="text-sm text-gray-600">IP: ${device.ip}</p>
                <p class="text-sm text-gray-600">Port: ${device.port}</p>
                ${device.username ? `<p class="text-sm text-gray-600">User: ${device.username}</p>` : ''}
                <button class="mt-2 w-full flex justify-center py-1 px-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                    Connect
                </button>
            `;

            // Add click handler for the connect button
            const connectBtn = card.querySelector('button');
            connectBtn.addEventListener('click', function() {
                document.getElementById('client-ip').value = device.ip;
                document.getElementById('client-port').value = device.port;
                document.getElementById('connect-button').click();
            });

            devicesContainer.appendChild(card);
        });

        // Show the container
        container.classList.remove('hidden');
    }

    function updateConnectionInfo() {
        if (!activeConnection) return;

        document.getElementById('connection-status').textContent = activeConnection.status === 'active' ? 'Connected' : 'Disconnected';
        document.getElementById('connection-mode').textContent = activeConnection.mode === 'server' ? 'Server' : 'Client';
        document.getElementById('connection-ip').textContent = activeConnection.ip;
        document.getElementById('connection-port').textContent = activeConnection.port;
    }

    function startWebSocketConnection() {
        // Close existing connection if any
        if (window.connectWebSocket) {
            window.connectWebSocket.close();
        }

        // Get authentication token if available
        const token = getAuthToken() || '';

        // Create new WebSocket connection with authentication token if available
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = token
            ? `${protocol}//${window.location.host}/api/v1/connect/ws?token=${token}`
            : `${protocol}//${window.location.host}/api/v1/connect/ws`;

        const ws = new WebSocket(wsUrl);

        ws.onopen = function() {
            console.log('WebSocket connection established');
        };

        ws.onmessage = function(event) {
            const data = JSON.parse(event.data);

            if (data.type === 'file') {
                // Handle received file
                addFileTransferHistoryEntry({
                    filename: data.filename,
                    size: data.size,
                    direction: 'received',
                    status: 'success',
                    time: new Date()
                });
            }
        };

        ws.onclose = function() {
            console.log('WebSocket connection closed');
        };

        ws.onerror = function(error) {
            console.error('WebSocket error:', error);
        };

        // Store the WebSocket connection
        window.connectWebSocket = ws;
    }

    function addFileTransferHistoryEntry(entry) {
        // Add to history array
        fileTransferHistory.push(entry);

        // Update UI
        const historyContainer = document.getElementById('file-history');
        const row = document.createElement('tr');

        row.innerHTML = `
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">${entry.filename}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${formatFileSize(entry.size)}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${entry.direction === 'sent' ? 'Sent' : 'Received'}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                <span class="status-${entry.status}">${entry.status === 'success' ? '✓ Success' : entry.status === 'error' ? '✗ Failed' : '⟳ Pending'}</span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${formatTime(entry.time)}</td>
        `;

        historyContainer.prepend(row);

        // Show the container
        document.getElementById('file-history-container').classList.remove('hidden');
    }

    function formatFileSize(size) {
        const units = ['B', 'KB', 'MB', 'GB', 'TB'];
        let i = 0;
        while (size >= 1024 && i < units.length - 1) {
            size /= 1024;
            i++;
        }
        return `${size.toFixed(1)} ${units[i]}`;
    }

    function formatTime(time) {
        return time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    }

    // Function to handle authentication errors
    function handleAuthError() {
        // Clear authentication data
        localStorage.removeItem('lumo_token');
        localStorage.removeItem('lumo_refresh_token');
        localStorage.removeItem('lumo_username');
        localStorage.removeItem('lumo_token_expiry');
        // Store current URL for redirect after login
        localStorage.setItem('lumo_redirect_after_login', window.location.href);
        // Redirect to main page
        window.location.href = '../';
    }

    // Function to send large files using chunked transfer
    async function sendLargeFileWithChunkedTransfer(file) {
        try {
            // Create progress bar if it doesn't exist
            let progressContainer = document.querySelector('.progress-container');
            let progressBar = document.querySelector('.progress-bar');

            if (!progressContainer) {
                progressContainer = document.createElement('div');
                progressContainer.className = 'progress-container';
                progressBar = document.createElement('div');
                progressBar.className = 'progress-bar';
                progressBar.style.width = '0%';
                progressContainer.appendChild(progressBar);

                // Insert after the file upload area
                const fileUploadArea = document.querySelector('.file-upload-area');
                fileUploadArea.parentNode.insertBefore(progressContainer, fileUploadArea.nextSibling);
            } else {
                progressBar.style.width = '0%';
                progressContainer.classList.remove('hidden');
            }

            // Step 1: Initialize upload
            const initResponse = await fetch('/api/v1/connect/upload/init', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    filename: file.name,
                    file_size: file.size
                })
            });

            if (!initResponse.ok) {
                throw new Error(`Failed to initialize upload: ${await initResponse.text()}`);
            }

            const uploadData = await initResponse.json();
            if (!uploadData || !uploadData.success) {
                throw new Error((uploadData && uploadData.error) || 'Failed to initialize upload');
            }

            if (!uploadData.upload_id || !uploadData.chunks) {
                throw new Error('Invalid response from server: missing upload ID or chunks information');
            }

            const uploadId = uploadData.upload_id;
            const chunkSize = uploadData.chunk_size || DefaultChunkSize;
            const totalChunks = uploadData.chunks.length;

            console.log(`Upload initialized with ID: ${uploadId}`);
            console.log(`Chunk size: ${chunkSize} bytes`);
            console.log(`Total chunks: ${totalChunks}`);

            // Step 2: Upload chunks
            for (let chunkId = 0; chunkId < totalChunks; chunkId++) {
                // Calculate progress
                const progress = Math.floor((chunkId + 1) * 100 / totalChunks);
                progressBar.style.width = `${progress}%`;

                // Calculate chunk boundaries
                const start = chunkId * chunkSize;
                const end = Math.min(start + chunkSize, file.size);

                // Read chunk
                const chunk = file.slice(start, end);

                // Upload chunk
                const chunkResponse = await fetch(
                    `/api/v1/connect/upload/chunk?upload_id=${uploadId}&chunk_id=${chunkId}`,
                    {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/octet-stream'
                        },
                        body: chunk
                    }
                );

                if (!chunkResponse.ok) {
                    throw new Error(`Failed to upload chunk ${chunkId}: ${await chunkResponse.text()}`);
                }

                const chunkResult = await chunkResponse.json();
                if (!chunkResult || !chunkResult.success) {
                    throw new Error((chunkResult && chunkResult.error) || `Failed to upload chunk ${chunkId}`);
                }
            }

            // Step 3: Complete upload
            const completeResponse = await fetch(
                `/api/v1/connect/upload/complete?upload_id=${uploadId}`,
                {
                    method: 'POST'
                }
            );

            if (!completeResponse.ok) {
                throw new Error(`Failed to complete upload: ${await completeResponse.text()}`);
            }

            const completeResult = await completeResponse.json();
            if (!completeResult || !completeResult.success) {
                throw new Error((completeResult && completeResult.error) || 'Failed to complete upload');
            }

            console.log(`Upload completed successfully!`);
            const filePath = completeResult.file_path || 'Unknown location';
            console.log(`File saved to: ${filePath}`);

            // Hide progress bar after completion
            setTimeout(() => {
                progressContainer.classList.add('hidden');
            }, 1000);

            // Add to file transfer history
            addFileTransferHistoryEntry({
                filename: file.name,
                size: file.size,
                direction: 'sent',
                status: 'success',
                time: new Date()
            });

            // Reset file upload
            fileToSend = null;
            document.getElementById('file-upload').value = '';
            document.querySelector('.file-upload-area p').textContent = 'or drag and drop';

            return filePath;
        } catch (error) {
            console.error('Chunked upload error:', error);

            // Hide progress bar on error
            const progressContainer = document.querySelector('.progress-container');
            if (progressContainer) {
                progressContainer.classList.add('hidden');
            }

            // Add to file transfer history with error status
            addFileTransferHistoryEntry({
                filename: file.name,
                size: file.size,
                direction: 'sent',
                status: 'error',
                time: new Date()
            });

            throw error;
        }
    }
});
