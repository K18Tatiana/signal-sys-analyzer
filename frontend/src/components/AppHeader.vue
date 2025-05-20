<template>
  <header class="app-header">
    <div class="logo">
      <router-link to="/">
        <img v-if="isDarkMode" alt="SignalSys Analyzer" src="../assets/logo-dark-mode.png" />
        <img v-else alt="SignalSys Analyzer" src="../assets/logo-light-mode.png" />
        <div class="logo-text">
          <span class="logo-name">SignalSys</span>
          <span class="logo-description">Analyzer</span>
        </div>
      </router-link>
    </div>

    <!-- Botón para el menú móvil -->
    <button class="mobile-menu-toggle" @click="openMobileMenu">
      <span class="menu-icon"></span>
    </button>

    <!-- Navegación para escritorio -->
    <nav class="desktop-nav">
      <div class="nav-links">
        <router-link to="/">
          <i class='icon bx bx-home-alt'></i> Inicio
        </router-link>
        <router-link to="/guide">
          <i class='icon bx bx-book-open'></i> Guía / Tutorial
        </router-link>
        <router-link to="/language">
          <i class='icon bx bx-globe'></i> Idioma
        </router-link>
        <ThemeToggle />
        <button class="btn btn-secondary">
          <router-link to="/login">
            <i class='icon bx bx-user'></i> Iniciar sesión
          </router-link>
        </button>
        <button class="btn btn-primary">
          <router-link to="/register">
            <i class='icon bx bx-user'></i> Registrarse
          </router-link>
        </button>
      </div>
    </nav>

    <!-- Menú móvil -->
    <MobileMenu
      :is-open="isMobileMenuOpen"
      @close="closeMobileMenu"
    />

    <div
      v-if="isMobileMenuOpen"
      class="mobile-menu-overlay"
      @click="closeMobileMenu"
    ></div>
  </header>
</template>

<script setup>
  import { inject, onMounted, onUnmounted, ref } from 'vue';
  import ThemeToggle from './ThemeToggle.vue';
  import MobileMenu from './MobileMenu.vue';

  const isDarkMode = inject('isDarkMode');
  const isMobileMenuOpen = ref(false);

  const openMobileMenu = () => {
    isMobileMenuOpen.value = true;
    document.body.classList.add('menu-open');
  };
  const closeMobileMenu = () => {
    isMobileMenuOpen.value = false;
    document.body.classList.remove('menu-open');
  };

  const handleResize = () => {
    if (window.innerWidth >= 992 && isMobileMenuOpen.value) {
      closeMobileMenu();
    }
  };

  onMounted(() => {
    window.addEventListener('resize', handleResize);
  });

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize);
  });
</script>

<style lang="scss">
@import '../styles/variables.scss';

.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background-color: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: relative;

  .dark-mode & {
    background-color: $background-color-dark-mode;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  }

  .logo {
    display: flex;
    align-items: center;

    a {
      display: flex;
      align-items: center;
      text-decoration: none;
    }

    img {
      height: 50px;
      margin-right: 8px;
    }

    .logo-text {
      display: flex;
      flex-direction: column;
      font-size: 1.2rem;
      font-family: $font-primary;
      color: $secondary-color-light-mode;

      .dark-mode & {
        color: $secondary-color-dark-mode;
      }
    }
  }

  // Botón del menú móvil
  .mobile-menu-toggle {
    display: none;
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.5rem;

    .menu-icon {
      display: block;
      width: 24px;
      height: 2px;
      background-color: $text-color-light-mode;
      position: relative;

      .dark-mode & {
        background-color: $text-color-dark-mode;
      }

      &::before,
      &::after {
        content: '';
        display: block;
        width: 24px;
        height: 2px;
        background-color: $text-color-light-mode;
        position: absolute;

        .dark-mode & {
          background-color: $text-color-dark-mode;
        }
      }

      &::before {
        top: -6px;
      }

      &::after {
        bottom: -6px;
      }
    }
  }

  // Navegación de escritorio
  .desktop-nav {
    .nav-links {
      display: flex;
      align-items: center;
      gap: 1.5rem;

      a {
        display: flex;
        align-items: center;
        text-decoration: none;
        color: $text-color-light-mode;
        font-size: 0.9rem;

        .dark-mode & {
          color: $text-color-dark-mode;
        }

        .icon {
          margin-right: 4px;
          font-size: 20px;
        }

        &.router-link-active {
          color: $primary-color-light-mode;
          font-weight: 600;

          .dark-mode & {
            color: $primary-color-dark-mode;
          }
        }
      }

      .btn {
        padding: 0.5rem 1rem;
        border-radius: 4px;
        font-size: 0.9rem;
        cursor: pointer;
        border: none;
        display: flex;
        align-items: center;

        &.btn-primary {
          background-color: $primary-color-light-mode;

          a {
            color: white !important;

            .dark-mode & {
              color: #000 !important;
            }
          }

          .dark-mode & {
            background-color: $primary-color-dark-mode;
          }
        }

        &.btn-secondary {
          background-color: transparent;
          color: $text-color-light-mode;

          .dark-mode & {
            color: $text-color-dark-mode;
          }
        }

        .icon {
          margin-right: 4px;
        }
      }
    }
  }

  .mobile-menu-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 999;
  }

  @media (max-width: 991px) {
    .desktop-nav {
      display: none;
    }

    .mobile-menu-toggle {
      display: block;
    }
  }

  @media (max-width: 576px) {
    padding: 0.75rem 1rem;

    .logo {
      img {
        height: 30px;
      }

      .logo-text {
        font-size: 1rem;
      }
    }
  }
}

</style>
