<template>
  <div v-if="results" class="results-display">
    <h2 class="results-title">RESULTADOS</h2>

    <div class="results-flex-container">
      <div class="collapsibles-column">
        <div class="collapsible-section">
          <div class="collapsible-header" @click="toggleSystemType">
            <h3>Tipo de sistema</h3>
            <span class="toggle-icon">{{ isSystemTypeOpen ? '▼' : '▶' }}</span>
          </div>
          <div v-show="isSystemTypeOpen" class="collapsible-content">
            <p class="system-name">Sistema {{ results.type }}</p>
            <p class="system-description">{{ results.description }}</p>
          </div>
        </div>

        <div class="collapsible-section">
          <div class="collapsible-header" @click="toggleSystemPoles">
            <h3>Polos del sistema</h3>
            <span class="toggle-icon">{{ isSystemPolesOpen ? '▼' : '▶' }}</span>
          </div>
          <div v-show="isSystemPolesOpen" class="collapsible-content">
            <div v-if="results.poles && results.poles.length > 0">
              <div v-for="(pole, index) in results.poles" :key="`pole-${index}`" class="pole-item">
                <p>{{ pole }}</p>
              </div>
            </div>
            <p v-else>No se han encontrado polos para este sistema.</p>
          </div>
        </div>
      </div>

      <div class="graph-column">
        <h3>Gráfica del Sistema</h3>
        <div class="system-graph">
          <img v-if="results.graphUrl" alt="Gráfica del sistema" class="graph-image" :src="results.graphUrl" />
          <div v-else-if="results.graphData" class="graph-container">
            <SystemChart
              ref="systemChart"
              :graph-data="results.graphData"
              :is-dark-mode="isPrinting ? false : isDarkMode"
            />
          </div>
          <div v-else class="graph-placeholder">
            <p>No hay datos disponibles para generar la gráfica.</p>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showActions" class="results-actions">
      <button v-if="showDownload" class="btn btn-primary" @click="downloadResults">
        <i class="icon bx bx-download"></i> Descargar Resultados
      </button>
      <button v-if="showReset" class="btn btn-secondary" @click="$emit('reset')">
        <i class="icon bx bx-repost"></i> Nuevo Análisis
      </button>
    </div>
  </div>
</template>

<script setup>
  import { inject, ref } from 'vue';
  import SystemChart from './SystemChart.vue';

  const props = defineProps({
    results: {
      type: Object,
      default: null,
    },
    showActions: {
      type: Boolean,
      default: true,
    },
    showDownload: {
      type: Boolean,
      default: true,
    },
    showReset: {
      type: Boolean,
      default: true,
    },
  });

  const isDarkMode = inject('isDarkMode');
  const isSystemTypeOpen = ref(true);
  const isSystemPolesOpen = ref(false);
  const isPrinting = ref(false);

  const toggleSystemType = () => {
    isSystemTypeOpen.value = !isSystemTypeOpen.value;
  };
  const toggleSystemPoles = () => {
    isSystemPolesOpen.value = !isSystemPolesOpen.value;
  };
  const downloadResults = async () => {
    if (!props.results) return;

    try {
      isPrinting.value = true;

      const systemChart = ref(null);

      isSystemTypeOpen.value = true;
      isSystemPolesOpen.value = true;

      await new Promise(resolve => setTimeout(resolve, 300));

      if (props.results.graphData && systemChart.value && systemChart.value.$el) {
        const svgElement = systemChart.value.$el.querySelector('svg');
        if (svgElement) {
          const svgClone = svgElement.cloneNode(true);
          svgClone.setAttribute('width', '100%');
          svgClone.setAttribute('height', 'auto');
          svgClone.style.maxHeight = '300px';

          const originalSvg = svgElement.outerHTML;

          svgElement.replaceWith(svgClone);

          setTimeout(() => {
            svgClone.outerHTML = originalSvg;
          }, 1000);
        }
      }

      window.print();

      setTimeout(() => {
        isPrinting.value = false;
      }, 1000);

    } catch (error) {
      console.error('Error al generar el PDF:', error);
      alert('Hubo un error al generar el PDF. Por favor, intenta nuevamente.');
      isPrinting.value = false;
    }
  }
</script>

<style lang="scss">
@import '../styles/variables.scss';

