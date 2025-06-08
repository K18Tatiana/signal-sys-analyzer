<template>
  <div class="documents-page">
    <div class="page-header">
      <h1 class="page-title">Mis Documentos</h1>
      <router-link class="btn btn-primary" to="/">
        <i class="icon-plus"></i>
        Nuevo Análisis
      </router-link>
    </div>

    <div v-if="loading" class="loading-indicator">
      <div class="spinner"></div>
      <p>Cargando documentos...</p>
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <div v-if="!loading && documents.length === 0" class="empty-state">
      <div class="empty-icon">📄</div>
      <h3>No tienes documentos aún</h3>
      <p>Sube tu primer archivo CSV para comenzar a analizar sistemas de control</p>
      <router-link class="btn btn-primary" to="/">
        Subir Primer Documento
      </router-link>
    </div>

    <div v-if="!loading && documents.length > 0" class="documents-grid">
      <div
        v-for="document in documents"
        :key="document.id"
        class="document-card"
      >
        <div class="document-header">
          <div class="document-info">
            <h3 class="document-name">{{ document.original_filename }}</h3>
            <p class="document-date">
              Subido: {{ formatDate(document.upload_date) }}
            </p>
            <p class="analysis-count">
              {{ document.analysis_count }} análisis realizados
            </p>
          </div>
          <div class="document-actions">
            <button
              class="btn btn-secondary btn-sm"
              @click="toggleDocumentDetails(document.id)"
            >
              {{ expandedDocuments.includes(document.id) ? 'Ocultar' : 'Ver Detalles' }}
            </button>
            <button
              class="btn btn-danger btn-sm"
              :disabled="deletingDocument === document.id"
              @click="deleteDocument(document.id)"
            >
              {{ deletingDocument === document.id ? 'Eliminando...' : 'Eliminar' }}
            </button>
          </div>
        </div>

        <!-- Detalles del documento (expandible) -->
        <div
          v-if="expandedDocuments.includes(document.id)"
          class="document-details"
        >
          <div class="new-analysis-section">
            <h4>Realizar Nuevo Análisis</h4>
            <div class="reanalysis-form">
              <div class="form-row">
                <div class="voltage-input">
                  <label>Voltaje de entrada:</label>
                  <input
                    v-model.number="newAnalysisVoltage[document.id]"
                    class="voltage-field"
                    max="12"
                    min="4"
                    placeholder="5.0"
                    step="0.1"
                    type="number"
                  />
                  <span class="unit">V</span>
                </div>
              </div>

              <div class="form-row">
                <div class="comment-input">
                  <label>Comentario (opcional):</label>
                  <textarea
                    v-model="newAnalysisComment[document.id]"
                    class="comment-textarea"
                    maxlength="500"
                    placeholder="Describe este análisis..."
                    rows="2"
                  ></textarea>
                  <small class="char-count">
                    {{ (newAnalysisComment[document.id] || '').length }}/500
                  </small>
                </div>
              </div>

              <div class="form-actions">
                <button
                  class="btn btn-primary"
                  :disabled="!isValidVoltage(newAnalysisVoltage[document.id]) || analyzingDocument === document.id"
                  @click="createNewAnalysis(document.id)"
                >
                  {{ analyzingDocument === document.id ? 'Analizando...' : 'Realizar Análisis' }}
                </button>
              </div>
            </div>
          </div>

          <div v-if="documentAnalyses[document.id]" class="analyses-section">
            <h4>Análisis Realizados</h4>
            <div
              v-for="analysis in documentAnalyses[document.id]"
              :key="analysis.analysis.id"
              class="analysis-item"
            >
              <div class="analysis-header" @click="toggleAnalysisExpansion(analysis.analysis.id)">
                <div class="analysis-info">
                  <strong>Voltaje: {{ analysis.analysis.input_voltage }}V</strong>
                  <br>
                  <span class="analysis-date">
                    {{ formatDate(analysis.analysis.created_at) }}
                  </span>
                  <p v-if="analysis.analysis.comment" class="analysis-comment">
                    <i class="comment-icon">💬</i>
                    <em>"{{ analysis.analysis.comment }}"</em>
                  </p>
                </div>
                <div class="analysis-actions">
                  <div class="analysis-status">
                    <span
                      :class="[
                        'status-badge',
                        analysis.analysis.is_processed ? 'status-completed' : 'status-processing'
                      ]"
                    >
                      <span v-if="!analysis.analysis.is_processed" class="processing-spinner"></span>
                      {{ analysis.analysis.is_processed ? 'Completado' : 'Procesando...' }}
                    </span>
                  </div>
                  <button class="expand-toggle">
                    {{ expandedAnalyses[analysis.analysis.id] ? '▼' : '▶' }}
                  </button>
                </div>
              </div>

              <!-- Resultados del análisis (expandibles) -->
              <div v-if="analysis.result && expandedAnalyses[analysis.analysis.id]" class="analysis-results">
                <ResultsDisplay
                  :compact="true"
                  :results="formatAnalysisResult(analysis)"
                  :show-actions="false"
                  :show-download="false"
                  :show-reset="false"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { onMounted, reactive, ref } from 'vue';
  import { authService } from '/src/services/auth.service.js';
  import ResultsDisplay from '/src/components/ResultsDisplay.vue';
  import Swal from 'sweetalert2';

  // Configuración api
  const API_URL = import.meta.env.VITE_API_URL

  // Estado
  const documents = ref([]);
  const loading = ref(false);
  const error = ref('');
  const expandedDocuments = ref([]);
  const documentAnalyses = reactive({});
  const deletingDocument = ref(null);
  const analyzingDocument = ref(null);

  // Estado para nuevos análisis
  const newAnalysisVoltage = reactive({});
  const newAnalysisComment = reactive({});

  const expandedAnalyses = reactive({});

  // Función para obtener headers de autorización
  const getAuthHeaders = () => {
    const headers = { 'Content-Type': 'application/json' };
    if (authService.isAuthenticated()) {
      const token = localStorage.getItem('token');
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }
    return headers;
  };

  // Cargar documentos del usuario
  const loadDocuments = async () => {
    if (!authService.isAuthenticated()) {
      error.value = 'Debes iniciar sesión para ver tus documentos';
      return;
    }

    loading.value = true;
    error.value = '';

    try {
      const response = await fetch(`${API_URL}/user/documents`, {
        method: 'GET',
        headers: getAuthHeaders(),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Error al cargar documentos');
      }

      const data = await response.json();
      documents.value = data;

      // Inicializar voltaje por defecto para nuevos análisis
      data.forEach(doc => {
        newAnalysisVoltage[doc.id] = 5.0;
        newAnalysisComment[doc.id] = '';
      });

    } catch (err) {
      error.value = err.message || 'Error al cargar documentos';
    } finally {
      loading.value = false;
    }
  };

  const toggleDocumentDetails = async documentId => {
    const index = expandedDocuments.value.indexOf(documentId);

    if (index > -1) {
      expandedDocuments.value.splice(index, 1);
    } else {
      expandedDocuments.value.push(documentId);
      if (!documentAnalyses[documentId]) {
        await loadDocumentAnalyses(documentId);
      }
    }
  };

  // Cargar análisis de un documento específico
  const loadDocumentAnalyses = async documentId => {
    try {
      const response = await fetch(`${API_URL}/user/documents/${documentId}`, {
        method: 'GET',
        headers: getAuthHeaders(),
      });

      if (!response.ok) {
        throw new Error('Error al cargar análisis del documento');
      }

      const data = await response.json();
      documentAnalyses[documentId] = data.analyses;

    } catch (err) {
      console.error('Error cargando análisis:', err);
    }
  };

  const toggleAnalysisExpansion = analysisId => {
    expandedAnalyses[analysisId] = !expandedAnalyses[analysisId];
  };

  const deleteDocument = async documentId => {
    const result = await Swal.fire({
      title: '¿Eliminar documento?',
      text: 'Esta acción no se puede deshacer',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#dc3545',
      cancelButtonColor: '#6c757d',
      confirmButtonText: 'Sí, eliminar',
      cancelButtonText: 'Cancelar',
      reverseButtons: true,
      customClass: {
        popup: 'sweet-alert-popup',
        title: 'sweet-alert-title',
        content: 'sweet-alert-content',
      },
    });

    if (!result.isConfirmed) {
      return;
    }

    deletingDocument.value = documentId;

    try {
      const response = await fetch(`${API_URL}/user/documents/${documentId}`, {
        method: 'DELETE',
        headers: getAuthHeaders(),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Error al eliminar documento');
      }

      // Remover documento de la lista
      documents.value = documents.value.filter(doc => doc.id !== documentId);

      // Limpiar estados relacionados
      delete documentAnalyses[documentId];
      delete newAnalysisVoltage[documentId];
      delete newAnalysisComment[documentId];
      Object.keys(expandedAnalyses).forEach(key => {
        if (key.toString().startsWith(documentId.toString())) {
          delete expandedAnalyses[key];
        }
      });
      const expandedIndex = expandedDocuments.value.indexOf(documentId);
      if (expandedIndex > -1) {
        expandedDocuments.value.splice(expandedIndex, 1);
      }

      // Alert de éxito
      await Swal.fire({
        title: '¡Eliminado!',
        text: 'El documento ha sido eliminado correctamente',
        icon: 'success',
        timer: 2000,
        showConfirmButton: false,
        toast: true,
        position: 'top-end',
      });

    } catch (err) {
      error.value = err.message || 'Error al eliminar documento';

      // Alert de error
      await Swal.fire({
        title: 'Error',
        text: err.message || 'Error al eliminar documento',
        icon: 'error',
        confirmButtonText: 'Entendido',
        confirmButtonColor: '#dc3545',
      });
    } finally {
      deletingDocument.value = null;
    }
  };

  // Crear nuevo análisis
  const createNewAnalysis = async documentId => {
    const voltage = newAnalysisVoltage[documentId];
    const comment = newAnalysisComment[documentId];

    if (!isValidVoltage(voltage)) {
      error.value = 'El voltaje debe estar entre 4 y 12 V';
      return;
    }

    analyzingDocument.value = documentId;

    try {
      const requestData = {
        document_id: documentId,
        input_voltage: voltage,
      };

      if (comment && comment.trim()) {
        requestData.comment = comment.trim();
      }

      const response = await fetch(`${API_URL}/analysis`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify(requestData),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Error al crear análisis');
      }

      const analysisResponse = await response.json();

      newAnalysisVoltage[documentId] = 5.0;
      newAnalysisComment[documentId] = '';

      await loadDocumentAnalyses(documentId);

      const docIndex = documents.value.findIndex(doc => doc.id === documentId);
      if (docIndex > -1) {
        documents.value[docIndex].analysis_count += 1;
      }

      await loadDocumentAnalyses(documentId);
      await pollAnalysisCompletion(analysisResponse.id, documentId);

    } catch (err) {
      error.value = err.message || 'Error al crear análisis';
    } finally {
      analyzingDocument.value = null;
    }
  };

  const pollAnalysisCompletion = async (analysisId, documentId, maxAttempts = 15) => {
    for (let attempt = 0; attempt < maxAttempts; attempt++) {
      try {
        await new Promise(resolve => setTimeout(resolve, 2000));

        const response = await fetch(`${API_URL}/analysis/${analysisId}`, {
          method: 'GET',
          headers: getAuthHeaders(),
        });

        if (!response.ok) {
          console.warn(`Intento ${attempt + 1}: Error al verificar análisis`);
          continue;
        }

        const data = await response.json();

        if (data.analysis_request && data.analysis_request.is_processed === true) {
          await loadDocumentAnalyses(documentId);
          return;
        }

        if (data.result) {
          await loadDocumentAnalyses(documentId);
          return;
        }

      } catch (err) {
        console.warn(`Intento ${attempt + 1}: Error en polling:`, err);

        if (attempt === maxAttempts - 1) {
          error.value = 'El análisis está tomando más tiempo del esperado. Recarga la página para ver si se completó.';
          return;
        }
      }
    }

    console.warn('Análisis no completado en el tiempo esperado');
    error.value = 'El análisis está tomando más tiempo del esperado. Recarga la página para ver si se completó.';
  };

  const isValidVoltage = voltage => {
    return voltage >= 4 && voltage <= 12;
  };

  const formatDate = dateString => {
    return new Date(dateString).toLocaleString('es-ES', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatAnalysisResult = analysis => {
    if (!analysis.result) return null;

    const resultData = analysis.result;

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

    return {
      type: resultData.system_type || 'Sistema no identificado',
      description: resultData.description || getSystemDescription(resultData.system_type), // CAMBIO AQUÍ
      poles,
      graphData: processGraphData(resultData.graph_data),
      technicalSummary: resultData.technical_summary || null, // NUEVO CAMPO
      rawData: resultData.raw_data,
    };
  };

  const getSystemDescription = systemType => {
    const descriptions = {
      'subamortiguado': 'Sistema subamortiguado - presenta oscilaciones antes de estabilizarse.',
      'sobreamortiguado': 'Sistema sobreamortiguado - se aproxima gradualmente al valor final sin oscilaciones.',
    };
    return descriptions[systemType] || 'Sistema identificado con características específicas según el análisis realizado.';
  };

  const processGraphData = graphDataString => {
    if (!graphDataString) return null;

    try {
      const graphData = typeof graphDataString === 'string' ? JSON.parse(graphDataString) : graphDataString;

      if (graphData.time && graphData.output) {
        return {
          labels: graphData.time.map(t => t.toFixed(3)),
          datasets: [
            {
              label: 'Respuesta del Sistema',
              data: graphData.output,
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
    return null;
  };

  onMounted(() => {
    loadDocuments();
  });
</script>

<style lang="scss" scoped>
@import '/src/styles/variables.scss';

.documents-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;

  .page-title {
    font-size: 2rem;
    font-weight: 700;
    color: $primary-color-light-mode;
    margin: 0;

    .dark-mode & {
      color: $primary-color-dark-mode;
    }
  }
}

.loading-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 3rem;

  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid rgba($primary-color-light-mode, 0.2);
    border-radius: 50%;
    border-top-color: $primary-color-light-mode;
    animation: spin 1s ease-in-out infinite;
    margin-bottom: 1rem;
  }
}

.error-message {
  background-color: #fff3f3;
  border-left: 4px solid #e74c3c;
  padding: 1rem;
  margin-bottom: 1.5rem;
  border-radius: 4px;
  color: #e74c3c;
}

.empty-state {
  text-align: center;
  padding: 4rem 2rem;

  .empty-icon {
    font-size: 4rem;
    margin-bottom: 1rem;
  }

  h3 {
    color: $secondary-color-light-mode;
    margin-bottom: 0.5rem;

    .dark-mode & {
      color: $secondary-color-dark-mode;
    }
  }

  p {
    color: $text-color-light-mode;
    margin-bottom: 2rem;

    .dark-mode & {
      color: $text-color-dark-mode;
    }
  }
}

.documents-grid {
  display: grid;
  gap: 1.5rem;
}

.document-card {
  background-color: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  border: 1px solid #e0e0e0;

  .dark-mode & {
    background-color: $background-color-dark-mode;
    border-color: #404040;
  }
}

.document-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
  gap: 20px;

  .document-info {
    flex: 1;

    .document-name {
      font-size: 1.2rem;
      font-weight: 600;
      margin: 0 0 0.5rem 0;
      color: $primary-color-light-mode;

      .dark-mode & {
        color: $primary-color-dark-mode;
      }
    }

    .document-date,
    .analysis-count {
      font-size: 0.9rem;
      color: $secondary-color-light-mode;
      margin: 0.25rem 0;

      .dark-mode & {
        color: $secondary-color-dark-mode;
      }
    }
  }

  .document-actions {
    display: flex;
    gap: 0.5rem;
    flex-shrink: 0;
  }
}

.document-details {
  border-top: 1px solid #e0e0e0;
  padding-top: 1rem;

  .dark-mode & {
    border-color: #404040;
  }

  h4 {
    font-size: 1.1rem;
    font-weight: 600;
    margin: 0 0 1rem 0;
    color: $primary-color-light-mode;

    .dark-mode & {
      color: $primary-color-dark-mode;
    }
  }
}

.new-analysis-section {
  margin-bottom: 2rem;
  padding: 1.5rem;
  background-color: #f0f8ff;
  border-radius: 8px;
  border: 1px solid #e3f2fd;

  .dark-mode & {
    background-color: rgba(20, 112, 175, 0.1);
    border-color: rgba(20, 112, 175, 0.2);
  }

  .reanalysis-form {
    .form-row {
      margin-bottom: 1rem;

      &:last-child {
        margin-bottom: 0;
      }
    }

    .voltage-input {
      display: flex;
      align-items: center;
      gap: 0.5rem;

      label {
        font-weight: 500;
        color: $text-color-light-mode;
        min-width: 140px;

        .dark-mode & {
          color: $text-color-dark-mode;
        }
      }

      .voltage-field {
        width: 100px;
        padding: 0.5rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-size: 0.9rem;

        .dark-mode & {
          background-color: rgba(255, 255, 255, 0.1);
          border-color: #555;
          color: $text-color-dark-mode;
        }
      }

      .unit {
        font-weight: 500;
        color: $secondary-color-light-mode;

        .dark-mode & {
          color: $secondary-color-dark-mode;
        }
      }
    }

    .comment-input {
      label {
        display: block;
        font-weight: 500;
        color: $text-color-light-mode;
        margin-bottom: 0.5rem;

        .dark-mode & {
          color: $text-color-dark-mode;
        }
      }

      .comment-textarea {
        width: 100%;
        padding: 0.5rem;
        border: 1px solid #ddd;
        border-radius: 4px;
        font-family: inherit;
        font-size: 0.9rem;
        resize: vertical;

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
        font-size: 0.75rem;
        color: $secondary-color-light-mode;

        .dark-mode & {
          color: $secondary-color-dark-mode;
        }
      }
    }

    .form-actions {
      display: flex;
      justify-content: flex-end;
      padding-top: 0.5rem;
    }
  }
}

.analyses-section {
  .analysis-item {
    border: 1px solid #e0e0e0;
    border-radius: 4px;
    padding: 1rem;
    margin-bottom: 1rem;

    .dark-mode & {
      border-color: #404040;
    }

    .analysis-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.5rem;
      cursor: pointer; // AGREGAR esta línea
      transition: background-color 0.2s; // AGREGAR esta línea

      &:hover {
        background-color: rgba(0, 0, 0, 0.02);
        border-radius: 4px;

        .dark-mode & {
          background-color: rgba(255, 255, 255, 0.02);
        }
      }

      .analysis-info {
        flex: 1;

        .analysis-date {
          font-size: 0.85rem;
          color: $secondary-color-light-mode;

          .dark-mode & {
            color: $secondary-color-dark-mode;
          }
        }

        .analysis-comment {
          margin: 0.5rem 0 0 0;
          padding: 0.5rem;
          background-color: rgba(255, 255, 255, 0.5);
          border-radius: 4px;
          font-size: 0.9rem;
          color: $text-color-light-mode;

          .dark-mode & {
            background-color: rgba(255, 255, 255, 0.05);
            color: $text-color-dark-mode;
          }

          .comment-icon {
            margin-right: 0.5rem;
          }

          em {
            color: $secondary-color-light-mode;

            .dark-mode & {
              color: $secondary-color-dark-mode;
            }
          }
        }
      }

      .analysis-actions {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        flex-shrink: 0;
      }

      .expand-toggle {
        background: none;
        border: none;
        font-size: 0.9rem;
        color: $secondary-color-light-mode;
        cursor: pointer;
        padding: 0.25rem;
        border-radius: 2px;
        transition: all 0.2s;

        &:hover {
          background-color: rgba(0, 0, 0, 0.1);
          color: $primary-color-light-mode;
        }

        .dark-mode & {
          color: $secondary-color-dark-mode;

          &:hover {
            background-color: rgba(255, 255, 255, 0.1);
            color: $primary-color-dark-mode;
          }
        }
      }

      .status-badge {
        padding: 0.25rem 0.5rem;
        border-radius: 4px;
        font-size: 0.8rem;
        font-weight: 500;
        display: flex;
        align-items: center;
        gap: 0.5rem;

        &.status-completed {
          background-color: #d4edda;
          color: #155724;
        }

        &.status-processing {
          background-color: #fff3cd;
          color: #856404;

          .processing-spinner {
            width: 12px;
            height: 12px;
            border: 2px solid transparent;
            border-top: 2px solid #856404;
            border-radius: 50%;
            animation: spin 1s linear infinite;
          }
        }
      }
    }

    .analysis-results {
      margin-top: 1rem;
    }
  }
}

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.2s;

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  &.btn-sm {
    padding: 0.25rem 0.5rem;
    font-size: 0.875rem;
  }

  &.btn-primary {
    background-color: $primary-color-light-mode;
    color: white;

    &:hover:not(:disabled) {
      background-color: darken($primary-color-light-mode, 10%);
    }

    .dark-mode & {
      background-color: $primary-color-dark-mode;
      color: $background-color-dark-mode;
    }
  }

  &.btn-secondary {
    background-color: #6c757d;
    color: white;

    &:hover:not(:disabled) {
      background-color: darken(#6c757d, 10%);
    }
  }

  &.btn-danger {
    background-color: #dc3545;
    color: white;

    &:hover:not(:disabled) {
      background-color: darken(#dc3545, 10%);
    }
  }

  &.btn-link {
    background: none;
    color: $primary-color-light-mode;
    padding: 0.25rem;

    &:hover {
      text-decoration: underline;
    }

    .dark-mode & {
      color: $primary-color-dark-mode;
    }
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 768px) {
  .documents-page {
    padding: 1rem 0.5rem;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;

    .page-title {
      font-size: 1.5rem;
    }
  }

  .document-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;

    .document-actions {
      width: 100%;
      flex-wrap: wrap;
    }
  }

  .voltage-input {
    flex-wrap: wrap;
  }

  .analysis-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
}

:deep(.sweet-alert-popup) {
  font-family: inherit;
}

:deep(.sweet-alert-title) {
  color: $primary-color-light-mode;

  .dark-mode & {
    color: $primary-color-dark-mode;
  }
}

:deep(.sweet-alert-content) {
  color: $text-color-light-mode;

  .dark-mode & {
    color: $text-color-dark-mode;
  }
}
</style>
