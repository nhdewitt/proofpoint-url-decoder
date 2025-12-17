// Fallback for older browsers
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
    if (!tooltip) return;

    // Check if we are in "mobile mode" (screen width < 768px)
    // This decides if we use mouse coordinates or fixed toast position
    const isMobile = window.matchMedia("(max-width: 768px)").matches;

    document.querySelectorAll('.output').forEach(el => {
        el.style.cursor = 'pointer';
        
        el.onclick = async e => {
            const text = el.textContent.trim();
            
            // 1. Perform Copy
            if (navigator.clipboard?.writeText) {
                try { await navigator.clipboard.writeText(text); } 
                catch { fallbackCopy(text); }
            } else {
                fallbackCopy(text);
            }

            // 2. Show Tooltip
            
            // Reset animation
            tooltip.classList.remove('show');
            tooltip.style.opacity = '0';
            
            // Small delay to allow class removal to register
            requestAnimationFrame(() => {
                tooltip.classList.add('show');
                tooltip.style.opacity = '1';

                // Only set X/Y coordinates if we are on Desktop
                if (!isMobile) {
                    tooltip.style.left = e.pageX + 'px';
                    tooltip.style.top  = e.pageY + 'px';
                    tooltip.style.transform = 'translate(-50%, -150%)'; // Center above cursor
                } else {
                    // Mobile: Ensure we clear inline styles so CSS class controls position
                    tooltip.style.left = '';
                    tooltip.style.top = '';
                    tooltip.style.transform = ''; 
                }
            });

            // 3. Hide after delay
            setTimeout(() => {
                tooltip.classList.remove('show');
                tooltip.style.opacity = '0';
            }, 2000);
        };
    });
}

// Make globally available so json-api-request.js can call it
window.attachCopyHandlers = attachCopyHandlers;
document.addEventListener('DOMContentLoaded', attachCopyHandlers);