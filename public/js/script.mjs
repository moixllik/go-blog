class WebComponent extends HTMLElement {
    constructor() {
        super();
        this.render()
    }
    async render() {
        const res = await fetch('/html/' + this.getAttribute('uri') + '.html');
        if (res.status != 200) return;
        const html = await res.text();
        const shadow = this.attachShadow({mode: 'closed'});
        shadow.innerHTML = html;
    }
}

customElements.define('web-component', WebComponent);
