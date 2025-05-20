<template>
  <div class="not-found-page">
    <div class="not-found-content">
      <div class="error-code">404</div>
      <h1 class="error-title">Página no encontrada</h1>
      <p class="error-description">
        La página que está buscando no existe o ha sido movida.
      </p>
      <router-link class="home-button" to="/">
        <i class="icon home-icon"></i> Volver a la página principal
      </router-link>
    </div>
  </div>
</template>

<style lang="scss">
@import '../styles/variables.scss';

.not-found-page {
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
  padding: 4rem 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;

  .not-found-content {
    background-color: white;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    padding: 3rem 2rem;

    .dark-mode & {
      background-color: #2a3238;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
    }
  }

  .error-code {
    font-size: 5rem;
    font-weight: 700;
    color: $primary-color-light-mode;
    line-height: 1.1;

    .dark-mode & {
      color: $primary-color-dark-mode;
    }
  }

  .error-title {
    font-size: 1.8rem;
    margin-bottom: 1rem;
    color: $secondary-color-light-mode;

    .dark-mode & {
      color: $secondary-color-dark-mode;
    }
  }

  .error-description {
    font-size: 1.1rem;
    margin-bottom: 2rem;
    color: #666;

    .dark-mode & {
      color: #aaa;
    }
  }

  .home-button {
    display: inline-flex;
    align-items: center;
    padding: 0.75rem 1.5rem;
    background-color: $primary-color-light-mode;
    color: white;
    text-decoration: none;
    border-radius: 4px;
    font-weight: 600;
    transition: background-color 0.3s;

    .dark-mode & {
      background-color: $primary-color-dark-mode;
      color: $background-color-dark-mode;
    }

    &:hover {
      background-color: darken($primary-color-light-mode, 10%);
      text-decoration: none;

      .dark-mode & {
        background-color: darken($primary-color-dark-mode, 10%);
      }
    }

    .home-icon {
      margin-right: 0.5rem;

      &::before {
        content: "🏠";
      }
    }
  }
}
</style>
