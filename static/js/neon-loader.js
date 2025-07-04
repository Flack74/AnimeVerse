// Neon Loader Utility Functions

function createNeonLoader(text = 'Loading...') {
    return `
        <div class="neon-loader-container">
            <div class="neon-loader"></div>
            <p class="neon-loader-text">${text}</p>
        </div>
    `;
}

function showNeonLoader(elementId, text = 'Loading...') {
    const element = document.getElementById(elementId);
    if (element) {
        element.innerHTML = createNeonLoader(text);
    }
}

function showNeonOverlay(text = 'Loading...') {
    const overlay = document.createElement('div');
    overlay.id = 'neon-overlay';
    overlay.className = 'neon-loading-overlay';
    overlay.innerHTML = `
        <div class="neon-card-loading">
            <div class="neon-loader"></div>
            <p class="neon-loader-text">${text}</p>
        </div>
    `;
    document.body.appendChild(overlay);
}

function hideNeonOverlay() {
    const overlay = document.getElementById('neon-overlay');
    if (overlay) {
        overlay.remove();
    }
}