.results-display {
  width: 100%;
  max-width: 100%;
  margin: 2rem 0;
  animation: fadeIn 0.5s ease-in-out;

  .results-title {
    font-family: $font-primary;
    font-size: 1.5rem;
    margin-bottom: 1.5rem;
    color: $primary-color-light-mode;
    text-align: center;

    .dark-mode & {
      color: $primary-color-dark-mode;
    }
  }

  .results-flex-container {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    width: 100%;

    @media (min-width: 992px) {
      flex-direction: row;
      align-items: flex-start;
    }

    .collapsibles-column {
      width: 100%;

      @media (min-width: 992px) {
        margin-top: 40px;
        width: 50%;
        padding-right: 1.5rem;
      }
    }

    .graph-column {
      width: 100%;

      @media (min-width: 992px) {
        width: 50%;
      }
    }
  }

  .collapsible-section {
    margin-bottom: 1rem;
    border: 1px solid #eee;
    border-radius: 4px;
    overflow: hidden;

    .dark-mode & {
      border-color: #444;
    }

    .collapsible-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1rem;
      background-color: white;
      cursor: pointer;

      .dark-mode & {
        background-color: #2a3238;
      }

      h3 {
        margin: 0;
        font-size: 1rem;
        color: $secondary-color-light-mode;

        .dark-mode & {
          color: $secondary-color-dark-mode;
        }
      }

      .toggle-icon {
        font-size: 0.8rem;
        color: $secondary-color-light-mode;

        .dark-mode & {
          color: $secondary-color-dark-mode;
        }
      }

      &:hover {
        background-color: #f9f9f9;

        .dark-mode & {
          background-color: #343e46;
        }
      }
    }

    .collapsible-content {
      padding: 1rem;
      background-color: #f9f9f9;

      .dark-mode & {
        background-color: #343e46;
      }

      .system-name {
        font-weight: 600;
        margin-bottom: 0.5rem;
      }

      .system-description {
        margin-bottom: 0;
        line-height: 1.5;
      }

      .pole-item {
        margin-bottom: 0.5rem;
        padding: 0.5rem;
        background-color: white;
        border-radius: 4px;

        .dark-mode & {
          background-color: #2a3238;
        }

        &:last-child {
          margin-bottom: 0;
        }
      }
    }
  }

  .graph-column {
    h3 {
      font-size: 1rem;
      margin-bottom: 1rem;
      color: $secondary-color-light-mode;

      .dark-mode & {
        color: $secondary-color-dark-mode;
      }
    }

    .system-graph {
      border: 1px solid #eee;
      border-radius: 4px;
      overflow: hidden;
      background-color: white;

      .dark-mode & {
        background-color: #2a3238;
        border-color: #444;
      }

      .graph-image {
        width: 100%;
        height: auto;
        display: block;
      }

      .graph-container {
        width: 100%;
        height: 300px;
      }

      .graph-placeholder {
        height: 300px;
        display: flex;
        align-items: center;
        justify-content: center;
        background-color: #f9f9f9;

        .dark-mode & {
          background-color: #343e46;
        }

        p {
          color: #888;

          .dark-mode & {
            color: #aaa;
          }
        }
      }
    }
  }

  .results-actions {
    display: flex;
    justify-content: center;
    gap: 1rem;
    margin-top: 2rem;

    @media (max-width: 576px) {
      flex-direction: column;
      align-items: center;

      button {
        width: 100%;
        max-width: 250px;
      }
    }

    button {
      display: flex;
      align-items: center;
      justify-content: center;

      .icon {
        margin-right: 0.5rem;
        font-size: 20px;
      }
    }
  }
}

@media (max-width: 768px) {
  .results-display {
    .results-title {
      font-size: 1.2rem;
    }

    .collapsible-header {
      h3 {
        font-size: 0.9rem;
      }
    }

    .graph-column {
      h3 {
        font-size: 0.9rem;
      }

      .system-graph {
        .graph-container,
        .graph-placeholder {
          height: 250px;
        }
      }
    }
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@media print {
  body {
    background-color: white;
  }

  .app-header {
    box-shadow: none !important;

    .mobile-menu-toggle {
      display: none;
    }
    .desktop-nav {
      display: none;
    }
    .mobile-menu {
      display: none;
    }
  }

  .page-title {
    display: none;
  }
  .upload-section {
    display: none !important;
  }
  .results-display {
    margin: 0;

    .results-title {
      color: $primary-color-light-mode !important;
    }

    .results-flex-container {
      flex-direction: column;

      h3 {
        color: $secondary-color-light-mode !important;
      }
      p {
        color: #000 !important;
      }
    }

    .collapsibles-column {
      width: 100% !important;

      .collapsible-section {
        border: none;

        .toggle-icon {
          display: none;
        }
      }
    }

    .graph-column {
      width: 100% !important;

      canvas {
        width: 100% !important;
        height: 100% !important;
      }
      .graph-container {
        height: 500px !important;
      }
    }

    .results-actions {
      display: none;
    }
  }

  .app-footer {
      display: none;
    }
}
</style>
