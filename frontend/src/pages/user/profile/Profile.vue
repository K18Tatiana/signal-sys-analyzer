<template>
  <div class="profile-page">
    <div class="profile-container">
      <!-- Header del perfil -->
      <div class="profile-header">
        <div class="user-avatar">
          <span class="avatar-text">{{ userInitials }}</span>
        </div>
        <div class="header-info">
          <h1>Mi Perfil</h1>
          <p>Gestiona tu información personal y configuración de cuenta</p>
        </div>
      </div>

      <!-- Información del usuario (solo lectura) -->
      <div class="profile-card">
        <h2 class="card-title">Información de la cuenta</h2>
        <div class="info-grid">
          <div class="info-item">
            <label>Correo electrónico</label>
            <div class="info-value">
              <span>{{ user?.email }}</span>
              <span class="readonly-badge">Solo lectura</span>
            </div>
          </div>
          <div class="info-item">
            <label>Miembro desde</label>
            <div class="info-value">
              <span>{{ formatDate(user?.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Formulario de actualización -->
      <div class="profile-card">
        <h2 class="card-title">Actualizar información</h2>

        <!-- Mostrar mensajes de éxito o error -->
        <div v-if="successMessage" class="message success-message">
          <i class="bx bx-check-circle"></i>
          <span>{{ successMessage }}</span>
        </div>

        <div v-if="errorMessage" class="message error-message">
          <i class="bx bx-error-circle"></i>
          <span>{{ errorMessage }}</span>
        </div>

        <form class="profile-form" @submit.prevent="updateProfile">
          <!-- Sección del nombre de usuario -->
          <div class="form-section">
            <h3 class="section-title">Cambiar nombre de usuario</h3>
            <div class="current-value">
              <span class="label">Nombre actual:</span>
              <span class="value">{{ user?.username || 'No disponible' }}</span>
            </div>
            <div class="form-group">
              <label for="username">Nuevo nombre de usuario (opcional)</label>
              <input
                id="username"
                v-model="formData.username"
                class="form-input"
                :disabled="loading"
                placeholder="Deja vacío si no quieres cambiar el nombre"
                type="text"
              />
            </div>
          </div>

          <!-- Sección de contraseña -->
          <div class="form-section">
            <h3 class="section-title">Cambiar contraseña</h3>
            <div class="form-group">
              <label for="newPassword">Nueva contraseña (opcional)</label>
              <div class="password-input-container">
                <input
                  id="newPassword"
                  v-model="formData.new_password"
                  class="form-input"
                  :disabled="loading"
                  placeholder="Deja vacío si no quieres cambiar la contraseña"
                  :type="showNewPassword ? 'text' : 'password'"
                />
                <button
                  v-if="formData.new_password"
                  class="password-toggle"
                  :disabled="loading"
                  type="button"
                  @click="showNewPassword = !showNewPassword"
                >
                  <i :class="showNewPassword ? 'bx bx-hide' : 'bx bx-show'"></i>
                </button>
              </div>
            </div>
          </div>

          <!-- Confirmación con contraseña actual (solo si hay cambios) -->
          <div v-if="hasChanges" class="form-section confirmation-section">
            <h3 class="section-title">Confirmar cambios</h3>
            <div class="form-group">
              <label for="oldPassword">Contraseña actual *</label>
              <div class="password-input-container">
                <input
                  id="oldPassword"
                  v-model="formData.old_password"
                  class="form-input"
                  :disabled="loading"
                  placeholder="Ingresa tu contraseña actual para confirmar los cambios"
                  required
                  :type="showOldPassword ? 'text' : 'password'"
                />
                <button
                  v-if="formData.old_password"
                  class="password-toggle"
                  :disabled="loading"
                  type="button"
                  @click="showOldPassword = !showOldPassword"
                >
                  <i :class="showOldPassword ? 'bx bx-hide' : 'bx bx-show'"></i>
                </button>
              </div>
            </div>
          </div>

          <div class="form-actions">
            <button
              class="btn btn-secondary"
              :disabled="loading"
              type="button"
              @click="resetForm"
            >
              Cancelar
            </button>
            <button
              class="btn btn-primary"
              :disabled="loading || !canSubmit"
              type="submit"
            >
              <i v-if="loading" class="bx bx-loader-alt rotating"></i>
              <i v-else class="bx bx-save"></i>
              <span>{{ loading ? 'Actualizando...' : 'Guardar cambios' }}</span>
            </button>
          </div>
        </form>
      </div>

      <!-- Zona de peligro -->
      <div class="profile-card danger-zone">
        <h2 class="card-title danger-title">Zona de peligro</h2>
        <p class="danger-description">
          Las siguientes acciones son irreversibles. Por favor, procede con precaución.
        </p>
        <button
          class="btn btn-danger"
          :disabled="loading"
          @click="showDeleteConfirmation = true"
        >
          <i class="bx bx-trash"></i>
          <span>Eliminar cuenta</span>
        </button>
      </div>
    </div>

    <!-- Modal de confirmación para eliminar cuenta -->
    <div v-if="showDeleteConfirmation" class="modal-overlay" @click="showDeleteConfirmation = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Confirmar eliminación de cuenta</h3>
          <button class="modal-close" @click="showDeleteConfirmation = false">
            <i class="bx bx-x"></i>
          </button>
        </div>
        <div class="modal-body">
          <p>¿Estás seguro de que quieres eliminar tu cuenta permanentemente?</p>
          <p class="warning-text">Esta acción no se puede deshacer y perderás todos tus documentos y análisis.</p>

          <div class="form-group">
            <label for="deletePassword">Ingresa tu contraseña para confirmar</label>
            <input
              id="deletePassword"
              v-model="deletePassword"
              class="form-input"
              :disabled="loading"
              placeholder="Tu contraseña actual"
              type="password"
            />
          </div>
        </div>
        <div class="modal-actions">
          <button
            class="btn btn-secondary"
            :disabled="loading"
            @click="showDeleteConfirmation = false"
          >
            Cancelar
          </button>
          <button
            class="btn btn-danger"
            :disabled="loading || !deletePassword"
            @click="deleteAccount"
          >
            <i v-if="loading" class="bx bx-loader-alt rotating"></i>
            <i v-else class="bx bx-trash"></i>
            <span>{{ loading ? 'Eliminando...' : 'Eliminar cuenta' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { computed, onMounted, reactive, ref } from 'vue';
  import { useRouter } from 'vue-router';
  import { useAuth } from '@/stores/auth.js';

  const router = useRouter();
  const { user, logout, isAuthenticated } = useAuth();

  // Configuración api
  const API_URL = import.meta.env.VITE_API_URL

  // Estados reactivos
  const loading = ref(false);
  const successMessage = ref('');
  const errorMessage = ref('');
  const showNewPassword = ref(false);
  const showOldPassword = ref(false);
  const showDeleteConfirmation = ref(false);
  const deletePassword = ref('');

  // Datos del formulario
  const formData = reactive({
    username: '',
    new_password: '',
    old_password: '',
  });

  // Función para obtener headers de autorización
  const getAuthHeaders = () => {
    const headers = { 'Content-Type': 'application/json' };
    if (isAuthenticated.value) {
      const token = localStorage.getItem('token');
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }
    return headers;
  };

  // Calcular iniciales del usuario para el avatar
  const userInitials = computed(() => {
    if (!user.value || !user.value.username) return '';
    return user.value.username.charAt(0).toUpperCase();
  });

  // Validar si hay cambios pendientes
  const hasChanges = computed(() => {
    const usernameChanged = formData.username.trim() !== '' && formData.username !== user.value?.username;
    const passwordChanged = formData.new_password.trim() !== '';
    return usernameChanged || passwordChanged;
  });

  // Validar si se puede enviar el formulario
  const canSubmit = computed(() => {
    // Si no hay cambios, no se puede enviar
    if (!hasChanges.value) return false;

    // Si hay cambios, debe tener contraseña actual
    return formData.old_password.trim() !== '';
  });

  // Cargar perfil del usuario
  const loadUserProfile = async () => {
    try {
      const response = await fetch(`${API_URL}/profile`, {
        method: 'GET',
        headers: getAuthHeaders(),
      });

      if (!response.ok) {
        throw new Error('Error al cargar el perfil');
      }

      const userData = await response.json();
      // Actualizar el user en el store manualmente si es necesario
      Object.assign(user.value, userData);

      if (user.value) {
        formData.username = user.value.username || '';
      }
    } catch (error) {
      console.error('Error cargando perfil:', error);
      errorMessage.value = 'Error al cargar el perfil del usuario';
    }
  };

  // Actualizar perfil
  const updateProfile = async () => {
    // Validar que haya al menos un cambio
    if (!hasChanges.value) {
      errorMessage.value = 'No has realizado ningún cambio. Modifica tu nombre de usuario o contraseña para continuar.';
      return;
    }

    // Validar que tenga contraseña actual
    if (!formData.old_password.trim()) {
      errorMessage.value = 'Debes ingresar tu contraseña actual para confirmar los cambios.';
      return;
    }

    loading.value = true;
    errorMessage.value = '';
    successMessage.value = '';

    try {
      const updateData = {
        old_password: formData.old_password,
      };

      // Solo incluir username si cambió
      if (formData.username.trim() && formData.username !== user.value?.username) {
        updateData.username = formData.username;
      }

      // Solo incluir nueva contraseña si se proporcionó
      if (formData.new_password.trim()) {
        updateData.new_password = formData.new_password;
      }

      const response = await fetch(`${API_URL}/user/update`, {
        method: 'PUT',
        headers: getAuthHeaders(),
        body: JSON.stringify(updateData),
      });

      if (!response.ok) {
        let errorMsg = 'Error al actualizar el perfil';

        try {
          const errorData = await response.json();
          errorMsg = errorData.message || errorData.error || errorMsg;
        } catch (parseError) {
          console.error('Error parsing error response:', parseError);
        }

        throw new Error(errorMsg);
      }

      // Actualizar el perfil
      await loadUserProfile();

      successMessage.value = 'Perfil actualizado exitosamente';

      // Limpiar campos del formulario
      formData.username = '';
      formData.new_password = '';
      formData.old_password = '';

      // Ocultar mensaje de éxito después de 5 segundos
      setTimeout(() => {
        successMessage.value = '';
      }, 5000);

    } catch (error) {
      console.error('Error actualizando perfil:', error);
      errorMessage.value = error.message || 'Error al actualizar el perfil';
    } finally {
      loading.value = false;
    }
  };

  // Resetear formulario
  const resetForm = () => {
    formData.username = '';
    formData.new_password = '';
    formData.old_password = '';
    errorMessage.value = '';
    successMessage.value = '';
  };

  // Eliminar cuenta
  const deleteAccount = async () => {
    loading.value = true;
    errorMessage.value = '';

    try {
      const response = await fetch(`${API_URL}/user/delete`, {
        method: 'DELETE',
        headers: getAuthHeaders(),
        body: JSON.stringify({ password: deletePassword.value }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Error al eliminar la cuenta');
      }

      // Cerrar sesión y redirigir
      logout();
      router.push('/login');

    } catch (error) {
      console.error('Error eliminando cuenta:', error);
      errorMessage.value = error.message || 'Error al eliminar la cuenta';
    } finally {
      loading.value = false;
      showDeleteConfirmation.value = false;
      deletePassword.value = '';
    }
  };

  // Función para formatear fechas
  const formatDate = dateString => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('es-ES', {
      day: '2-digit',
      month: 'long',
      year: 'numeric',
    });
  };

  // Al montar el componente
  onMounted(async () => {
    if (!isAuthenticated.value) {
      router.push('/login');
      return;
    }

    await loadUserProfile();
  });
</script>

<style lang="scss" scoped>
@import '/src/styles/variables.scss';

.profile-page {
  min-height: calc(100vh - 70px);
  background-color: #f8f9fa;
  padding: 2rem 1rem;

  .dark-mode & {
    background-color: darken($background-color-dark-mode, 2%);
  }
}

.profile-container {
  max-width: 800px;
  margin: 0 auto;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-bottom: 2rem;

  .user-avatar {
    width: 80px;
    height: 80px;
    background-color: $primary-color-light-mode;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;

    .dark-mode & {
      background-color: $primary-color-dark-mode;
    }

    .avatar-text {
      color: white;
      font-size: 2.2rem;
      font-weight: 600;

      .dark-mode & {
        color: black;
      }
    }
  }

  .header-info {
    h1 {
      font-size: 2rem;
      margin: 0;
      color: $secondary-color-light-mode;
      font-weight: 700;

      .dark-mode & {
        color: $secondary-color-dark-mode;
      }
    }

    p {
      margin: 0.5rem 0 0;
      color: rgba($text-color-light-mode, 0.7);
      font-size: 1.1rem;

      .dark-mode & {
        color: rgba($text-color-dark-mode, 0.7);
      }
    }
  }

  @media (max-width: 576px) {
    flex-direction: column;
    text-align: center;

    .header-info h1 {
      font-size: 1.7rem;
    }
  }
}

.profile-card {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
  padding: 2rem;
  margin-bottom: 1.5rem;

  .dark-mode & {
    background-color: $background-color-dark-mode;
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
  }

  .card-title {
    font-size: 1.3rem;
    font-weight: 600;
    color: $secondary-color-light-mode;
    margin-top: 0;
    margin-bottom: 1.5rem;

    .dark-mode & {
      color: $secondary-color-dark-mode;
    }
  }

  &.danger-zone {
    border-left: 4px solid #e74c3c;

    .danger-title {
      color: #e74c3c;
    }

    .danger-description {
      color: rgba($text-color-light-mode, 0.7);
      margin-bottom: 1.5rem;

      .dark-mode & {
        color: rgba($text-color-dark-mode, 0.7);
      }
    }
  }
}

.info-grid {
  display: grid;
  gap: 1.5rem;

  .info-item {
    label {
      display: block;
      font-weight: 500;
      color: $text-color-light-mode;
      margin-bottom: 0.5rem;

      .dark-mode & {
        color: $text-color-dark-mode;
      }
    }

    .info-value {
      display: flex;
      align-items: center;
      gap: 1rem;

      span {
        color: rgba($text-color-light-mode, 0.8);

        .dark-mode & {
          color: rgba($text-color-dark-mode, 0.8);
        }
      }

      .readonly-badge {
        background-color: rgba($primary-color-light-mode, 0.1);
        color: $primary-color-light-mode;
        padding: 0.25rem 0.5rem;
        border-radius: 4px;
        font-size: 0.8rem;
        font-weight: 500;

        .dark-mode & {
          background-color: rgba($primary-color-dark-mode, 0.2);
          color: $primary-color-dark-mode;
        }
      }
    }
  }
}

.message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
  font-weight: 500;

  &.success-message {
    background-color: rgba(#27ae60, 0.1);
    color: #27ae60;
    border: 1px solid rgba(#27ae60, 0.2);
  }

  &.error-message {
    background-color: rgba(#e74c3c, 0.1);
    color: #e74c3c;
    border: 1px solid rgba(#e74c3c, 0.2);
  }

  i {
    font-size: 1.2rem;
  }
}

.profile-form {
  .form-group {
    margin-bottom: 1.5rem;

    label {
      display: block;
      font-weight: 500;
      color: $text-color-light-mode;
      margin-bottom: 0.5rem;

      .dark-mode & {
        color: $text-color-dark-mode;
      }
    }

    .form-input {
      width: 100%;
      padding: 0.8rem 1rem;
      border: 2px solid rgba($text-color-light-mode, 0.1);
      border-radius: 8px;
      font-size: 1rem;
      transition: border-color 0.2s;
      background-color: white;
      color: $text-color-light-mode;

      &:focus {
        outline: none;
        border-color: $primary-color-light-mode;
      }

      &:disabled {
        background-color: rgba($text-color-light-mode, 0.05);
        cursor: not-allowed;
      }

      .dark-mode & {
        background-color: lighten($background-color-dark-mode, 3%);
        border-color: rgba($text-color-dark-mode, 0.1);
        color: $text-color-dark-mode;

        &:focus {
          border-color: $primary-color-dark-mode;
        }

        &:disabled {
          background-color: rgba($text-color-dark-mode, 0.05);
        }
      }
    }

    .password-input-container {
      position: relative;

      .password-toggle {
        position: absolute;
        right: 1rem;
        top: 50%;
        transform: translateY(-50%);
        background: none;
        border: none;
        color: rgba($text-color-light-mode, 0.6);
        cursor: pointer;
        font-size: 1.2rem;
        padding: 0.2rem;

        &:hover {
          color: $primary-color-light-mode;
        }

        &:disabled {
          cursor: not-allowed;
          opacity: 0.5;
        }

        .dark-mode & {
          color: rgba($text-color-dark-mode, 0.6);

          &:hover {
            color: $primary-color-dark-mode;
          }
        }
      }
    }
  }
}

.form-section {
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid rgba($text-color-light-mode, 0.1);
  &:last-child {
    border-bottom: none;
    margin-bottom: 0;
  }
  &.confirmation-section {
    background-color: rgba($primary-color-light-mode, 0.02);
    padding: 1.5rem;
    border-radius: 8px;
    border: 1px solid rgba($primary-color-light-mode, 0.1);
    .dark-mode & {
      background-color: rgba($primary-color-dark-mode, 0.05);
      border-color: rgba($primary-color-dark-mode, 0.2);
    }
  }
  .dark-mode & {
    border-bottom-color: rgba($text-color-dark-mode, 0.1);
  }
  .section-title {
    font-size: 1.1rem;
    font-weight: 600;
    color: $secondary-color-light-mode;
    margin: 0 0 1rem 0;
    .dark-mode & {
      color: $secondary-color-dark-mode;
    }
  }
  .current-value {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    margin-bottom: 1rem;
    padding: 0.8rem 1rem;
    background-color: rgba($text-color-light-mode, 0.03);
    border-radius: 6px;
    border: 1px solid rgba($text-color-light-mode, 0.1);
    .dark-mode & {
      background-color: rgba($text-color-dark-mode, 0.05);
      border-color: rgba($text-color-dark-mode, 0.1);
    }
    .label {
      font-weight: 500;
      color: rgba($text-color-light-mode, 0.7);
      .dark-mode & {
        color: rgba($text-color-dark-mode, 0.7);
      }
    }
    .value {
      font-weight: 600;
      color: $primary-color-light-mode;
      .dark-mode & {
        color: $primary-color-dark-mode;
      }
    }
  }
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;

  @media (max-width: 576px) {
    flex-direction: column;
  }
}

.btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.8rem 1.5rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  text-decoration: none;

  &:disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }

  &.btn-primary {
    background-color: $primary-color-light-mode;
    color: white;

    &:hover:not(:disabled) {
      background-color: darken($primary-color-light-mode, 5%);
    }

    .dark-mode & {
      background-color: $primary-color-dark-mode;
      color: black;

      &:hover:not(:disabled) {
        background-color: lighten($primary-color-dark-mode, 5%);
      }
    }
  }

  &.btn-secondary {
    background-color: rgba($text-color-light-mode, 0.1);
    color: $text-color-light-mode;

    &:hover:not(:disabled) {
      background-color: rgba($text-color-light-mode, 0.15);
    }

    .dark-mode & {
      background-color: rgba($text-color-dark-mode, 0.1);
      color: $text-color-dark-mode;

      &:hover:not(:disabled) {
        background-color: rgba($text-color-dark-mode, 0.15);
      }
    }
  }

  &.btn-danger {
    background-color: rgba(#e74c3c, 0.1);
    color: #e74c3c;

    &:hover:not(:disabled) {
      background-color: rgba(#e74c3c, 0.15);
    }
  }

  .rotating {
    animation: rotate 1s linear infinite;
  }
}

@keyframes rotate {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

// Modal styles
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal-content {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  max-width: 500px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;

  .dark-mode & {
    background-color: $background-color-dark-mode;
  }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem;
  border-bottom: 1px solid rgba($text-color-light-mode, 0.1);

  .dark-mode & {
    border-bottom-color: rgba($text-color-dark-mode, 0.1);
  }

  h3 {
    margin: 0;
    color: $text-color-light-mode;

    .dark-mode & {
      color: $text-color-dark-mode;
    }
  }

  .modal-close {
    background: none;
    border: none;
    font-size: 1.5rem;
    color: rgba($text-color-light-mode, 0.6);
    cursor: pointer;
    padding: 0.25rem;

    &:hover {
      color: $text-color-light-mode;
    }

    .dark-mode & {
      color: rgba($text-color-dark-mode, 0.6);

      &:hover {
        color: $text-color-dark-mode;
      }
    }
  }
}

.modal-body {
  padding: 1.5rem;

  p {
    color: $text-color-light-mode;
    margin-bottom: 1rem;

    .dark-mode & {
      color: $text-color-dark-mode;
    }

    &.warning-text {
      color: #e74c3c;
      font-weight: 500;
    }
  }
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  padding: 1.5rem;
  border-top: 1px solid rgba($text-color-light-mode, 0.1);

  .dark-mode & {
    border-top-color: rgba($text-color-dark-mode, 0.1);
  }

  @media (max-width: 576px) {
    flex-direction: column;
  }
}
</style>
