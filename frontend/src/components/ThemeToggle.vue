<template>
  <button class="theme-toggle" :title="isDarkMode ? 'Cambiar a modo claro' : 'Cambiar a modo oscuro'" @click="toggleDarkMode">
    <i class="icon" :class="isDarkMode ? 'light-mode-icon bx  bx-sun' : 'dark-mode-icon bx  bx-moon'"></i>
    <span>{{ isDarkMode ? 'Modo claro' : 'Modo oscuro' }}</span>
  </button>
</template>

<script setup>
  import { inject } from 'vue';

  const isDarkMode = inject('isDarkMode');
  const toggleDarkMode = inject('toggleDarkMode');
</script>

<style lang="scss">
@import '../styles/variables.scss';

.theme-toggle {
  display: flex;
  align-items: center;
  background: none;
  border: none;
  color: $text-color-light-mode;
  cursor: pointer;
  font-size: 0.9rem;
  padding: 0.5rem;

  .dark-mode & {
    color: $text-color-dark-mode;
  }

  .icon {
    margin-right: 4px;
    font-size: 1.1rem;
  }

  &:hover {
    color: $primary-color-light-mode;

    .dark-mode & {
      color: $primary-color-dark-mode;
    }
  }
}
</style>
