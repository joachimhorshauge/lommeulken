function setTheme(theme) {
  const htmlElement = document.documentElement;

  htmlElement.classList.remove('light', 'dark');

  if (theme !== 'light') {
    htmlElement.classList.add(theme);
  }

  localStorage.setItem('theme', theme);
}

function initializeTheme() {
  const savedTheme = localStorage.getItem('theme') || 'light';
  setTheme(savedTheme);
}

function attachThemeSelectListener() {
  const themeSelect = document.getElementById('theme-select');
  if (themeSelect) {
    themeSelect.addEventListener('change', function (event) {
      setTheme(event.target.value);
    });

    themeSelect.value = localStorage.getItem('theme') || 'light';
  }
}

document.addEventListener('DOMContentLoaded', () => {
  initializeTheme(); 
  attachThemeSelectListener(); 
});
