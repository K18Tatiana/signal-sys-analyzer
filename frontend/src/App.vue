<!-- App.vue -->
<template>
  <div class="app" :class="{ 'dark-mode': isDarkMode }">
    <AppHeader />
    <main class="main-content">
      <router-view />
    </main>
    <AppFooter />
  </div>
</template>

<script setup>
  import { provide } from 'vue';
  import { useDarkMode } from './composables/useDarkMode';
  import AppHeader from './components/AppHeader.vue';
  import AppFooter from './components/AppFooter.vue';

  const { isDarkMode, toggleDarkMode } = useDarkMode();
  provide('isDarkMode', isDarkMode);
  provide('toggleDarkMode', toggleDarkMode);
</script>

<style lang="scss">
@import './styles/variables.scss';
@import './styles/globals.scss';

.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: $background-color-light-mode;
  color: $text-color-light-mode;
  transition: background-color 0.3s, color 0.3s;

  &.dark-mode {
    background-color: $background-color-dark-mode;
    color: $text-color-dark-mode;
  }

  .main-content {
    flex: 1;
    padding: 2rem;
    padding-top: 3.5rem;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
}
</style>
