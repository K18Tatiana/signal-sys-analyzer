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

      <div v-if="isUserAuthenticated" class="comment-field">
        <label for="analysis-comment">Comentario del análisis (opcional):</label>
        <textarea
          id="analysis-comment"
          v-model="analysisComment"
          class="comment-textarea"
          maxlength="500"
          placeholder="Describe este análisis o agrega notas relevantes..."
          rows="3"
        />
        <small class="char-count">{{ analysisComment.length }}/500</small>
      </div>

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
      <p>{{ processingMessage }}</p>
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
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
  import { authService } from '../services/auth.service.js';

  // Configuración api
  const API_URL = import.meta.env.VITE_API_URL

  // Estado
  const selectedFile = ref(null);
  const inputVoltage = ref(5);
  const voltageError = ref('');
  const error = ref('');
  const isProcessing = ref(false);
  const processingMessage = ref('');
  const results = ref(null);
  const uploadedDocumentId = ref(null);
  const analysisComment = ref('');

  // Función para obtener headers de autorización si el usuario está logueado
  const getAuthHeaders = () => {
    const headers = {};

    // Si el usuario está autenticado, agregar el token
    if (authService.isAuthenticated()) {
      const token = localStorage.getItem('token');
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }

    return headers;
  };

  // Validaciones
  const handleFileError = msg => {
    error.value = msg;
  };

  const handleVoltageValidation = msg => {
    voltageError.value = msg;
  };

  // Verificar si el formulario es válido
  const isFormValid = computed(() => {
    return selectedFile.value &&
      inputVoltage.value >= 4 &&
      inputVoltage.value <= 12 &&
      !voltageError.value;
  });

  // Verificar si el usuario está autenticado
  const isUserAuthenticated = computed(() => {
    return authService.isAuthenticated();
  });

  // Subir documento al servidor (con autenticación opcional)
  const uploadDocument = async () => {
    const formData = new FormData();
    formData.append('file', selectedFile.value);

    try {
      // Crear objeto de configuración para fetch
      const fetchConfig = {
        method: 'POST',
        body: formData,
      };

      // Agregar headers de autorización si el usuario está logueado
      const authHeaders = getAuthHeaders();
      if (Object.keys(authHeaders).length > 0) {
        fetchConfig.headers = authHeaders;
      }

      const response = await fetch(`${API_URL}/documents`, fetchConfig);

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Error al subir el archivo');
      }

      const data = await response.json();
      return data;
    } catch (err) {
      throw new Error(`Error de conexión: ${err.message}`);
    }
  };

  // Crear solicitud de análisis (con autenticación opcional)
  const createAnalysisRequest = async documentId => {
    try {
      const requestBody = {
        document_id: parseInt(documentId),
        input_voltage: parseFloat(inputVoltage.value),
      };

      // Solo agregar comentario si el usuario está autenticado
      if (isUserAuthenticated.value && analysisComment.value.trim()) {
        requestBody.comment = analysisComment.value.trim();
      }

      const fetchConfig = {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders(),
        },
        body: JSON.stringify(requestBody),
      };

      const response = await fetch(`${API_URL}/analysis`, fetchConfig);

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Error al crear la solicitud de análisis');
      }

      const data = await response.json();
      return data;
    } catch (err) {
      throw new Error(`Error al crear análisis: ${err.message}`);
    }
  };

  // Obtener resultado del análisis (con autenticación opcional)
  const getAnalysisResult = async analysisId => {
    try {
      const fetchConfig = {
        method: 'GET',
        headers: getAuthHeaders(),
      };

      const response = await fetch(`${API_URL}/analysis/${analysisId}`, fetchConfig);

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Error al obtener el resultado');
      }

      const data = await response.json();
      return data;
    } catch (err) {
      throw new Error(`Error al obtener resultado: ${err.message}`);
    }
  };

  // Polling para verificar si el análisis está completo
  const pollAnalysisResult = async (analysisId, maxAttempts = 10) => {
    for (let attempt = 0; attempt < maxAttempts; attempt++) {
      try {
        const result = await getAnalysisResult(analysisId);

        if (result.analysis_request && result.analysis_request.is_processed === true) {
          return result;
        }

        if (result.result) {
          return result;
        }

        const timeElapsed = (attempt + 1) * 2;
        processingMessage.value = `Procesando análisis... (${attempt + 1}/${maxAttempts}) - ${timeElapsed}s transcurridos`;

        await new Promise(resolve => setTimeout(resolve, 2000));

      } catch (err) {
        if (attempt === maxAttempts - 1) {
          throw new Error(`Error al obtener resultado del análisis: ${err.message}`);
        }

        processingMessage.value = `Reintentando... (${attempt + 1}/${maxAttempts})`;
        await new Promise(resolve => setTimeout(resolve, 2000));
      }
    }

    throw new Error(`El análisis está tomando más tiempo del esperado (${maxAttempts * 2}s). Por favor, intenta más tarde.`);
  };

  // Método principal para procesar datos
  const processData = async () => {
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

      // Mostrar información de estado de autenticación
      const isLoggedIn = authService.isAuthenticated();
      const currentUser = authService.getCurrentUser();

      if (isLoggedIn && currentUser) {
        processingMessage.value = `Subiendo archivo (como ${currentUser.username})...`;
      } else {
        processingMessage.value = 'Subiendo archivo (como usuario anónimo)...';
      }

      // Paso 1: Subir documento
      const documentResponse = await uploadDocument();
      uploadedDocumentId.value = documentResponse.id;

      processingMessage.value = 'Iniciando análisis del sistema...';

      // Paso 2: Crear solicitud de análisis
      const analysisResponse = await createAnalysisRequest(documentResponse.id);

      processingMessage.value = 'Identificando sistema de control...';

      // Paso 3: Esperar y obtener resultado
      const finalResult = await pollAnalysisResult(analysisResponse.id);

      // Paso 4: Procesar y mostrar resultados
      if (finalResult.result) {
        const resultData = finalResult.result;

        let poles = [];
        if (resultData.poles && resultData.poles.polos) {
          poles = resultData.poles.polos.map(pole => {
            if (pole.imag === 0) {
              return pole.real.toString();
            } else if (pole.imag > 0) {
              return `${pole.real} + ${pole.imag}j`;
            } else {
              return `${pole.real} - ${Math.abs(pole.imag)}j`;
            }
          });
        }

        results.value = {
          type: resultData.system_type || 'Sistema no identificado',
          description: getSystemDescription(resultData.system_type, resultData), // CAMBIO AQUÍ
          poles,
          graphData: processGraphData(resultData.raw_data, resultData.graph_data),
          technicalSummary: resultData.technical_summary || null, // NUEVO CAMPO
          analysisId: analysisResponse.id,
          documentId: documentResponse.id,
          rawData: resultData.raw_data,
          userInfo: isLoggedIn ? {
            isAuthenticated: true,
            username: currentUser?.username,
          } : {
            isAuthenticated: false,
          },
        };
      } else {
        throw new Error('No se obtuvieron resultados del análisis. Estructura de respuesta inesperada.');
      }

    } catch (err) {
      error.value = err.message || 'Ha ocurrido un error durante el procesamiento.';
      results.value = null;
    } finally {
      isProcessing.value = false;
      processingMessage.value = '';
    }
  };

  const getSystemDescription = (systemType, resultData = null) => {
    // Si hay datos del resultado con descripción del backend, usarla
    if (resultData && resultData.description) {
      return resultData.description;
    }

    // Fallback a descripciones estáticas si no hay descripción del backend
    const descriptions = {
      'subamortiguado': 'Sistema subamortiguado - presenta oscilaciones antes de estabilizarse. Consulte los datos técnicos para más detalles.',
      'sobreamortiguado': 'Sistema sobreamortiguado - se aproxima gradualmente al valor final sin oscilaciones. Consulte los datos técnicos para más detalles.',
    };

    return descriptions[systemType] || 'Sistema identificado con características específicas según el análisis realizado.';
  };

  const processGraphData = (rawData, graphDataString) => {
    if (graphDataString) {
      try {
        const graphData = typeof graphDataString === 'string' ? JSON.parse(graphDataString) : graphDataString;
        if (graphData.time && graphData.output) {
          const reducedData = reducePointsForVisualization(graphData.time, graphData.output, 150);
          return {
            labels: reducedData.time.map(t => t.toFixed(3)),
            datasets: [
              {
                label: 'Respuesta del Sistema',
                data: reducedData.output,
                borderColor: '#1470AF',
                backgroundColor: 'rgba(20, 112, 175, 0.1)',
                tension: 0.2,
                pointRadius: 0,
                pointHoverRadius: 4,
                borderWidth: 2,
              },
            ],
          };
        }
      } catch (err) {
        console.warn('Error procesando graph_data:', err);
      }
    }
    return null;
  };

  const reducePointsForVisualization = (timeArray, outputArray, maxPoints) => {
    if (timeArray.length <= maxPoints) {
      return { time: timeArray, output: outputArray };
    }

    const reducedTime = [];
    const reducedOutput = [];
    const step = Math.floor(timeArray.length / maxPoints);

    reducedTime.push(timeArray[0]);
    reducedOutput.push(outputArray[0]);

    for (let i = step; i < timeArray.length - step; i += step) {
      let selectedIndex = i;
      let maxVariation = 0;

      const windowStart = Math.max(0, i - Math.floor(step/2));
      const windowEnd = Math.min(timeArray.length - 1, i + Math.floor(step/2));

      for (let j = windowStart; j <= windowEnd; j++) {
        if (j > 0 && j < outputArray.length - 1) {
          const variation = Math.abs(outputArray[j+1] - 2*outputArray[j] + outputArray[j-1]);
          if (variation > maxVariation) {
            maxVariation = variation;
            selectedIndex = j;
          }
        }
      }

      reducedTime.push(timeArray[selectedIndex]);
      reducedOutput.push(outputArray[selectedIndex]);
    }

    reducedTime.push(timeArray[timeArray.length - 1]);
    reducedOutput.push(outputArray[outputArray.length - 1]);

    return { time: reducedTime, output: reducedOutput };
  };

  const resetForm = () => {
    selectedFile.value = null;
    inputVoltage.value = 5;
    analysisComment.value = '';
    voltageError.value = '';
    error.value = '';
    results.value = null;
    uploadedDocumentId.value = null;
    processingMessage.value = '';
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
      text-align: center;

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

.comment-field {
  width: 100%;
  max-width: 500px;

  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: $text-color-light-mode;

    .dark-mode & {
      color: $text-color-dark-mode;
    }
  }

  .comment-textarea {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-family: inherit;
    font-size: 0.9rem;
    resize: vertical;
    min-height: 80px;

    &:focus {
      outline: none;
      border-color: $primary-color-light-mode;
      box-shadow: 0 0 0 2px rgba($primary-color-light-mode, 0.2);
    }

    .dark-mode & {
      background-color: rgba(255, 255, 255, 0.1);
      border-color: #555;
      color: $text-color-dark-mode;

      &:focus {
        border-color: $primary-color-dark-mode;
        box-shadow: 0 0 0 2px rgba($primary-color-dark-mode, 0.2);
      }
    }
  }

  .char-count {
    display: block;
    text-align: right;
    margin-top: 0.25rem;
    font-size: 0.8rem;
    color: $secondary-color-light-mode;

    .dark-mode & {
      color: $secondary-color-dark-mode;
    }
  }
}
</style>
