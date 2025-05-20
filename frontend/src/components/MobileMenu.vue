<template>
  <div class="mobile-menu" :class="{ 'is-open': isOpen, 'dark-mode': isDarkMode }">
    <div class="mobile-menu-header">
      <div class="logo">
        <img v-if="isDarkMode" alt="SignalSys Analyzer" src="../assets/logo-dark-mode.png" />
        <img v-else alt="SignalSys Analyzer" src="../assets/logo-light-mode.png" />
        <div class="logo-text">
          <span class="logo-name">SignalSys</span>
          <span class="logo-description">Analyzer</span>
        </div>
      </div>
      <button class="close-button" @click="$emit('close')">
        <span class="close-icon">✕</span>
      </button>
    </div>

    <nav class="mobile-nav">
      <ul class="mobile-nav-list">
        <li class="mobile-nav-item">
          <router-link to="/" @click="$emit('close')">
            <i class='icon bx bx-home-alt'></i> Inicio
          </router-link>
        </li>
        <li class="mobile-nav-item">
          <router-link to="/guide" @click="$emit('close')">
            <i class='icon bx bx-book-open'></i>  Guía / Tutorial
          </router-link>
        </li>
        <li class="mobile-nav-item">
          <router-link to="/language" @click="$emit('close')">
            <i class='icon bx bx-globe'></i>  Idioma
          </router-link>
        </li>
        <li class="mobile-nav-item">
          <button class="theme-toggle-mobile" @click="toggleTheme">
            <i class="icon" :class="isDarkMode ? 'light-mode-icon bx  bx-sun' : 'dark-mode-icon bx  bx-moon'"></i>
            {{ isDarkMode ? 'Modo claro' : 'Modo oscuro' }}
          </button>
        </li>
        <li class="mobile-nav-item">
          <router-link class="login-link" to="/login" @click="$emit('close')">
            <i class='icon bx bx-user'></i>  Iniciar sesión
          </router-link>
        </li>
      </ul>
    </nav>
  </div>
</template>

<script setup>
  import { inject } from 'vue';

  defineProps({
    isOpen: {
      type: Boolean,
      default: false,
    },
  });

  defineEmits(['close']);

  const isDarkMode = inject('isDarkMode');
  const toggleDarkMode = inject('toggleDarkMode');

  const toggleTheme = () => {
    toggleDarkMode();
  };
</script>

<style lang="scss">
@import '../styles/variables.scss';

.mobile-menu {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: white;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  transform: translateX(-100%);
  transition: transform 0.3s ease-in-out;
  overflow-y: auto;

  &.is-open {
    transform: translateX(0);
  }

  &.dark-mode {
    background-color: #1E2529;
  }

  .mobile-menu-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem;
    border-bottom: 1px solid #eee;

    .dark-mode & {
      border-bottom-color: #444;
    }

    .logo {
      display: flex;
      align-items: center;

      img {
        height: 36px;
        margin-right: 8px;
      }

      .logo-text {
        display: flex;
        flex-direction: column;

        .logo-name {
          font-family: $font-primary;
          font-size: 1.2rem;
          font-weight: 700;
          color: $primary-color-light-mode;

          .dark-mode & {
            color: $primary-color-dark-mode;
          }
        }

        .logo-description {
          font-family: $font-primary;
          font-size: 0.9rem;
          color: $secondary-color-light-mode;

          .dark-mode & {
            color: $secondary-color-dark-mode;
          }
        }
      }
    }

    .close-button {
      background: none;
      border: none;
      font-size: 1.5rem;
      color: #888;
      cursor: pointer;
      padding: 0.5rem;
      display: flex;
      align-items: center;
      justify-content: center;

      &:hover {
        color: #666;
      }

      .dark-mode & {
        color: #aaa;

        &:hover {
          color: #ccc;
        }
      }
    }
  }

  .mobile-nav {
    flex: 1;
    padding: 1rem;

    .mobile-nav-list {
      list-style: none;
      padding: 0;
      margin: 0;
    }

    .mobile-nav-item {
      margin-bottom: 1rem;

      a, button {
        display: flex;
        align-items: center;
        padding: 0.75rem;
        background: none;
        border: none;
        font-family: $font-primary;
        font-size: 1rem;
        color: $text-color-light-mode;
        text-decoration: none;
        width: 100%;
        text-align: left;
        border-radius: 4px;
        cursor: pointer;

        .dark-mode & {
          color: $text-color-dark-mode;
        }

        &:hover, &.router-link-active {
          background-color: #f5f5f5;

          .dark-mode & {
            background-color: #2a3238;
          }
        }

        &.router-link-active {
          color: $primary-color-light-mode;
          font-weight: 600;

          .dark-mode & {
            color: $primary-color-dark-mode;
          }
        }

        .icon {
          margin-right: 0.75rem;
          font-size: 20px;
        }
      }

      .theme-toggle-mobile {
        color: $text-color-light-mode;

        .dark-mode & {
          color: $text-color-dark-mode;
        }
      }

      .login-link {
        margin-top: 2rem;
        color: $primary-color-light-mode;
        font-weight: 600;

        .dark-mode & {
          color: $primary-color-dark-mode;
        }
      }
    }
  }
}

body.menu-open {
  overflow: hidden;
}
</style>
