<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card">
        <div class="login-header">
          <h1 class="login-title">Iniciar Sesión</h1>
          <p class="login-subtitle">Accede a tu cuenta para gestionar tus análisis</p>
        </div>

        <div v-if="registeredSuccessfully" class="alert alert-success">
          <i class="bx bx-check-circle alert-icon"></i>
          <span>Registro exitoso. Ya puedes iniciar sesión con tus credenciales.</span>
        </div>

        <div v-if="errorMessage" class="alert alert-error">
          <i class="bx bx-error-circle alert-icon"></i>
          <span>{{ errorMessage }}</span>
        </div>

        <form class="login-form" @submit.prevent="handleLogin">
          <div class="form-group">
            <label for="email">Correo electrónico <span class="required">*</span></label>
            <div class="input-container">
              <i class="bx bx-envelope input-icon"></i>
              <input
                id="email"
                v-model="formData.email"
                autocomplete="email"
                placeholder="tucorreo@ejemplo.com"
                required
                type="email"
              >
            </div>
            <span v-if="errors.email" class="error-message">{{ errors.email }}</span>
          </div>

          <div class="form-group">
            <div class="label-row">
              <label for="password">Contraseña <span class="required">*</span></label>
              <router-link class="forgot-password" to="/forgot-password">
                ¿Olvidaste tu contraseña?
              </router-link>
            </div>
            <div class="input-container">
              <i class="bx bx-lock-alt input-icon"></i>
              <input
                id="password"
                v-model="formData.password"
                autocomplete="current-password"
                placeholder="Ingresa tu contraseña"
                required
                :type="showPassword ? 'text' : 'password'"
              >
              <button
                :aria-label="showPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'"
                class="toggle-password"
                type="button"
                @click="showPassword = !showPassword"
              >
                <i class="bx" :class="showPassword ? 'bx-hide' : 'bx-show'"></i>
              </button>
            </div>
            <span v-if="errors.password" class="error-message">{{ errors.password }}</span>
          </div>

          <div class="form-group remember-row">
            <label class="checkbox-container">
              <input id="remember" v-model="formData.remember" type="checkbox">
              <span class="checkbox-custom"></span>
              <span class="checkbox-text">Recordarme en este dispositivo</span>
            </label>
          </div>

          <button class="submit-button" :disabled="isLoading" type="submit">
            <span v-if="isLoading" class="loading-spinner"></span>
            <span v-else>Iniciar sesión</span>
          </button>
        </form>

        <div class="auth-redirect">
          <p>¿No tienes una cuenta? <router-link to="/register">Regístrate</router-link></p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { onMounted, reactive, ref } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { useAuth } from '@/stores/auth.js';

  const router = useRouter();
  const route = useRoute();
  const { login } = useAuth();

  const formData = reactive({
    email: '',
    password: '',
    remember: false,
  });

  const errors = reactive({
    email: '',
    password: '',
  });

  const isLoading = ref(false);
  const showPassword = ref(false);
  const errorMessage = ref('');
  const registeredSuccessfully = ref(false);

  onMounted(() => {
    if (route.query.registered === 'true') {
      registeredSuccessfully.value = true;
    }
  });

  const handleLogin = async () => {
    errors.email = '';
    errors.password = '';
    errorMessage.value = '';
    registeredSuccessfully.value = false;

    let isValid = true;

    if (!formData.email) {
      errors.email = 'El correo electrónico es obligatorio';
      isValid = false;
    } else if (!isValidEmail(formData.email)) {
      errors.email = 'Correo electrónico no válido';
      isValid = false;
    }

    if (!formData.password) {
      errors.password = 'La contraseña es obligatoria';
      isValid = false;
    }

    if (!isValid) return;

    try {
      isLoading.value = true;

      await login({
        email: formData.email,
        password: formData.password,
      });

      if (formData.remember) {
        localStorage.setItem('remember_me', 'true');
      } else {
        localStorage.removeItem('remember_me');
      }

      const redirectPath = route.query.redirect || '/dashboard';
      router.push(redirectPath);
    } catch (error) {
      console.error('Error al iniciar sesión:', error);
      errorMessage.value = error.message || 'Error al iniciar sesión. Por favor, inténtalo de nuevo.';
    } finally {
      isLoading.value = false;
    }
  };

  const isValidEmail = email => {
    const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(String(email).toLowerCase());
  };
</script>

<style lang="scss" scoped>
@import '/src/styles/variables.scss';

.login-page {
  min-height: calc(100vh - 70px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem 1rem;
  background-color: #f8f9fa;

  .dark-mode & {
    background-color: darken($background-color-dark-mode, 2%);
  }
}

.login-container {
  width: 100%;
  max-width: 450px;
}

.login-card {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
  padding: 2.5rem;

  .dark-mode & {
    background-color: $background-color-dark-mode;
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.2);
  }

  @media (max-width: 576px) {
    padding: 2rem 1.5rem;
  }
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;

  .login-title {
    font-family: $font-primary;
    font-size: 2rem;
    font-weight: 700;
    color: $secondary-color-light-mode;
    margin-bottom: 0.5rem;

    .dark-mode & {
      color: $secondary-color-dark-mode;
    }
  }

  .login-subtitle {
    font-size: 1rem;
    color: rgba($text-color-light-mode, 0.7);

    .dark-mode & {
      color: rgba($text-color-dark-mode, 0.7);
    }
  }
}

