document.addEventListener('DOMContentLoaded', () => {
    const button = document.querySelector('.button-grab-get');
    const input = document.querySelector('#youtubeURL');
    const resultDiv = document.querySelector('#grab-result');

    // Listen for Enter key inside input
    input.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
        event.preventDefault();
        button.click();
        }
    });

    button.addEventListener('htmx:configRequest', (event) => {
        const button = event.target;
        const input = document.querySelector('#youtubeURL');
        const resultDiv = document.querySelector('#grab-result');

        if (button.disabled) return;

        // Disable button and input during loading
        button.disabled = true;
        input.disabled = true;

        // Show loading indicator
        if (resultDiv) {
            // Remove old spinners if any (prevent duplicates)
            const oldSpinner = resultDiv.querySelector('.result-loading');
            if (oldSpinner) oldSpinner.remove();

            // Append spinner *inside* existing content
            resultDiv.innerHTML = `
                <div class="result-loading">
                    <div class="result-spinner"></div>
                    <span class="result-loading-text">Loading...</span>
                </div>
            `;
        }
    });

    document.body.addEventListener('htmx:afterOnLoad', (event) => {
        button.disabled = false;
        input.disabled = false;

        if (event.detail.xhr.status === 200) {
            if (input) input.value = '';
        } else {
            // Show error message inside result div
            if (resultDiv) {
                resultDiv.innerHTML = `
                    <div class="grab-result-item">    
                        <span class="result-error">Error: ${event.detail.xhr.responseText}</span>
                    </div>
                `;
            } 
        }
    });
});

