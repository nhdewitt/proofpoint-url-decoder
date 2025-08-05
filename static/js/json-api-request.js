document.addEventListener('DOMContentLoaded', () => {
    const form  = document.querySelector('.decode-form');
    const input = document.getElementById('input');
    const list  = document.getElementById('results');

    form.addEventListener('submit', async e => {
        e.preventDefault();

        const urls = input.value
            .split('\n')
            .map(s => s.trim())
            .filter(Boolean);

        const res = await fetch('/api/decode', {
            method:     'POST',
            headers:    {'Content-Type':'application/json'},
            body:       JSON.stringify({urls}),
        });

        if (!res.ok) {
            list.innerHTML = `<li class="error">Server error: ${res.statusText}</li>`;
            return;
        }

        const {results, errors} = await res.json();

        list.innerHTML = '';

        results.forEach((link, i) => {
            const li = document.createElement('li');
            li.classList.add('result-item');

            const inputDiv = document.createElement('div');
            inputDiv.classList.add('input-url');
            inputDiv.textContent = urls[i];
            li.appendChild(inputDiv);

            if (errors && errors[i]) {
                const err = document.createElement('div');
                err.classList.add('error');
                err.textContent = errors[i];
                li.appendChild(err);
            } else {
                const pre = document.createElement('pre');
                pre.classList.add('output');
                pre.textContent = link;
                li.appendChild(pre);
            }

            list.appendChild(li);
        });

        input.value = '';
        const card = document.querySelector('.card.container');
        if (card) card.classList.add('has-results');

        window.attachCopyHandlers?.();
    });
});