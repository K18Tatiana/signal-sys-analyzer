// stores/app.js
import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

export const useAppStore = defineStore('app', () => {
  // Estado
  const selectedFile = ref(null);
  const inputVoltage = ref(5);
  const isProcessing = ref(false);
  const processingResults = ref(null);
  const processingError = ref(null);

  // Getters
  const hasFile = computed(() => selectedFile.value !== null);
  const isFormValid = computed(() => inputVoltage.value > 0 && hasFile.value);
  const hasResults = computed(() => processingResults.value !== null);
  const hasError = computed(() => processingError.value !== null);

  // Acciones
  function setFile (file) {
    selectedFile.value = file;
    // Resetear resultados anteriores al cambiar el archivo
    processingResults.value = null;
    processingError.value = null;
  }

  function setVoltage (voltage) {
    inputVoltage.value = parseFloat(voltage);
  }

  async function processData () {
    if (!isFormValid.value) return;

    isProcessing.value = true;
    processingError.value = null;

    try {
      // Simulación de procesamiento de datos
      // En un caso real, aquí se haría una llamada a una API
      await new Promise(resolve => setTimeout(resolve, 2000));

      // Simulación de resultado exitoso
      processingResults.value = {
        type: 'First Order System',
        parameters: {
          timeConstant: 1.5,
          gain: 2.3,
        },
        response: {
          riseTime: 3.3,
          settlingTime: 7.5,
          steadyStateValue: inputVoltage.value * 2.3,
        },
      };
    } catch (error) {
      processingError.value = error.message || 'Ha ocurrido un error durante el procesamiento.';
    } finally {
      isProcessing.value = false;
    }
  }

  function reset () {
    selectedFile.value = null;
    inputVoltage.value = 5;
    processingResults.value = null;
    processingError.value = null;
  }

  return {
    // Estado
    selectedFile,
    inputVoltage,
    isProcessing,
    processingResults,
    processingError,

    // Getters
    hasFile,
    isFormValid,
    hasResults,
    hasError,

    // Acciones
    setFile,
    setVoltage,
    processData,
    reset,
  };
});
