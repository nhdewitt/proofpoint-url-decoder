function fallbackCopy(str) {
const ta = document.createElement('textarea');
ta.value = str;
Object.assign(ta.style, {
    position:'fixed', top:0, left:0,
    width:'1px', height:'1px',
    padding:0, border:'none', outline:'none', opacity:0
});
document.body.appendChild(ta);
ta.focus(); ta.select();
document.execCommand('copy');
document.body.removeChild(ta);
}

function attachCopyHandlers() {
    const tooltip = document.getElementById('copyTooltip');
    document.querySelectorAll('pre.output').forEach(el => {
        el.style.cursor = 'pointer';
        el.onclick = async e => {
            const text = el.textContent.trim();
            if (navigator.clipboard?.writeText) {
                try { await navigator.clipboard.writeText(text) }
                catch { fallbackCopy(text) }
            } else {
                fallbackCopy(text);
            }
            tooltip.textContent     = '✅ Copied to clipboard!';
            tooltip.style.left      = e.pageX + 'px';
            tooltip.style.top       = e.pageY + 'px';
            tooltip.style.opacity   = '1';
            setTimeout(() => { tooltip.style.opacity = '0' }, 1500);
        }
    });
}

window.attachCopyHandlers = attachCopyHandlers;

document.addEventListener('DOMContentLoaded', attachCopyHandlers);