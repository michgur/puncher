function clientY(e) {
    return e.clientY || e.changedTouches[0].clientY;
}

class BottomSheet extends HTMLElement {
    static observedAttributes = ['content'];

    connectedCallback() {
        const template = document.getElementById('bottom-sheet-template');
        this.appendChild(template.content.cloneNode(true));
        this.modal = this.firstElementChild;
        this.modalBody = this.querySelector('.modal-body');
        this.attributeChangedCallback('content', null, this.getAttribute('content'));
        this.scrim = this.querySelector('.scrim');
        this.scrim.addEventListener('mousedown', (e) => {
            e.preventDefault();
            this.close();
        });
        this.scrim.addEventListener('touchstart', (e) => {
            e.preventDefault();
            this.close();
        });
        this.grabber = this.querySelector('.grabber');
        this.grabber.addEventListener('click', this.close.bind(this));
        this.grabber.addEventListener('mousedown', this.grabberMouseDown.bind(this));
        this.grabber.addEventListener('touchstart', this.grabberMouseDown.bind(this));
    }

    grabberMouseDown(e) {
        const startY = clientY(e);
        const mouseMove = (e) => {
            const delta = Math.max((clientY(e) - startY) * 0.6, -20);
            this.modal.style.transform = `translateY(${delta}px)`;
        };
        const mouseUp = (e) => {
            this.modal.style.transform = '';
            if (clientY(e) - startY > 75) {
                this.close();
            }
            document.removeEventListener('mousemove', mouseMove);
            document.removeEventListener('mouseup', mouseUp);
            document.removeEventListener('touchmove', mouseMove);
            document.removeEventListener('touchend', mouseUp);
        };
        document.addEventListener('mousemove', mouseMove);
        document.addEventListener('mouseup', mouseUp);
        document.addEventListener('touchmove', mouseMove);
        document.addEventListener('touchend', mouseUp);
    }

    attributeChangedCallback(name, oldValue, newValue) {
        if (this.modalBody && name === 'content' && newValue !== oldValue) {
            const content = document.querySelector(newValue);
            this.modalBody.appendChild(content);
        }
    }

    open() {
        this.modal.ariaHidden = false;
        setTimeout(() => {
            this.modal.ariaExpanded = true;
        });
    }

    close() {
        this.modal.ariaExpanded = false;
        setTimeout(() => {
            this.modal.ariaHidden = true;
        }, 300)
    }
}

customElements.define('bottom-sheet', BottomSheet);
