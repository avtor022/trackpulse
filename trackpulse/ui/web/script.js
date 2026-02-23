// WebSocket connection for live updates
let socket;
const leaderboardBody = document.getElementById('leaderboard-body');
const statusValue = document.getElementById('status-value');
const timerValue = document.getElementById('timer-value');

function initWebSocket() {
    // Connect to the WebSocket server
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUri = `${protocol}//${window.location.host}/ws`;
    
    socket = new WebSocket(wsUri);
    
    socket.onopen = function(event) {
        console.log('Connected to WebSocket server');
    };
    
    socket.onmessage = function(event) {
        try {
            const data = JSON.parse(event.data);
            updateDisplay(data);
        } catch (e) {
            console.error('Error parsing WebSocket message:', e);
        }
    };
    
    socket.onclose = function(event) {
        console.log('Disconnected from WebSocket server. Attempting to reconnect...');
        // Reconnect after 3 seconds
        setTimeout(initWebSocket, 3000);
    };
    
    socket.onerror = function(error) {
        console.error('WebSocket error:', error);
    };
}

function updateDisplay(data) {
    if (data.type === 'race_update') {
        updateRaceInfo(data.data);
    } else if (data.event === 'lap_update') {
        updateLapInfo(data);
    }
}

function updateRaceInfo(raceData) {
    // Update race title if available
    if (raceData.race_title) {
        document.getElementById('race-title').textContent = raceData.race_title;
    }
    
    // Update race status
    if (raceData.status) {
        statusValue.textContent = raceData.status;
        statusValue.className = `status-${raceData.status}`;
    }
    
    // Update timer if elapsed time is provided
    if (raceData.elapsed_ms !== undefined) {
        const elapsedSeconds = Math.floor(raceData.elapsed_ms / 1000);
        const minutes = Math.floor(elapsedSeconds / 60);
        const seconds = elapsedSeconds % 60;
        timerValue.textContent = `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    }
    
    // Update leaderboard if participants data is available
    if (raceData.participants) {
        updateLeaderboard(raceData.participants);
    }
}

function updateLapInfo(lapData) {
    // Update the specific participant's lap information
    // This could involve updating just one row in the leaderboard
    // For now, we'll trigger a full refresh of the leaderboard
    // In a real implementation, we'd find the specific row and update it
    console.log('Lap update received:', lapData);
}

function updateLeaderboard(participants) {
    // Sort participants by laps (descending) and best lap time (ascending)
    const sortedParticipants = [...participants].sort((a, b) => {
        // Primary sort: by number of laps (descending)
        if (b.number_of_laps !== a.number_of_laps) {
            return b.number_of_laps - a.number_of_laps;
        }
        // Secondary sort: by best lap time (ascending)
        return a.best_lap_time_ms - b.best_lap_time_ms;
    });
    
    // Clear the current leaderboard
    leaderboardBody.innerHTML = '';
    
    // Add each participant to the leaderboard
    sortedParticipants.forEach((participant, index) => {
        const row = document.createElement('tr');
        
        // Format lap times from milliseconds to MM:SS.mmm
        const bestLapFormatted = formatLapTime(participant.best_lap_time_ms);
        const lastLapFormatted = formatLapTime(participant.last_lap_time_ms);
        
        row.innerHTML = `
            <td>${participant.position || (index + 1)}</td>
            <td>${participant.racer_number}</td>
            <td>${participant.full_name}</td>
            <td>${participant.model_name}</td>
            <td>${participant.number_of_laps}</td>
            <td>${bestLapFormatted}</td>
            <td>${lastLapFormatted}</td>
        `;
        
        // Add position-based styling
        if (index === 0) row.classList.add('pos-1');
        else if (index === 1) row.classList.add('pos-2');
        else if (index === 2) row.classList.add('pos-3');
        
        leaderboardBody.appendChild(row);
    });
}

function formatLapTime(ms) {
    if (ms === 0 || ms === null || ms === undefined) {
        return '--:--.---';
    }
    
    const minutes = Math.floor(ms / 60000);
    const seconds = Math.floor((ms % 60000) / 1000);
    const milliseconds = ms % 1000;
    
    return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}.${milliseconds.toString().padStart(3, '0')}`;
}

// Initialize the page when loaded
document.addEventListener('DOMContentLoaded', function() {
    initWebSocket();
});