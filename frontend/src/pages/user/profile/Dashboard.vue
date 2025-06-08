<template>
  <div class="dashboard-page">
    <div class="dashboard-container">
      <div class="welcome-card">
        <div class="welcome-header">
          <div class="user-avatar">
            <span class="avatar-text">{{ userInitials }}</span>
          </div>
          <div class="welcome-text">
            <h1>Bienvenido, {{ user?.username }}</h1>
            <p>Esta es tu área personal donde puedes gestionar tus análisis y documentos.</p>
          </div>
        </div>
      </div>

      <div class="dashboard-grid">
        <!-- Estadísticas -->
        <div class="stats-card dashboard-card">
          <h2 class="card-title">Resumen</h2>
          <div class="stats-grid">
            <div class="stat-item">
              <div class="stat-icon">
                <i class="bx bx-file"></i>
              </div>
              <div class="stat-info">
                <h3>{{ stats.documentsCount }}</h3>
                <p>Documentos</p>
              </div>
            </div>
            <div class="stat-item">
              <div class="stat-icon">
                <i class="bx bx-analyse"></i>
              </div>
              <div class="stat-info">
                <h3>{{ stats.analysisCount }}</h3>
                <p>Análisis</p>
              </div>
            </div>
            <div class="stat-item">
              <div class="stat-icon">
                <i class="bx bx-calendar"></i>
              </div>
              <div class="stat-info">
                <h3>{{ formatDate(user?.created_at) }}</h3>
                <p>Miembro desde</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Acciones rápidas -->
        <div class="actions-card dashboard-card">
          <h2 class="card-title">Acciones rápidas</h2>
          <div class="action-buttons">
            <router-link class="action-button primary" to="/home">
              <i class="bx bx-plus"></i>
              <span>Nuevo análisis</span>
            </router-link>
            <router-link class="action-button secondary" to="/user/documents">
              <i class="bx bx-file"></i>
              <span>Ver documentos</span>
            </router-link>
            <router-link class="action-button secondary" to="/profile">
              <i class="bx bx-user"></i>
              <span>Editar perfil</span>
            </router-link>
          </div>
        </div>

        <!-- Actividad reciente -->
        <div class="activity-card dashboard-card">
          <h2 class="card-title">Actividad reciente</h2>
          <div v-if="recentActivity.length === 0" class="empty-state">
            <i class="bx bx-calendar-check empty-icon"></i>
            <p>No tienes actividad reciente</p>
          </div>
          <ul v-else class="activity-list">
            <li v-for="(activity, index) in recentActivity" :key="index" class="activity-item">
              <div class="activity-icon">
                <i :class="getActivityIcon(activity.type)"></i>
              </div>
              <div class="activity-content">
                <p class="activity-text">{{ activity.description }}</p>
                <p class="activity-date">{{ formatActivityDate(activity.date) }}</p>
              </div>
            </li>
          </ul>
        </div>

        <!-- Documentos recientes -->
        <div class="documents-card dashboard-card">
          <h2 class="card-title">Documentos recientes</h2>
          <div v-if="recentDocuments.length === 0" class="empty-state">
            <i class="bx bx-file empty-icon"></i>
            <p>No tienes documentos aún</p>
            <router-link class="upload-button" to="/home">
              <i class="bx bx-upload"></i>
              <span>Subir documento</span>
            </router-link>
          </div>
          <ul v-else class="document-list">
            <li v-for="(document, index) in recentDocuments" :key="index" class="document-item">
              <div class="document-icon">
                <i class="bx bx-file"></i>
              </div>
              <div class="document-info">
                <p class="document-name">{{ document.name }}</p>
                <p class="document-date">{{ formatDate(document.date) }}</p>
              </div>
              <router-link class="document-action" to="/user/documents">
                <i class="bx bx-chevron-right"></i>
              </router-link>
            </li>
          </ul>
        </div>
      </div>

      <div class="logout-section">
        <button class="logout-button" @click="handleLogout">
          <i class="bx bx-log-out"></i>
          <span>Cerrar sesión</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { computed, onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router';
  import { useAuth } from '@/stores/auth.js';

  // Configuración api
  const API_URL = import.meta.env.VITE_API_URL

  const router = useRouter();
  const { user, logout, isAuthenticated } = useAuth();

  // Estados
  const stats = ref({
    documentsCount: 0,
    analysisCount: 0,
  });

  const recentActivity = ref([]);
  const recentDocuments = ref([]);
  const loading = ref(true);
  const error = ref('');

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

  // Cargar estadísticas del usuario
  const loadUserStats = async () => {
    try {
      const response = await fetch(`${API_URL}/user/stats`, {
        method: 'GET',
        headers: getAuthHeaders(),
      });

      if (!response.ok) {
        throw new Error('Error al cargar estadísticas');
      }

      const data = await response.json();
      stats.value = {
        documentsCount: data.documents_count,
        analysisCount: data.analysis_count,
      };
    } catch (err) {
      console.error('Error cargando estadísticas:', err);
      error.value = 'Error al cargar estadísticas';
    }
  };

  // Cargar actividad reciente
  const loadRecentActivity = async () => {
    try {
      const response = await fetch(`${API_URL}/user/recent-activity`, {
        method: 'GET',
        headers: getAuthHeaders(),
      });

      if (!response.ok) {
        throw new Error('Error al cargar actividad reciente');
      }

      const data = await response.json();
      recentActivity.value = data;
    } catch (err) {
      console.error('Error cargando actividad:', err);
    }
  };

  // Cargar documentos recientes
  const loadRecentDocuments = async () => {
    try {
      const response = await fetch(`${API_URL}/user/recent-documents`, {
        method: 'GET',
        headers: getAuthHeaders(),
      });

      if (!response.ok) {
        throw new Error('Error al cargar documentos recientes');
      }

      const data = await response.json();
      recentDocuments.value = data;
    } catch (err) {
      console.error('Error cargando documentos:', err);
    }
  };

  // Función principal para cargar datos del dashboard
  const loadDashboardData = async () => {
    loading.value = true;
    try {
      await Promise.all([
        loadUserStats(),
        loadRecentActivity(),
        loadRecentDocuments(),
      ]);
    } catch (err) {
      error.value = 'Error al cargar datos del dashboard: ' + err;
    } finally {
      loading.value = false;
    }
  };

  // Calcular iniciales del usuario para el avatar
  const userInitials = computed(() => {
    if (!user.value || !user.value.username) return '';
    return user.value.username.charAt(0).toUpperCase();
  });

  onMounted(async () => {
    // Verificar si hay usuario autenticado
    if (!isAuthenticated.value) {
      router.push('/login');
      return;
    }

    // Cargar datos del dashboard
    await loadDashboardData();
  });

  // Función para cerrar sesión
  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  // Funciones de utilidad para formatear fechas
  const formatDate = dateString => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('es-ES', {
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    });
  };

  const formatActivityDate = dateString => {
    if (!dateString) return '';
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now - date);
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

    if (diffDays === 1) {
      return 'Ayer';
    } else if (diffDays < 7) {
      return `Hace ${diffDays} días`;
    } else {
      return date.toLocaleDateString('es-ES', {
        day: '2-digit',
        month: 'short',
      });
    }
  };

  // Función para determinar el icono según el tipo de actividad
  const getActivityIcon = type => {
    switch (type) {
      case 'document':
        return 'bx bx-file';
      case 'analysis':
        return 'bx bx-analyse';
      case 'login':
        return 'bx bx-log-in';
      default:
        return 'bx bx-activity';
    }
  };
