<template>
  <div
    class="upload-area"
    :class="{ 'dragging': isDragging, 'has-file': hasFile }"
    @click="triggerFileInput"
    @dragleave.prevent="onDragLeave"
    @dragover.prevent="onDragOver"
    @drop.prevent="onDrop"
  >
    <div class="upload-icon">
      <span v-if="!hasFile" class="cloud-upload-icon"></span>
      <span v-else class="file-icon"></span>
    </div>

    <p class="upload-text">
      <template v-if="hasFile">
        {{ shortenFileName(selectedFile.name) }}
        <button class="remove-file-btn" @click.stop="removeFile">
          <span class="close-icon"></span>
        </button>
      </template>
      <template v-else>
        Cargue o arrastre algún archivo .csv
      </template>
    </p>

    <input
      ref="fileInput"
      accept=".csv"
      class="file-input"
      type="file"
      @change="onFileSelected"
    />
  </div>

  <div v-if="error" class="error-message">
    {{ error }}
  </div>
</template>

<script setup>
  import { computed, ref } from 'vue';

  const props = defineProps({
    modelValue: {
      type: Object,
      default: null,
    },
  });

  const emit = defineEmits(['update:modelValue', 'error']);

  const fileInput = ref(null);
  const isDragging = ref(false);
  const error = ref('');
  const selectedFile = computed(() => props.modelValue);
  const hasFile = computed(() => selectedFile.value !== null);

  const isValidFile = file => {
    if (!file) return false;

    const fileName = file.name.toLowerCase();
    const isCSV = fileName.endsWith('.csv');

    if (!isCSV) {
      error.value = 'Solo se permiten archivos CSV (.csv)';
      emit('error', error.value);
      return false;
    }

    error.value = '';
    return true;
  };

  const triggerFileInput = () => {
    if (!hasFile.value) {
      fileInput.value.click();
    }
  };

  const onDragOver = () => {
    isDragging.value = true;
  };

  const onDragLeave = () => {
    isDragging.value = false;
  };

  const onDrop = event => {
    isDragging.value = false;
    if (event.dataTransfer.files.length) {
      const file = event.dataTransfer.files[0];
      if (isValidFile(file)) {
        emit('update:modelValue', file);
      }
    }
  };

  const onFileSelected = event => {
    if (event.target.files.length) {
      const file = event.target.files[0];
      if (isValidFile(file)) {
        emit('update:modelValue', file);
      } else {
        if (fileInput.value) {
          fileInput.value.value = '';
        }
      }
    }
  };

  const removeFile = () => {
    emit('update:modelValue', null);
    error.value = '';
    if (fileInput.value) {
      fileInput.value.value = '';
    }
  };

  const shortenFileName = fileName => {
    if (!fileName) return '';
    if (fileName.length <= 20) return fileName;

    const extension = fileName.split('.').pop();
    const name = fileName.substring(0, fileName.length - extension.length - 1);

    if (name.length <= 16) return fileName;

    return name.substring(0, 16) + '...' + (extension ? '.' + extension : '');
  };
</script>

<style lang="scss">
@import '../styles/variables.scss';

.upload-area {
  width: 100%;
  max-width: 350px;
  height: 180px;
  border: 2px dashed #ccc;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
  margin: 0 auto;

  &:hover {
    border-color: $primary-color-light-mode;

    .dark-mode & {
      border-color: $primary-color-dark-mode;
    }
  }

  &.dragging {
    border-color: $primary-color-light-mode;
    background-color: rgba(20, 112, 175, 0.05);

    .dark-mode & {
      border-color: $primary-color-dark-mode;
      background-color: rgba(183, 226, 255, 0.05);
    }
  }

  &.has-file {
    border-color: $primary-color-light-mode;
    background-color: rgba(20, 112, 175, 0.05);

    .dark-mode & {
      border-color: $primary-color-dark-mode;
      background-color: rgba(183, 226, 255, 0.05);
    }
  }

  .upload-icon {
    color: $primary-color-light-mode;
    margin-bottom: 0.5rem;

    .dark-mode & {
      color: $primary-color-dark-mode;
    }

    .cloud-upload-icon {
      display: inline-block;
      width: 50px;
      height: 50px;
      background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%231470AF' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M8 17l4-4 4 4'/%3E%3Cpath d='M12 12v9'/%3E%3Cpath d='M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3'/%3E%3C/svg%3E");
      background-size: contain;
      background-repeat: no-repeat;
      background-position: center;

      .dark-mode & {
        background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23B7E2FF' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M8 17l4-4 4 4'/%3E%3Cpath d='M12 12v9'/%3E%3Cpath d='M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3'/%3E%3C/svg%3E");
      }
    }

    .file-icon {
      display: inline-block;
      width: 40px;
      height: 40px;
      background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%231470AF' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z'/%3E%3Cpolyline points='14 2 14 8 20 8'/%3E%3Cline x1='16' y1='13' x2='8' y2='13'/%3E%3Cline x1='16' y1='17' x2='8' y2='17'/%3E%3Cpolyline points='10 9 9 9 8 9'/%3E%3C/svg%3E");
      background-size: contain;
      background-repeat: no-repeat;
      background-position: center;

      .dark-mode & {
        background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='%23B7E2FF' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpath d='M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z'/%3E%3Cpolyline points='14 2 14 8 20 8'/%3E%3Cline x1='16' y1='13' x2='8' y2='13'/%3E%3Cline x1='16' y1='17' x2='8' y2='17'/%3E%3Cpolyline points='10 9 9 9 8 9'/%3E%3C/svg%3E");
      }
    }
  }

  .upload-text {
    color: #666;
    text-align: center;

    .dark-mode & {
      color: #aaa;
    }

    .remove-file-btn {
      background: none;
      border: none;
      color: #e74c3c;
      cursor: pointer;
      font-size: 1rem;
      margin-left: 0.5rem;
      display: flex;
      align-items: center;

      .close-icon::before {
        content: "✕";
      }

      &:hover {
        color: darken(#e74c3c, 10%);
      }
    }
  }

  .file-input {
    position: absolute;
    width: 100%;
    height: 100%;
    opacity: 0;
    cursor: pointer;
    top: 0;
    left: 0;
  }
}

.error-message {
  margin-top: 0.5rem;
  color: #e74c3c;
  font-size: 0.85rem;
  text-align: center;

  .dark-mode & {
    color: #ff6b6b;
  }
}

@media (max-width: 576px) {
  .upload-area {
    height: 150px;
    max-width: 250px;

    .upload-icon {
      .cloud-upload-icon,
      .file-icon {
        width: 30px;
        height: 30px;
      }
    }

    .upload-text {
      font-size: 0.8rem;
    }
  }
}
</style>