// Alertas
.alert {
  margin-bottom: 1.5rem;
  padding: 1rem;
  border-radius: 8px;
  display: flex;
  align-items: center;

  &.alert-error {
    background-color: rgba(#e74c3c, 0.1);
    border-left: 4px solid #e74c3c;
    color: #e74c3c;
  }

  &.alert-success {
    background-color: rgba(#2ecc71, 0.1);
    border-left: 4px solid #2ecc71;
    color: #2ecc71;
  }

  .alert-icon {
    margin-right: 0.75rem;
    font-size: 1.2rem;
  }
}

.login-form {
  .form-group {
    margin-bottom: 1.5rem;

    label {
      display: block;
      margin-bottom: 0.5rem;
      font-weight: 500;
      font-size: 0.95rem;
      color: $text-color-light-mode;

      .dark-mode & {
        color: $text-color-dark-mode;
      }
    }

    .label-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.5rem;

      .forgot-password {
        font-size: 0.85rem;
        color: $primary-color-light-mode;
        text-decoration: none;

        &:hover {
          text-decoration: underline;
        }

        .dark-mode & {
          color: $primary-color-dark-mode;
        }
      }
    }

    .input-container {
      position: relative;

      .input-icon {
        position: absolute;
        left: 1rem;
        top: 50%;
        transform: translateY(-50%);
        font-size: 1.2rem;
        color: rgba($text-color-light-mode, 0.5);

        .dark-mode & {
          color: rgba($text-color-dark-mode, 0.5);
        }
      }

      input {
        width: 100%;
        padding: 0.8rem 1rem 0.8rem 2.8rem;
        border: 1px solid rgba(0, 0, 0, 0.1);
        border-radius: 8px;
        font-size: 1rem;
        background-color: white;
        transition: border-color 0.2s, box-shadow 0.2s;

        &:focus {
          outline: none;
          border-color: $primary-color-light-mode;
          box-shadow: 0 0 0 3px rgba($primary-color-light-mode, 0.2);
        }

        &::placeholder {
          color: rgba($text-color-light-mode, 0.4);
        }

        .dark-mode & {
          background-color: lighten($background-color-dark-mode, 5%);
          border-color: rgba(255, 255, 255, 0.1);
          color: $text-color-dark-mode;

          &:focus {
            border-color: $primary-color-dark-mode;
            box-shadow: 0 0 0 3px rgba($primary-color-dark-mode, 0.2);
          }

          &::placeholder {
            color: rgba($text-color-dark-mode, 0.4);
          }
        }
      }

      .toggle-password {
        position: absolute;
        right: 1rem;
        top: 50%;
        transform: translateY(-50%);
        background: none;
        border: none;
        cursor: pointer;
        color: rgba($text-color-light-mode, 0.6);
        font-size: 1.2rem;
        padding: 0;

        .dark-mode & {
          color: rgba($text-color-dark-mode, 0.6);
        }

        &:hover {
          color: $primary-color-light-mode;

          .dark-mode & {
            color: $primary-color-dark-mode;
          }
        }
      }
    }

    .error-message {
      display: block;
      margin-top: 0.5rem;
      color: #e74c3c;
      font-size: 0.85rem;
    }

    &.remember-row {
      display: flex;
      align-items: center;
    }
  }

  .checkbox-container {
    display: flex;
    align-items: center;
    cursor: pointer;

    input[type="checkbox"] {
      position: absolute;
      opacity: 0;
      cursor: pointer;
      height: 0;
      width: 0;

      &:checked ~ .checkbox-custom {
        background-color: $primary-color-light-mode;
        border-color: $primary-color-light-mode;

        &:after {
          opacity: 1;
        }

        .dark-mode & {
          background-color: $primary-color-dark-mode;
          border-color: $primary-color-dark-mode;
        }
      }
    }

    .checkbox-custom {
      position: relative;
      height: 18px;
      width: 18px;
      background-color: white;
      border: 2px solid rgba(0, 0, 0, 0.2);
      border-radius: 4px;
      transition: all 0.2s;

      .dark-mode & {
        background-color: lighten($background-color-dark-mode, 5%);
        border-color: rgba(255, 255, 255, 0.2);
      }

      &:after {
        content: "";
        position: absolute;
        left: 5px;
        top: 1px;
        width: 5px;
        height: 10px;
        border: solid white;
        border-width: 0 2px 2px 0;
        transform: rotate(45deg);
        opacity: 0;
        transition: opacity 0.2s;
      }
    }

    .checkbox-text {
      margin-left: 0.5rem;
      font-size: 0.9rem;
      color: $text-color-light-mode;

      .dark-mode & {
        color: $text-color-dark-mode;
      }
    }
  }

  .submit-button {
    width: 100%;
    padding: 0.9rem;
    background-color: $primary-color-light-mode;
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    justify-content: center;
    align-items: center;

    &:hover {
      background-color: darken($primary-color-light-mode, 5%);
    }

    &:disabled {
      opacity: 0.7;
      cursor: not-allowed;
    }

    .dark-mode & {
      background-color: $primary-color-dark-mode;
      color: black;

      &:hover {
        background-color: lighten($primary-color-dark-mode, 5%);
      }
    }
  }
}

.loading-spinner {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 0.8s linear infinite;

  .dark-mode & {
    border: 2px solid rgba(0, 0, 0, 0.3);
    border-top-color: black;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.required {
  color: #e74c3c;
  margin-left: 2px;
}

.auth-redirect {
  margin-top: 2rem;
  text-align: center;
  font-size: 0.95rem;
  color: $text-color-light-mode;

  .dark-mode & {
    color: $text-color-dark-mode;
  }

  a {
    color: $primary-color-light-mode;
    text-decoration: none;
    font-weight: 600;

    &:hover {
      text-decoration: underline;
    }

    .dark-mode & {
      color: $primary-color-dark-mode;
    }
  }
}
</style>