</script>

<style lang="scss" scoped>
@import '/src/styles/variables.scss';

.dashboard-page {
  min-height: calc(100vh - 70px);
  background-color: #f8f9fa;
  padding: 2rem 1rem;

  .dark-mode & {
    background-color: darken($background-color-dark-mode, 2%);
  }
}

.dashboard-container {
  max-width: 1200px;
  margin: 0 auto;
}

.welcome-card {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
  padding: 2rem;
  margin-bottom: 2rem;

  .dark-mode & {
    background-color: $background-color-dark-mode;
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.2);
  }

  .welcome-header {
    display: flex;
    align-items: center;
    gap: 1.5rem;

    .user-avatar {
      width: 64px;
      height: 64px;
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
        font-size: 1.8rem;
        font-weight: 600;

        .dark-mode & {
          color: black;
        }
      }
    }

    .welcome-text {
      h1 {
        font-size: 1.8rem;
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

        .dark-mode & {
          color: rgba($text-color-dark-mode, 0.7);
        }
      }
    }

    @media (max-width: 576px) {
      flex-direction: column;
      text-align: center;

      .welcome-text {
        h1 {
          font-size: 1.5rem;
        }
      }
    }
  }
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
  margin-bottom: 2rem;

  @media (max-width: 992px) {
    grid-template-columns: 1fr;
  }

  .dashboard-card {
    background-color: white;
    border-radius: 12px;
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
    padding: 1.5rem;

    .dark-mode & {
      background-color: $background-color-dark-mode;
      box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
    }

    .card-title {
      font-size: 1.2rem;
      font-weight: 600;
      color: $secondary-color-light-mode;
      margin-top: 0;
      margin-bottom: 1.5rem;

      .dark-mode & {
        color: $secondary-color-dark-mode;
      }
    }
  }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 1rem;

    .stat-icon {
      width: 48px;
      height: 48px;
      border-radius: 8px;
      background-color: rgba($primary-color-light-mode, 0.1);
      display: flex;
      align-items: center;
      justify-content: center;

      i {
        font-size: 1.5rem;
        color: $primary-color-light-mode;
      }

      .dark-mode & {
        background-color: rgba($primary-color-dark-mode, 0.2);

        i {
          color: $primary-color-dark-mode;
        }
      }
    }

    .stat-info {
      h3 {
        margin: 0;
        font-size: 1.5rem;
        font-weight: 700;
        color: $text-color-light-mode;

        .dark-mode & {
          color: $text-color-dark-mode;
        }
      }

      p {
        margin: 0.2rem 0 0;
        font-size: 0.9rem;
        color: rgba($text-color-light-mode, 0.7);

        .dark-mode & {
          color: rgba($text-color-dark-mode, 0.7);
        }
      }
    }
  }
}

