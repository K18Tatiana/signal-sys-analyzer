<template>
  <div class="register-page">
    <div class="register-container">
      <div class="register-card">
        <div class="register-header">
          <h1 class="register-title">Crear Cuenta</h1>
          <p class="register-subtitle">Regístrate para analizar sistemas dinámicos</p>
        </div>

        <div v-if="errorMessage" class="alert alert-error">
          <i class="bx bx-error-circle alert-icon"></i>
          <span>{{ errorMessage }}</span>
        </div>

        <div v-if="successMessage" class="alert alert-success">
          <i class="bx bx-check-circle alert-icon"></i>
          <span>{{ successMessage }}</span>
        </div>

        <form class="register-form" @submit.prevent="handleRegister">
          <div class="form-group">
            <label for="username">Nombre de usuario <span class="required">*</span></label>
            <div class="input-container">
              <i class="bx bx-user input-icon"></i>
              <input
                id="username"
                v-model="formData.name"
                autocomplete="username"
                placeholder="Ingresa tu nombre de usuario"
                required
                type="text"
              >
            </div>
            <span v-if="errors.name" class="error-message">{{ errors.name }}</span>
          </div>

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
            <label for="password">Contraseña <span class="required">*</span></label>
            <div class="input-container">
              <i class="bx bx-lock-alt input-icon"></i>
              <input
                id="password"
                v-model="formData.password"
                autocomplete="new-password"
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
            <div v-if="formData.password" class="password-strength">
              <div class="password-strength-meter">
                <div
                  class="password-strength-bar"
                  :class="passwordStrength.level"
                  :style="{ width: passwordStrength.percentage + '%' }"
                ></div>
              </div>
              <span class="password-strength-text" :class="passwordStrength.level">
                {{ passwordStrengthText }}
              </span>
            </div>
            <ul v-if="formData.password" class="password-requirements">
              <li :class="{ valid: hasMinLength }">Mínimo 8 caracteres</li>
              <li :class="{ valid: hasLowercase }">Al menos una letra minúscula</li>
              <li :class="{ valid: hasUppercase }">Al menos una letra mayúscula</li>
              <li :class="{ valid: hasNumber }">Al menos un número</li>
            </ul>
          </div>

          <div class="form-group">
            <label for="passwordConfirm">Confirmar contraseña <span class="required">*</span></label>
            <div class="input-container">
              <i class="bx bx-lock-alt input-icon"></i>
              <input
                id="passwordConfirm"
                v-model="formData.passwordConfirm"
                autocomplete="new-password"
                placeholder="Confirma tu contraseña"
                required
                :type="showConfirmPassword ? 'text' : 'password'"
              >
              <button
                :aria-label="showConfirmPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'"
                class="toggle-password"
                type="button"
                @click="showConfirmPassword = !showConfirmPassword"
              >
                <i class="bx" :class="showConfirmPassword ? 'bx-hide' : 'bx-show'"></i>
              </button>
            </div>
            <span v-if="errors.passwordConfirm" class="error-message">{{ errors.passwordConfirm }}</span>
          </div>

          <div class="form-group terms-group">
            <label class="checkbox-container">
              <input id="terms" v-model="formData.termsAccepted" required type="checkbox">
              <span class="checkbox-custom"></span>
              <span class="checkbox-label">
                Acepto los
                <router-link target="_blank" to="/legal/terms">términos y condiciones</router-link>,
                y la
                <router-link target="_blank" to="/legal/privacy">política de privacidad</router-link>
                <span class="required">*</span>
              </span>
            </label>
            <span v-if="errors.termsAccepted" class="error-message">{{ errors.termsAccepted }}</span>
          </div>

          <button class="submit-button" :disabled="isLoading || !formData.termsAccepted" type="submit">
            <span v-if="isLoading" class="loading-spinner"></span>
            <span v-else>Registrarse</span>
          </button>
        </form>

        <div class="auth-redirect">
          <p>¿Ya tienes una cuenta? <router-link to="/login">Inicia sesión</router-link></p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { computed, reactive, ref, watch } from 'vue';
  import { useRouter } from 'vue-router';
  import { useAuth } from '@/stores/auth.js';

  const router = useRouter();
  const { register } = useAuth();

  const formData = reactive({
    name: '',
    email: '',
    password: '',
    passwordConfirm: '',
    termsAccepted: false,
  });

  const errors = reactive({
    name: '',
    email: '',
    password: '',
    passwordConfirm: '',
    termsAccepted: '',
  });

  const isLoading = ref(false);
  const showPassword = ref(false);
  const showConfirmPassword = ref(false);
  const errorMessage = ref('');
  const successMessage = ref('');

  const hasMinLength = computed(() => formData.password.length >= 8);
  const hasLowercase = computed(() => /[a-z]/.test(formData.password));
  const hasUppercase = computed(() => /[A-Z]/.test(formData.password));
  const hasNumber = computed(() => /[0-9]/.test(formData.password));

  const passwordStrength = computed(() => {
    if (!formData.password) return { level: 'none', percentage: 0 };

    let strength = 0;
    if (hasMinLength.value) strength += 25;
    if (hasLowercase.value) strength += 25;
    if (hasUppercase.value) strength += 25;
    if (hasNumber.value) strength += 25;

    let level = 'weak';
    if (strength >= 100) level = 'strong';
    else if (strength >= 75) level = 'good';
    else if (strength >= 50) level = 'medium';

    return {
      level,
      percentage: strength,
    };
  });

  const passwordStrengthText = computed(() => {
    const strength = passwordStrength.value.level;
    switch (strength) {
      case 'weak': return 'Débil';
      case 'medium': return 'Media';
      case 'good': return 'Buena';
      case 'strong': return 'Fuerte';
      default: return '';
    }
  });

  watch(formData, () => {
    Object.keys(errors).forEach(key => {
      errors[key] = '';
    });
    errorMessage.value = '';
    successMessage.value = '';
  }, { deep: true });

  const handleRegister = async () => {
    let isValid = true;

    if (!formData.name.trim()) {
      errors.name = 'El nombre de usuario es obligatorio';
      isValid = false;
    }

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
    } else if (passwordStrength.value.level === 'weak') {
      errors.password = 'La contraseña es demasiado débil';
      isValid = false;
    }

    if (!formData.passwordConfirm) {
      errors.passwordConfirm = 'Debes confirmar tu contraseña';
      isValid = false;
    } else if (formData.password !== formData.passwordConfirm) {
      errors.passwordConfirm = 'Las contraseñas no coinciden';
      isValid = false;
    }

    if (!formData.termsAccepted) {
      errors.termsAccepted = 'Debes aceptar los términos y condiciones';
      isValid = false;
    }

    if (!isValid) return;

    try {
      isLoading.value = true;

      await register({
        name: formData.name,
        email: formData.email,
        password: formData.password,
      });

      successMessage.value = '¡Registro exitoso! Serás redirigido al inicio de sesión...';

      Object.keys(formData).forEach(key => {
        if (key !== 'termsAccepted') formData[key] = '';
      });

      setTimeout(() => {
        router.push('/login?registered=true');
      }, 2000);
    } catch (error) {
      console.error('Error al registrar:', error);
      errorMessage.value = error.message || 'Error al crear la cuenta. Por favor, inténtalo de nuevo.';

      if (error.message?.includes('correo') || error.message?.includes('email')) {
        errors.email = 'Este correo electrónico ya está registrado';
      }

      if (error.message?.includes('usuario') || error.message?.includes('username')) {
        errors.name = 'Este nombre de usuario ya está registrado';
      }
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

.register-page {
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

.register-container {
  width: 100%;
  max-width: 500px;
}

.register-card {
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

.register-header {
  text-align: center;
  margin-bottom: 2rem;

  .register-title {
    font-family: $font-primary;
    font-size: 2rem;
    font-weight: 700;
    color: $secondary-color-light-mode;
    margin-bottom: 0.5rem;

    .dark-mode & {
      color: $secondary-color-dark-mode;
    }
  }

  .register-subtitle {
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

.register-form {
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

    &.terms-group {
      margin-top: 1.5rem;
    }
  }

  .checkbox-container {
    display: flex;
    align-items: flex-start;
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
      flex-shrink: 0;
      height: 18px;
      width: 18px;
      margin-top: 4px;
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

    .checkbox-label {
      margin-left: 0.5rem;
      font-size: 0.9rem;
      line-height: 1.5;
      color: $text-color-light-mode;

      a {
        color: $primary-color-light-mode;
        text-decoration: none;

        &:hover {
          text-decoration: underline;
        }
      }

      .dark-mode & {
        color: $text-color-dark-mode;

        a {
          color: $primary-color-dark-mode;
        }
      }
    }
  }

  .password-strength {
    margin-top: 0.75rem;

    .password-strength-meter {
      width: 100%;
      height: 5px;
      background-color: rgba(0, 0, 0, 0.1);
      border-radius: 3px;
      overflow: hidden;

      .dark-mode & {
        background-color: rgba(255, 255, 255, 0.1);
      }

      .password-strength-bar {
        height: 100%;
        border-radius: 3px;
        transition: width 0.3s ease;

        &.weak {
          background-color: #e74c3c;
        }

        &.medium {
          background-color: #f39c12;
        }

        &.good {
          background-color: #3498db;
        }

        &.strong {
          background-color: #2ecc71;
        }
      }
    }

    .password-strength-text {
      display: block;
      margin-top: 0.25rem;
      font-size: 0.8rem;
      text-align: right;

      &.weak {
        color: #e74c3c;
      }

      &.medium {
        color: #f39c12;
      }

      &.good {
        color: #3498db;
      }

      &.strong {
        color: #2ecc71;
      }
    }
  }

  .password-requirements {
    margin-top: 0.5rem;
    padding-left: 1.2rem;
    font-size: 0.8rem;
    color: rgba($text-color-light-mode, 0.7);

    .dark-mode & {
      color: rgba($text-color-dark-mode, 0.7);
    }

    li {
      margin-bottom: 0.25rem;

      &.valid {
        color: #2ecc71;

        &::marker {
          content: "✓ ";
        }
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
    margin-top: 2rem;

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
