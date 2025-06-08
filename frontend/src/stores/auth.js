import { computed, ref } from 'vue';
import { authService } from '../services/auth.service.js';

const user = ref(null);
const token = ref(null);

const initializeAuth = () => {
  const storedToken = localStorage.getItem('token');
  const storedUser = localStorage.getItem('user');

  if (storedToken) {
    token.value = storedToken;
  }

  if (storedUser) {
    try {
      user.value = JSON.parse(storedUser);
    } catch (error) {
      console.error('Error parsing stored user:', error);
      localStorage.removeItem('user');
    }
  }
};

const isAuthenticated = computed(() => {
  return !!token.value;
});

// Funciones para manejar la autenticación
const setAuth = authData => {
  token.value = authData.token;
  user.value = authData.user;

  localStorage.setItem('token', authData.token);
  localStorage.setItem('user', JSON.stringify(authData.user));
};

const clearAuth = () => {
  token.value = null;
  user.value = null;

  localStorage.removeItem('token');
  localStorage.removeItem('user');
};

const useAuth = () => {
  return {
    user: computed(() => user.value),
    token: computed(() => token.value),
    isAuthenticated,

    async login (credentials) {
      try {
        const response = await authService.login(credentials);
        setAuth(response);
        return response;
      } catch (error) {
        throw error;
      }
    },

    async register (userData) {
      try {
        const response = await authService.register(userData);
        if (response.token) {
          setAuth(response);
        }
        return response;
      } catch (error) {
        throw error;
      }
    },

    logout () {
      authService.logout();
      clearAuth();
    },

    getCurrentUser () {
      return user.value;
    },

    initialize: initializeAuth,
  };
};

export { useAuth }