.action-buttons {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }

  .action-button {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 1.2rem;
    border-radius: 8px;
    text-decoration: none;
    transition: all 0.2s;

    i {
      font-size: 1.8rem;
    }

    span {
      font-size: 0.9rem;
      font-weight: 500;
    }

    &.primary {
      background-color: $primary-color-light-mode;
      color: white;

      &:hover {
        background-color: darken($primary-color-light-mode, 5%);
        transform: translateY(-3px);
      }

      .dark-mode & {
        background-color: $primary-color-dark-mode;
        color: black;

        &:hover {
          background-color: lighten($primary-color-dark-mode, 5%);
        }
      }
    }

    &.secondary {
      background-color: rgba($primary-color-light-mode, 0.1);
      color: $primary-color-light-mode;

      &:hover {
        background-color: rgba($primary-color-light-mode, 0.15);
        transform: translateY(-3px);
      }

      .dark-mode & {
        background-color: rgba($primary-color-dark-mode, 0.2);
        color: $primary-color-dark-mode;

        &:hover {
          background-color: rgba($primary-color-dark-mode, 0.25);
        }
      }
    }
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  text-align: center;

  .empty-icon {
    font-size: 3rem;
    color: rgba($text-color-light-mode, 0.2);
    margin-bottom: 1rem;

    .dark-mode & {
      color: rgba($text-color-dark-mode, 0.2);
    }
  }

  p {
    margin: 0;
    color: rgba($text-color-light-mode, 0.6);
    font-size: 0.95rem;

    .dark-mode & {
      color: rgba($text-color-dark-mode, 0.6);
    }
  }

  .upload-button {
    margin-top: 1rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background-color: rgba($primary-color-light-mode, 0.1);
    color: $primary-color-light-mode;
    border: none;
    border-radius: 6px;
    padding: 0.7rem 1.2rem;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    text-decoration: none;

    &:hover {
      background-color: rgba($primary-color-light-mode, 0.15);
    }

    .dark-mode & {
      background-color: rgba($primary-color-dark-mode, 0.2);
      color: $primary-color-dark-mode;

      &:hover {
        background-color: rgba($primary-color-dark-mode, 0.25);
      }
    }
  }
}

