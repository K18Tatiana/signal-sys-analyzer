<template>
  <div class="home-page">
    <h1 class="page-title">IDENTIFICACIÓN DE SISTEMAS DE CONTROL</h1>

    <div class="upload-section" :class="{ 'hidden': isProcessing }">
      <UploadArea
        v-model="selectedFile"
        @error="handleFileError"
      />

      <InputField
        v-model="inputVoltage"
        :error="voltageError"
        label="Ingrese el voltaje de entrada:"
        :max="12"
        :min="4"
        placeholder="5"
        :show-range-alert="true"
        type="number"
        unit="V"
        @validation="handleVoltageValidation"
      />

      <button
        class="submit-button"
        :disabled="!isFormValid"
        @click="processData"
      >
        Enviar
      </button>
    </div>

    <div v-if="isProcessing" class="processing-indicator">
      <div class="spinner"></div>
      <p>Identificando sistema de control...</p>
    </div>

    <ResultsDisplay
      v-if="results"
      :results="results"
      @reset="resetForm"
    />
  </div>
</template>

<script setup>
  import { computed, ref } from 'vue';
  import UploadArea from '../components/UploadArea.vue';
  import InputField from '../components/InputField.vue';
  import ResultsDisplay from '../components/ResultsDisplay.vue';

  // Estado
  const selectedFile = ref(null);
  const inputVoltage = ref(5);
  const voltageError = ref('');
  const error = ref('');
  const isProcessing = ref(false);
  const results = ref(null);

  // Validaciones
  const handleFileError = msg => {
    error.value = msg;
  };

  const handleVoltageValidation = msg => {
    voltageError.value = msg;
  };

  // Verificar si el formulario es válido
  const isFormValid = computed(() => {
    return selectedFile.value && inputVoltage.value >= 4 && inputVoltage.value <= 12 && !voltageError.value;
  });

  // Métodos para procesar datos
  const processData = async () => {
    // No procesar si el formulario no es válido
    if (!isFormValid.value) {
      if (!selectedFile.value) {
        error.value = 'Debe seleccionar un archivo CSV';
        return;
      }

      if (inputVoltage.value < 4 || inputVoltage.value > 12 || voltageError.value) {
        voltageError.value = `El voltaje debe estar entre 4 y 12 V`;
        return;
      }

      return;
    }

    try {
      isProcessing.value = true;
      error.value = '';

      // Simular procesamiento (en una aplicación real, esto sería una llamada a la API)
      await new Promise(resolve => setTimeout(resolve, 1500));

      // Simular resultado
      results.value = {
        type: 'Sistema subamortiguado',
        description: 'El sistema presenta oscilaciones antes de estabilizarse. Esto indica que los polos del sistema tienen una parte imaginaria significativa, lo cual es típico en sistemas donde la respuesta transitoria incluye sobrepaso. Este tipo de comportamiento es común en sistemas que requieren rapidez de respuesta, como servomecanismos.',
        poles: [
          '-0.5 + 1.2j',
          '-0.5 - 1.2j',
          '-2.0',
        ],
        // Datos para generar el gráfico en el frontend
        graphData: {
          labels: ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '10'],
          datasets: [
            {
              label: 'Respuesta',
              data: [0, 0.2, 0.6, 1.1, 1.3, 1.2, 1.1, 1.0, 1.0, 1.0, 1.0],
              borderColor: '#1470AF',
              backgroundColor: 'rgba(20, 112, 175, 0.1)',
              tension: 0.4,
            },
          ],
        },
      };
    } catch (err) {
      error.value = err.message || 'Ha ocurrido un error durante el procesamiento.';
      results.value = null;
    } finally {
      isProcessing.value = false;
    }
  };

  const resetForm = () => {
    selectedFile.value = null;
    inputVoltage.value = 5;
    voltageError.value = '';
    error.value = '';
    results.value = null;
  };
</script>

<style lang="scss">
@import '../styles/variables.scss';

.home-page {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 1rem;

  .page-title {
    font-family: $font-primary;
    font-size: 2rem;
    font-weight: 700;
    margin-bottom: 4rem;
    text-align: center;
    color: $primary-color-light-mode;

    .dark-mode & {
      color: $primary-color-dark-mode;
    }
  }

  .error-message {
    width: 100%;
    max-width: 600px;
    background-color: #fff3f3;
    border-left: 4px solid #e74c3c;
    padding: 1rem;
    margin-bottom: 1.5rem;
    border-radius: 4px;
    color: #e74c3c;

    .dark-mode & {
      background-color: rgba(231, 76, 60, 0.2);
      color: #ff6b6b;
    }
  }

  .upload-section {
    width: 100%;
    max-width: 500px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
    margin-bottom: 2rem;

    &.hidden {
      display: none;
    }
  }

  .processing-indicator {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 2rem;

    .spinner {
      width: 50px;
      height: 50px;
      border: 5px solid rgba($primary-color-light-mode, 0.2);
      border-radius: 50%;
      border-top-color: $primary-color-light-mode;
      animation: spin 1s ease-in-out infinite;
      margin-bottom: 1rem;

      .dark-mode & {
        border-color: rgba($primary-color-dark-mode, 0.2);
        border-top-color: $primary-color-dark-mode;
      }
    }

    p {
      font-size: 1rem;
      color: $secondary-color-light-mode;

      .dark-mode & {
        color: $secondary-color-dark-mode;
      }
    }
  }

  .submit-button {
    padding: 0.75rem 2rem;
    background-color: $primary-color-light-mode;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: background-color 0.3s;
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 150px;
    margin-top: 0.5rem;

    &:hover:not(:disabled) {
      background-color: darken($primary-color-light-mode, 10%);
    }

    &:disabled {
      background-color: #ccc;
      cursor: not-allowed;
    }

    .dark-mode & {
      background-color: $primary-color-dark-mode;
      color: $background-color-dark-mode;

      &:hover:not(:disabled) {
        background-color: darken($primary-color-dark-mode, 10%);
      }

      &:disabled {
        background-color: #555;
        color: #888;
      }
    }
  }

  .btn-sm {
    padding: 0.25rem 0.5rem;
    font-size: 0.875rem;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 768px) {
  .home-page {
    .page-title {
      font-size: 1.5rem;
    }
  }
}

@media (max-width: 576px) {
  .home-page {
    .page-title {
      font-size: 1.2rem;
    }
  }
}
</style>
