(() => {
    function nextDigit(el) {
        if (el.nextElementSibling?.tagName === 'INPUT') {
            el.nextElementSibling.focus();
            return el.nextElementSibling;
        }
        el.blur();
    }

    function previousDigit(el) {
        if (el.previousElementSibling?.tagName === 'INPUT') {
            el.previousElementSibling.focus();
            return el.previousElementSibling;
        }
        el.blur();
    }

    class OtpInput extends HTMLElement {
        trySubmit() {
            if (this.digits.every(el => el.value !== '')) {
                this.form.requestSubmit();
                this.digits.forEach(el => el.disabled = true);
            }
        }

        otpKeyPress(e) {
            // custom handling - prevent default
            e.preventDefault();
            // not a digit - ignore
            if (!e.key.match(/^\d$/)) {
                return;
            }

            // fill input with digit
            e.target.value = e.key;
            // form state changed - remove error
            this.form.classList.remove('error');

            if (nextDigit(e.target) === undefined) {
                this.trySubmit();
            }
        }


        otpKeyDown(e) {
            switch (e.code) {
                case 'ArrowLeft':
                    previousDigit(e.target);
                    break;
                case 'ArrowRight':
                    nextDigit(e.target);
                    break;
                case 'Backspace':
                    if (e.target.value === '') {
                        // current digit is empty - go back
                        previousDigit(e.target);
                    } else {
                        // clear current digit
                        e.target.value = '';
                    }
                    // form state changed - remove error
                    this.form.classList.remove('error');
                    break;
            }
        }

        otpPaste(e) {
            // custom handling - prevent default
            e.preventDefault();
            // get clipboard data
            const data = e.clipboardData.getData('text');
            // fill inputs with clipboard data
            let i = 0;
            let digit = e.target;
            while (
                i < data.length &&
                data[i].match(/\d/) &&
                digit !== undefined
            ) {
                digit.value = data[i];
                digit = nextDigit(digit);
                i++;
            }
            // form state changed - remove error
            this.form.classList.remove('error');
            // try to submit
            this.trySubmit();
        }

        connectedCallback() {
            // # of digits in OTP
            const numberOfDigits = (+this.getAttribute('digits')) || 4;
            // fill with existing value (for presisting state)
            const value = this.getAttribute('value') || '';
            // HTTP GET request path
            let path = (this.getAttribute('path') || '') + "/";

            // get digit & form templates
            const formTemplate = document.querySelector('template#otp-form-template');
            const digitTemplate = document.querySelector('template#otp-digit-template');

            // create form from template
            this.appendChild(formTemplate.content.cloneNode(true));
            this.form = this.querySelector('form');
            // add error class if needed
            if (this.hasAttribute('error')) this.form.classList.add('error');

            // create digits
            this.digits = [];
            for (let i = 0; i < numberOfDigits; i++) {
                // create digit input
                const input = digitTemplate.content.cloneNode(true).firstElementChild;
                input.name = `digit${i}`;
                // bind event listeners
                input.addEventListener('keydown', this.otpKeyDown.bind(this));
                input.addEventListener('keypress', this.otpKeyPress.bind(this));
                input.addEventListener('paste', this.otpPaste.bind(this));
                // add input value to path variable
                path += `{${input.name}}`;
                // prefill input from value attribute
                if (value.length > i) input.value = value[i];
                // autofocus first empty input
                else if (value.length === i) input.autofocus = true;
                // append input to form
                this.form.appendChild(input);
                this.digits.push(input);
            }

            // update hx-get path to use digit inputs
            this.form.setAttribute('hx-get', path);
        }
    }

    customElements.define('otp-input', OtpInput);
})();