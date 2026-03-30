document.addEventListener('DOMContentLoaded', () => {
    const authForm = document.getElementById('authForm');
    const password = document.getElementById('password');
    const confirmPassword = document.getElementById('confirmPassword');
    const errorMessage = document.getElementById('error');

    // Display error on non-200 + non-503
    document.body.addEventListener('htmx:afterOnLoad', (event) => {
        if (event.detail.elt === authForm) {
            if (event.detail.xhr.status !== 200 &&
                event.detail.xhr.status !== 503) {
                if (errorMessage) {
                    let text = event.detail.xhr.responseText;
                    try {
                        const data = JSON.parse(text);
                        if (data && typeof data === "object" && "message" in data) {
                            errorMessage.textContent = data.message;
                        }
                    } catch (e) {}
                }
            }
        }
    });

    if (confirmPassword) {
        authForm.addEventListener('submit', function(e) {
            if (password.value !== confirmPassword.value) {
                e.preventDefault();
                errorMessage.textContent = "Passwords do not match";
            }
        });
    }
});