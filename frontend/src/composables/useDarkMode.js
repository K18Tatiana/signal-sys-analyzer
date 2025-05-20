import { onMounted, ref, watch } from 'vue';

export function useDarkMode () {
  const isDarkMode = ref(false);

  const toggleDarkMode = () => {
    isDarkMode.value = !isDarkMode.value;
    localStorage.setItem('dark-mode', isDarkMode.value);
    updateBodyClass();
  };

  const updateBodyClass = () => {
    if (isDarkMode.value) {
      document.body.classList.add('dark-mode');
    } else {
      document.body.classList.remove('dark-mode');
    }
  };

  onMounted(() => {
    const savedMode = localStorage.getItem('dark-mode');

    if (savedMode !== null) {
      isDarkMode.value = savedMode === 'true';
    } else {
      const prefersDarkMode = window.matchMedia('(prefers-color-scheme: dark)');
      isDarkMode.value = prefersDarkMode.matches;

      prefersDarkMode.addEventListener('change', e => {
        isDarkMode.value = e.matches;
        localStorage.setItem('dark-mode', isDarkMode.value);
      });
    }

    updateBodyClass();
  });

  watch(isDarkMode, () => {
    updateBodyClass();
  });

  return {
    isDarkMode,
    toggleDarkMode,
  };
}