.activity-list, .document-list {
  list-style: none;
  padding: 0;
  margin: 0;

  .activity-item, .document-item {
    display: flex;
    align-items: flex-start;
    gap: 1rem;
    padding: 1rem;
    border-bottom: 1px solid rgba(0, 0, 0, 0.05);

    &:last-child {
      border-bottom: none;
    }

    .dark-mode & {
      border-bottom-color: rgba(255, 255, 255, 0.05);
    }
  }
}

.activity-item {
  .activity-icon {
    flex-shrink: 0;
    width: 40px;
    height: 40px;
    border-radius: 8px;
    background-color: rgba($primary-color-light-mode, 0.1);
    display: flex;
    align-items: center;
    justify-content: center;

    i {
      font-size: 1.2rem;
      color: $primary-color-light-mode;
    }

    .dark-mode & {
      background-color: rgba($primary-color-dark-mode, 0.2);

      i {
        color: $primary-color-dark-mode;
      }
    }
  }

  .activity-content {
    flex: 1;

    .activity-text {
      margin: 0;
      font-size: 0.95rem;
      color: $text-color-light-mode;

      .dark-mode & {
        color: $text-color-dark-mode;
      }
    }

    .activity-date {
      margin: 0.3rem 0 0;
      font-size: 0.8rem;
      color: rgba($text-color-light-mode, 0.6);

      .dark-mode & {
        color: rgba($text-color-dark-mode, 0.6);
      }
    }
  }
}

.document-item {
  .document-icon {
    flex-shrink: 0;
    width: 40px;
    height: 40px;
    border-radius: 8px;
    background-color: rgba($primary-color-light-mode, 0.1);
    display: flex;
    align-items: center;
    justify-content: center;

    i {
      font-size: 1.2rem;
      color: $primary-color-light-mode;
    }

    .dark-mode & {
      background-color: rgba($primary-color-dark-mode, 0.2);

      i {
        color: $primary-color-dark-mode;
      }
    }
  }

  .document-info {
    flex: 1;

    .document-name {
      margin: 0;
      font-size: 0.95rem;
      font-weight: 500;
      color: $text-color-light-mode;

      .dark-mode & {
        color: $text-color-dark-mode;
      }
    }

    .document-date {
      margin: 0.3rem 0 0;
      font-size: 0.8rem;
      color: rgba($text-color-light-mode, 0.6);

      .dark-mode & {
        color: rgba($text-color-dark-mode, 0.6);
      }
    }
  }

  .document-action {
    flex-shrink: 0;
    background: none;
    border: none;
    color: rgba($text-color-light-mode, 0.5);
    font-size: 1.2rem;
    cursor: pointer;
    padding: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: color 0.2s;
    text-decoration: none;

    &:hover {
      color: $primary-color-light-mode;
    }

    .dark-mode & {
      color: rgba($text-color-dark-mode, 0.5);

      &:hover {
        color: $primary-color-dark-mode;
      }
    }
  }
}

.logout-section {
  display: flex;
  justify-content: center;
  margin-top: 2rem;

  .logout-button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background-color: rgba(#e74c3c, 0.1);
    color: #e74c3c;
    border: none;
    border-radius: 6px;
    padding: 0.7rem 1.5rem;
    font-size: 0.95rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background-color: rgba(#e74c3c, 0.15);
    }
  }
}
</style>
