document.addEventListener('DOMContentLoaded', () => {
    const form  = document.querySelector('.decode-form');
    const input = document.getElementById('input');
    const list  = document.getElementById('results');
    const card  = document.querySelector('.card'); // Used to show/hide empty states

    if (!form) return; // Exit if we are on a page without the form

    form.addEventListener('submit', async e => {
        e.preventDefault();

        // 1. Prepare data
        const urls = input.value
            .split('\n')
            .map(s => s.trim())
            .filter(Boolean);

        if (urls.length === 0) return;

        // 2. Fetch
        const res = await fetch('/api/decode', {
            method:     'POST',
            headers:    {'Content-Type':'application/json'},
            body:       JSON.stringify({urls}),
        });

        if (!res.ok) {
            list.innerHTML = `<li class="error">Server error: ${res.statusText}</li>`;
            list.style.display = 'block'; // Ensure list is visible
            return;
        }

        const {results, errors} = await res.json();

        // 3. Clear and Rebuild List
        list.innerHTML = '';
        list.style.display = 'flex'; // Make sure the list is visible (CSS :empty hides it)

        results.forEach((link, i) => {
            const li = document.createElement('li');
            li.classList.add('result-item');
            
            // A. Create "Source Input" Label
            const srcLabel = document.createElement('div');
            srcLabel.className = 'input-url';
            srcLabel.textContent = 'Source Input';
            srcLabel.style.borderBottom = 'none'; // Visual tweak
            srcLabel.style.marginBottom = '0';
            li.appendChild(srcLabel);

            // B. Create Input Value
            const inputDiv = document.createElement('div');
            inputDiv.className = 'input-url';
            inputDiv.style.color = 'var(--text-main)';
            inputDiv.style.marginBottom = '1rem';
            inputDiv.textContent = urls[i];
            li.appendChild(inputDiv);

            // C. Handle Error vs Success
            if (errors && errors[i]) {
                const err = document.createElement('div');
                err.classList.add('error');
                err.textContent = `Error: ${errors[i]}`;
                li.appendChild(err);
            } else {
                // Label
                const outLabel = document.createElement('div');
                outLabel.className = 'input-url';
                outLabel.textContent = 'Decoded Output';
                outLabel.style.borderBottom = 'none';
                outLabel.style.marginBottom = '0';
                li.appendChild(outLabel);

                // Output Block
                const pre = document.createElement('pre');
                pre.classList.add('output');
                pre.textContent = link;
                li.appendChild(pre);
            }

            list.appendChild(li);
        });

        // 4. Cleanup
        input.value = '';
        
        // Re-attach copy handlers to the new elements
        if (window.attachCopyHandlers) window.attachCopyHandlers();
    });
});