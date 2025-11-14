document.addEventListener('DOMContentLoaded', () => {
    const formGrab = document.querySelector('#form-grab');
    const buttonGrab = document.querySelector('.button-grab-get');
    const inputURL = document.querySelector('#youtubeURL');
    const resultDivInfo = document.querySelector('#grab-result-info');

    // Listen for Enter key inside input
    inputURL.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
        event.preventDefault();
        buttonGrab.click();
        }
    });

    htmx.on('#form-grab', 'htmx:beforeRequest', () => {
        if (inputURL) inputURL.value = '';
        if (resultDivInfo) resultDivInfo.innerHTML = '';
    });

    document.body.addEventListener('htmx:afterOnLoad', (event) => {
        if (event.detail.elt === formGrab) {
            if (inputURL) inputURL.value = '';
            if (event.detail.xhr.status !== 200) {
                if (resultDivInfo) {
                    resultDivInfo.innerHTML = `
                        <div class="div-grab-result-info-item">    
                            <span class="result-failed">Error: ${event.detail.xhr.responseText}</span>
                        </div>
                    `;
                } 
            }
        }
    });

    document.body.addEventListener('htmx:beforeRequest', function(evt) {
        const div = document.getElementById("grab-result-item-replaceme");
        if (div) {
            const onlyOne = div.dataset.onlyOne;
            const d = new Date();
            d.setTime(d.getTime() + (7*24*60*60*1000));
            document.cookie = "resultItemsOnlyOne=" + encodeURIComponent(onlyOne) + ";expires=" + d.toUTCString() + ";path=/";
        }
    });

});

