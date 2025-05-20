import axios from 'axios';

// Configura la URL base según tu entorno
const API_URL = process.env.VUE_APP_API_URL || 'http://localhost:8081/api';

// Crea una instancia personalizada de axios
const apiClient = axios.create({
  baseURL: API_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
});

// Interceptor para añadir el token de autenticación a las solicitudes
apiClient.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// Servicio de autenticación
export const authService = {
  /**
   * Registra a un nuevo usuario
   * @param {Object} userData - Datos del usuario (username, email, contraseña)
   * @returns {Promise} - Promesa con la respuesta del servidor
   */
  register: async userData => {
    try {
      // Adaptamos el campo name a username que espera el servidor
      const requestData = {
        username: userData.name,
        email: userData.email,
        password: userData.password,
      };

      const response = await apiClient.post('/register', requestData);

      // Si el registro es exitoso y el backend devuelve un token, lo guardamos
      if (response.data.token) {
        localStorage.setItem('token', response.data.token);
        localStorage.setItem('user', JSON.stringify(response.data.user));
      }

      return response.data;
    } catch (error) {
      throw handleError(error);
    }
  },

  /**
   * Inicia sesión de un usuario
   * @param {Object} credentials - Credenciales (email, contraseña)
   * @returns {Promise} - Promesa con la respuesta del servidor que incluye el token
   */
  login: async credentials => {
    try {

      console.log('Intentando iniciar sesión con:', credentials);
      console.log('URL completa:', `${API_URL}/login`);

      const response = await apiClient.post('/login', credentials);
      console.log('Respuesta exitosa:', response.data);

      // Si la autenticación es exitosa, guarda el token en localStorage
      if (response.data.token) {
        localStorage.setItem('token', response.data.token);
        localStorage.setItem('user', JSON.stringify(response.data.user));
      }

      return response.data;
    } catch (error) {
      console.error('Error completo:', error);

      if (error.response) {
        console.error('Respuesta de error:', error.response.data);
      } else if (error.request) {
        console.error('Objeto request:', error.request);
      }

      throw handleError(error);
    }
  },

  /**
   * Cierra la sesión del usuario
   */
  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  /**
   * Verifica si el usuario está autenticado
   * @returns {boolean} - True si el usuario está autenticado
   */
  isAuthenticated: () => {
    return !!localStorage.getItem('token');
  },

  /**
   * Obtiene el usuario actual
   * @returns {Object|null} - Datos del usuario o null si no hay usuario
   */
  getCurrentUser: () => {
    const userJson = localStorage.getItem('user');
    return userJson ? JSON.parse(userJson) : null;
  },
};

/**
 * Maneja los errores de las solicitudes HTTP
 * @param {Error} error - Error de axios
 * @returns {Object} - Error formateado
 */
function handleError (error) {
  let errorMessage = 'Ha ocurrido un error inesperado';

  if (error.response) {
    // El servidor respondió con un código de estado fuera del rango 2xx
    if (error.response.data && error.response.data.message) {
      errorMessage = error.response.data.message;
    } else if (error.response.data && typeof error.response.data === 'string') {
      // Si el servidor devuelve un string directo (como "Email o contraseña incorrectos")
      errorMessage = error.response.data;
    } else {
      switch (error.response.status) {
        case 400:
          errorMessage = 'Solicitud incorrecta';
          break;
        case 401:
          errorMessage = 'Email o contraseña incorrectos';
          break;
        case 403:
          errorMessage = 'Acceso prohibido';
          break;
        case 404:
          errorMessage = 'Recurso no encontrado';
          break;
        case 409:
          errorMessage = 'El correo electrónico ya está registrado';
          break;
        case 500:
          errorMessage = 'Error interno del servidor';
          break;
        default:
          errorMessage = `Error ${error.response.status}`;
      }
    }
  } else if (error.request) {
    // La solicitud fue realizada pero no se recibió respuesta
    errorMessage = 'No se recibió respuesta del servidor';
  }

  return {
    message: errorMessage,
    originalError: error,
  };
}

export default authService;
