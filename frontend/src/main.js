// main.js
import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';
import { useAuth } from './stores/auth.js';
import './styles/globals.scss';

// Creación de la aplicación Vue
const app = createApp(App);

// Crear y usar Pinia para el estado global
const pinia = createPinia();
app.use(pinia);

// Uso del enrutador
app.use(router);

// Inicializar el estado de autenticación al arrancar la aplicación
const { initialize } = useAuth();
initialize();

// Montaje de la aplicación en el elemento con id 'app'
app.mount('#app');
