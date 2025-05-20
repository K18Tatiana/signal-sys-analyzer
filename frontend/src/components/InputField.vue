<template>
  <div class="input-field">
    <label v-if="label" class="input-label" :for="id">{{ label }}</label>

    <div class="input-wrapper" :class="{ 'with-unit': unit }">
      <input
        :id="id"
        class="input-control"
        :class="{ 'error-border': error }"
        :disabled="disabled"
        :max="max"
        :min="min"
        :placeholder="placeholder"
        :required="required"
        :step="step"
        :type="type"
        :value="modelValue"
        @input="updateValue"
      />
      <span v-if="unit" class="input-unit">{{ unit }}</span>
    </div>

    <div v-if="error" class="input-error">{{ error }}</div>
    <div v-if="hint" class="input-hint">{{ hint }}</div>
  </div>
</template>

<script setup>
  const props = defineProps({
    modelValue: {
      type: [String, Number],
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    id: {
      type: String,
      default: () => `input-${Math.random().toString(36).substr(2, 9)}`,
    },
    type: {
      type: String,
      default: 'text',
    },
    placeholder: {
      type: String,
      default: '',
    },
    unit: {
      type: String,
      default: '',
    },
    min: {
      type: [Number, String],
      default: null,
    },
    max: {
      type: [Number, String],
      default: null,
    },
    step: {
      type: [Number, String],
      default: null,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    error: {
      type: String,
      default: '',
    },
    hint: {
      type: String,
      default: '',
    },
    minValue: {
      type: [Number, String],
      default: 4,
    },
    maxValue: {
      type: [Number, String],
      default: 12,
    },
    showRangeAlert: {
      type: Boolean,
      default: true,
    },
  });

  const emit = defineEmits(['update:modelValue', 'validation']);

  const updateValue = event => {
    let value = event.target.value;

    if (props.type === 'number') {
      value = value === '' ? '' : Number(value);

      if (value !== '' && props.minValue !== null && props.maxValue !== null) {
        const minVal = Number(props.minValue);
        const maxVal = Number(props.maxValue);

        if (value < minVal || value > maxVal) {
          emit('validation', `El valor debe estar entre ${minVal} y ${maxVal} ${props.unit || ''}`);
        } else {
          emit('validation', '');
        }
      }
    }

    emit('update:modelValue', value);
  };
</script>

<style lang="scss">
@import '../styles/variables.scss';

.input-field {
  width: 100%;
  margin-bottom: 1rem;
  text-align: center;

  .input-label {
    display: block;
    margin-bottom: 0.5rem;
    text-align: center;
  }

  .input-wrapper {
    position: relative;
    display: inline-flex;
    width: auto;

    &.with-unit {
      .input-control {
        border-radius: 4px 0 0 4px;
        width: 60px;
        text-align: center;

        &.error-border {
          border-color: #e74c3c;

          &:focus {
            box-shadow: 0 0 0 0.2rem rgba(231, 76, 60, 0.25);
          }
        }
      }

      .input-unit {
        display: flex;
        align-items: center;
        padding: 0 0.75rem;
        background-color: #f5f5f5;
        border: 1px solid #ccc;
        border-left: none;
        border-radius: 0 4px 4px 0;
        font-size: 0.9rem;

        .dark-mode & {
          background-color: #444;
          border-color: #555;
          color: #fff;
        }
      }
    }
  }

  .input-control {
    padding: 0.5rem 0.75rem;
    font-size: 1rem;
    line-height: 1.5;
    color: $text-color-light-mode;
    background-color: #fff;
    border: 1px solid #ccc;
    border-radius: 4px;
    transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;

    &:focus {
      outline: 0;
      border-color: $primary-color-light-mode;
      box-shadow: 0 0 0 0.2rem rgba($primary-color-light-mode, 0.25);

      .dark-mode & {
        border-color: $primary-color-dark-mode;
        box-shadow: 0 0 0 0.2rem rgba($primary-color-dark-mode, 0.25);
      }
    }

    &:disabled {
      background-color: #e9ecef;
      opacity: 1;

      .dark-mode & {
        background-color: #2a3238;
      }
    }

    .dark-mode & {
      background-color: #2a3238;
      border-color: #555;
      color: $text-color-dark-mode;
    }

    // Eliminar los botones incremento/decremento en campos numéricos
    &[type="number"] {
      -moz-appearance: textfield;
    }

    &[type="number"]::-webkit-outer-spin-button,
    &[type="number"]::-webkit-inner-spin-button {
      -webkit-appearance: none;
      margin: 0;
    }

    &.error-border {
      border-color: #e74c3c;

      &:focus {
        box-shadow: 0 0 0 0.2rem rgba(231, 76, 60, 0.25);
      }

      .dark-mode & {
        border-color: #ff6b6b;

        &:focus {
          box-shadow: 0 0 0 0.2rem rgba(255, 107, 107, 0.25);
        }
      }
    }
  }

  .input-error {
    margin-top: 0.25rem;
    font-size: 0.8rem;
    color: #e74c3c;
    text-align: center;

    .dark-mode & {
      color: #ff6b6b;
    }
  }

  .range-alert {
    margin-top: 0.5rem;
    font-size: 0.8rem;
    color: #3498db;
    background-color: rgba(52, 152, 219, 0.1);
    border: 1px solid rgba(52, 152, 219, 0.3);
    border-radius: 4px;
    padding: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;

    .dark-mode & {
      color: #5dade2;
      background-color: rgba(93, 173, 226, 0.1);
      border-color: rgba(93, 173, 226, 0.3);
    }

    .info-icon {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 16px;
      height: 16px;
      border-radius: 50%;
      background-color: #3498db;
      color: white;
      font-size: 10px;
      font-weight: bold;
      margin-right: 6px;

      &::before {
        content: "i";
      }

      .dark-mode & {
        background-color: #5dade2;
      }
    }
  }

  .input-hint {
    margin-top: 0.25rem;
    font-size: 0.8rem;
    color: #666;
    text-align: center;

    .dark-mode & {
      color: #aaa;
    }
  }
}

@media (max-width: 768px) {
  .input-field {
    .input-wrapper {
      &.with-unit {
        .input-control {
          width: 50px;
        }
      }
    }
  }
}
</style>
