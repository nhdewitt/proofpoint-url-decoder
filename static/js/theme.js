document.addEventListener('DOMContentLoaded', () => {
  const toggle     = document.getElementById('themeToggle');
  const stored     = localStorage.getItem('theme');
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  const theme      = stored || (prefersDark ? 'dark' : 'light');

  if (theme === 'dark') {
    document.body.classList.add('dark');
    if (toggle) toggle.checked = true;
  }

  if (toggle) {
    toggle.addEventListener('change', () => {
      if (toggle.checked) {
        document.body.classList.add('dark');
        localStorage.setItem('theme', 'dark');
      } else {
        document.body.classList.remove('dark');
        localStorage.setItem('theme', 'light');
      }
    });
  }
});