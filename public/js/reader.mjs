class MyDate extends HTMLElement {
    constructor() {
        super();
        this.innerHTML = this.innerHTML.slice(0, 10);
    }
}

class MyMarkdown extends HTMLElement {
    constructor() {
        super();
        this.matches = [];
    }
    connectedCallback() {
        this.add(/^\#\# (.*)/gm, '<h2>$1</h2>');
        this.add(/^\#\#\# (.*)/gm, '<h3>$1</h3>');

        this.add(/^\`\`\`$/gm, '\n</my-code></pre>');
        this.add(/^\`\`\`(.*)/gm, '<pre><my-code data-lang="$1">');

        this.add(/\`(.*?)\`/g, '<code>$1</code>');

        this.add(/\!\[(.*?)\]\((.*?)\)/g, '<figure><img alt="$1" src="$2" width="100%" height="100%"></figure>');

        this.add(/\[(.*?)\]\((.*?)\)/g, '<a href="$2" target="_blank">$1</a>');

        this.format();
    }
    add(match, result) {
        this.matches.push([match, result]);
    }
    format() {
        let text = this.innerText;
        this.matches.forEach(it => {
            text = text.replace(it[0], it[1]);
        });
        this.innerHTML = text;
    }
}

class MyCode extends HTMLElement {
    constructor() {
        super();
        this.matches = [];
        this.keywords = false;
        this.language = this.getAttribute('data-lang');
        this.text = this.innerText;
        this.regexp = ['\\b(', ')\\b', 'g'];
    }
    connectedCallback() {
        this.loadMatches();
        this.format();
        this.copyTool();
    }
    add(match, result) {
        this.matches.push([match, result]);
    }
    format() {
        if (this.text.length == 0) return;
        this.matches.forEach(it => {
            this.text = this.text.replace(it[0], it[1]);
        });
        this.innerHTML = this.text;
    }
    loadMatches() {
        this['lang_' + this.language]?.apply(this);
        if (this.keywords) this.keywords.split(' ').forEach(it => {
            if (it.length == 0) return;
            this.add(new RegExp(this.regexp[0] + it + this.regexp[1], this.regexp[2]), '<b>' + it + ' </b>');
        });
    }
    copyTool() {
        const $elm = document.createElement('div');
        const text = this.innerText.trim();
        $elm.className = 'copy';
        $elm.addEventListener('click', _ => {
            navigator.clipboard.writeText(text).then(_ => {
                $elm.style.color = 'greenyellow';
            });
        });
        this.parentNode.appendChild($elm);
    }
    lang_bash() {
        this.add(/^(.*)\#(.*)/gm, '$1<i>#$2</i>');
        this.add(/\"(.*?)\"/g, '<b>"$1"</b>');
    }
    lang_v() {
        this.keywords = 'module import fn';
        this.add(/^\/\/(.*)/gm, '<i>//$1</i>');
        this.add(/\'(.*?)\'/g, '<span>\'$1\'</span>');
    }
    lang_yml() {
        this.add(/^(.*):/gm, '<span>$1:</span>');
    }
    lang_docker() {
        this.keywords = 'FROM RUN WORKDIR COPY ENTRYPOINT'
    }
    lang_json() {
        this.add(/\"(.*?)\"/g, '<span>"$1"</span>');
    }
    lang_math() {
        this.style.fontStyle = 'italic';
        this.add(/([^\p{L}\p{N}\s\n])/gu, '<em>$1</em>');
    }
    lang_python() {
        this.keywords = 'print lambda def import form return class if for range';
    }
    lang_javascript() {
        this.keywords = 'class var const let extends function if for while';
    }
    lang_toml() {
        this.add(/^\[(.*)\]$/gm, '<b>[$1]</b>');
        this.add(/\"(.*?)\"/g, '<span>"$1"</span>');
    }
    lang_table() {
        const $table = document.createElement('table');
        const align = [];
        this.text.split('\n').forEach(row => {
            const rows = row.trim().split('|');
            if (rows.length == 0) return;
            const $tr = document.createElement('tr');
            const elm = ('-:.'.includes(row.trim()[1])) ? 'th' : 'td';
            rows.forEach((column, idx) => {
                const field = column.trim();
                if (field.length == 0) return;
                const $elm = document.createElement(elm);
                const mod = ('-:.'.includes(field[0])) ? field[0] : false;

                if (mod == '-') align[idx] = 'left';
                else if (mod == ':') align[idx] = 'center';
                else if (mod == '.') align[idx] = 'right';

                $elm.style.textAlign = align[idx];

                if (mod) $elm.innerHTML = field.substring(1).trim();
                else $elm.innerHTML = field;
                $tr.appendChild($elm);
            });
            if ($tr.children.length > 0) $table.appendChild($tr);
        });
        if ($table.children.length > 0) {
            this.text = '';
            this.innerHTML = '';
            this.appendChild($table);
        }
    }
}

customElements.define('my-date', MyDate);
customElements.define('my-code', MyCode);
customElements.define('my-markdown', MyMarkdown);
