const root = document.documentElement;
const saved = localStorage.getItem('theme');
if (saved) root.dataset.theme = saved;
else if (matchMedia('(prefers-color-scheme: dark)').matches) root.dataset.theme = 'dark';

document.getElementById('theme').addEventListener('click', () => {
  const next = root.dataset.theme === 'dark' ? 'light' : 'dark';
  root.dataset.theme = next;
  localStorage.setItem('theme', next);
});

document.getElementById('year').textContent = new Date().getFullYear();
