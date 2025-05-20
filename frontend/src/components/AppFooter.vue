<template>
  <footer class="app-footer">
    <div class="footer-content">
      <div class="footer-logo">
        <img v-if="isDarkMode" alt="SignalSys Analyzer" src="../assets/logo-dark-mode.png" />
        <img v-else alt="SignalSys Analyzer" src="../assets/logo-light-mode.png" />
        <div class="logo-text">SignalSys Analyzer</div>
      </div>
      <div class="footer-links">
        <div class="footer-section">
          <h3>ACERCA DE</h3>
          <ul>
            <li><router-link to="/about/history">Nuestra historia</router-link></li>
            <li><router-link to="/about/mission-vision">Misión y Visión</router-link></li>
          </ul>
        </div>

        <div class="footer-section">
          <h3>SOPORTE</h3>
          <ul>
            <li><router-link to="/support/help">Ayuda / FAQ</router-link></li>
            <li><router-link to="/support/contact">Contacto</router-link></li>
            <li><router-link to="/support/feedback">Enviar feedback</router-link></li>
          </ul>
        </div>

        <div class="footer-section">
          <h3>LEGAL</h3>
          <ul>
            <li><router-link to="/legal/privacy">Seguridad y Privacidad de Datos</router-link></li>
            <li><router-link to="/legal/terms">Términos y Condiciones</router-link></li>
          </ul>
        </div>

        <div class="footer-section">
          <h3>CONECTA CON NOSOTROS</h3>
          <div class="social-icons">
            <a aria-label="Instagram" href="https://www.instagram.com" rel="noopener noreferrer" target="_blank">
              <i class="icon bx bxl-instagram"></i>
            </a>
            <a aria-label="LinkedIn" href="https://linkedin.com/" rel="noopener noreferrer" target="_blank">
              <i class="icon bx bxl-linkedin"></i>
            </a>
            <a aria-label="Facebook" href="https://facebook.com/" rel="noopener noreferrer" target="_blank">
              <i class="icon bx bxl-facebook"></i>
            </a>
          </div>
        </div>
      </div>
    </div>
    <div class="footer-background"></div>
  </footer>
</template>

<script setup>
  import { inject } from 'vue';

  const isDarkMode = inject('isDarkMode');
</script>

<style lang="scss">
@import '../styles/variables.scss';

.app-footer {
  position: relative;

  .dark-mode & {
    background-color: $background-color-dark-mode;
  }

  .footer-content {
    background-color: white;
    border: 1px solid #eaeaea;
    padding: 4rem;
    margin: auto;
    width: 93%;
    border-radius: 20px;
    position: relative;
    z-index: 10;

    .dark-mode & {
      background-color: $background-color-dark-mode;
      border: 1px solid #000000;
    }
  }

  .footer-background {
    position: absolute;
    background-color: $primary-color-light-mode;
    top: 25%;
    height: 92%;
    width: 100%;
    z-index: 1;
    .dark-mode & {
      background-color: $secondary-color-dark-mode;
    }
  }

  .footer-logo {
    display: flex;
    align-items: center;
    margin-bottom: 3rem;

    img {
      height: 36px;
      margin-right: 10px;
    }

    .logo-text {
      font-family: $font-primary;
      font-size: 1.2rem;
      font-weight: 700;
      color: $primary-color-light-mode;

      .dark-mode & {
        color: $primary-color-dark-mode;
      }
    }
  }

  .footer-links {
    display: flex;
    justify-content: space-between;
    flex-wrap: wrap;

    .footer-section {
      flex: 1;
      min-width: 200px;

      h3 {
        font-family: $font-primary;
        font-size: 0.9rem;
        font-weight: 700;
        margin-bottom: 1rem;
        color: $primary-color-light-mode;
        margin-top: 20px;

        .dark-mode & {
          color: $primary-color-dark-mode;
        }
      }

      ul {
        list-style: none;
        padding: 0;
        margin: 0;

        li {
          margin-bottom: 0.5rem;

          a {
            text-decoration: none;
            color: $secondary-color-light-mode;
            font-size: 0.85rem;
            transition: color 0.2s;
            &:hover {
              color: $primary-color-light-mode;
            }

            .dark-mode & {
              color: #fff;

              &:hover {
                color: $primary-color-dark-mode;
              }
            }
          }
        }
      }

      .social-icons {
        display: flex;
        gap: 1rem;

        .icon {
          color: $secondary-color-light-mode;
          font-size: 2.3rem;
          font-weight: 500;
          &:hover {
            color: $primary-color-light-mode;
          }
          .dark-mode & {
            color: #fff;
            &:hover {
              color: $primary-color-dark-mode;
            }
          }
        }
      }
    }
  }
}
@media (max-width: 576px) {
  .app-footer {
    .footer-content {
      padding: 3rem;
    }
    .footer-background {
      height: 85%;
    }
  }
}
</style>
