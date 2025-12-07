// document.documentElement.classList.add('dark');

document.addEventListener('DOMContentLoaded', () => {
    const formGrab = document.querySelector('#form-grab');
    const buttonGrab = document.querySelector('.button-grab-get');
    const inputURL = document.querySelector('#youtubeURL');
    const resultDivInfo = document.querySelector('#grab-result-info');
    const radios = document.querySelectorAll('input[name="formatType"]');

    //-----------------------------------------------------------------
    //  Restore saved radio selection (formatType)
    //-----------------------------------------------------------------
    const getCookie = (name) => {
        return document.cookie
            .split('; ')
            .find(row => row.startsWith(name + "="))
            ?.split('=')[1];
    };

    const savedFormat = getCookie('formatType');
    if (savedFormat) {
        const savedRadio = document.querySelector(`input[name="formatType"][value="${savedFormat}"]`);
        if (savedRadio) savedRadio.checked = true;
    }

    //-----------------------------------------------------------------
    //  Save radio choice on change
    //-----------------------------------------------------------------
    radios.forEach(radio => {
        radio.addEventListener('change', () => {
            const date = new Date();
            date.setFullYear(date.getFullYear() + 1); // store 1 year
            document.cookie = "formatType=" + encodeURIComponent(radio.value)
                + ";expires=" + date.toUTCString()
                + ";path=/";
        });
    });

    //-----------------------------------------------------------------
    //  Submit on Enter
    //-----------------------------------------------------------------
    inputURL.addEventListener('keydown', (event) => {
        if (event.key === 'Enter') {
        event.preventDefault();
        buttonGrab.click();
        }
    });

    //-----------------------------------------------------------------
    //  Before HTMX Request — clear input + result
    //-----------------------------------------------------------------
    htmx.on('#form-grab', 'htmx:beforeRequest', () => {
        if (inputURL) inputURL.value = '';
        if (resultDivInfo) resultDivInfo.innerHTML = '';
    });

    //-----------------------------------------------------------------
    //  After HTMX Response — display error if not 200
    //-----------------------------------------------------------------
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

    //-----------------------------------------------------------------
    //  Save "resultItemsOnlyOne" cookie before any HTMX request
    //-----------------------------------------------------------------
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

