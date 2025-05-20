<template>
  <div class="feedback-page">
    <div class="feedback-header">
      <h1 class="page-title">Enviar Feedback</h1>
      <p class="subtitle">Tu opinión es fundamental para mejorar SignalSysAnalyzer. Cuéntanos tu experiencia.</p>
    </div>

    <div class="feedback-container">
      <form class="feedback-form" :class="{ 'is-submitting': isSubmitting }" @submit.prevent="submitFeedback">
        <div class="form-section rating-section">
          <h3>¿Cómo calificarías tu experiencia con SignalSysAnalyzer?</h3>
          <p class="section-description">Tu valoración nos ayuda a saber qué tan satisfecho estás con nuestra herramienta.</p>

          <div class="star-rating">
            <template v-for="star in 5" :key="star">
              <div
                class="star-container"
                @click="selectRating(star)"
                @mouseleave="hoverRating = 0"
                @mouseover="hoverRating = star"
              >
                <i
                  class="bx"
                  :class="[
                    (star <= (hoverRating || rating)) ? 'bxs-star filled' : 'bx-star',
                    { 'selected': star <= rating }
                  ]"
                ></i>
                <span class="star-tooltip" :class="{ 'active': hoverRating === star }">
                  {{ getRatingLabel(star) }}
                </span>
              </div>
            </template>
          </div>

          <div v-if="rating > 0" class="rating-label">
            <strong>Tu valoración:</strong> {{ getRatingLabel(rating) }} ({{ rating }}/5)
          </div>
        </div>

        <div class="form-section">
          <h3>¿Tienes algún comentario o sugerencia?</h3>
          <p class="section-description">Cuéntanos qué podemos mejorar o qué características te gustaría ver en futuras versiones.</p>

          <div class="textarea-wrapper">
            <textarea
              ref="textareaElement"
              v-model="feedbackText"
              :class="{ error: errors.feedbackText }"
              placeholder="Escribe aquí tus comentarios o sugerencias..."
              rows="6"
              @blur="setActiveTextarea(false)"
              @focus="setActiveTextarea(true)"
            ></textarea>

            <div v-if="isTextareaActive && !feedbackText.length" class="feedback-suggestions">
              <p class="suggestion-prompt">Algunas ideas para tu feedback:</p>
              <ul class="suggestion-list">
                <li @mousedown.prevent="insertSuggestion('Me gustaría que la herramienta incluyera...')">
                  <i class="bx bx-plus-circle"></i> Me gustaría que la herramienta incluyera...
                </li>
                <li @mousedown.prevent="insertSuggestion('Tuve problemas con...')">
                  <i class="bx bx-error-circle"></i> Tuve problemas con...
                </li>
                <li @mousedown.prevent="insertSuggestion('Lo que más me gusta es...')">
                  <i class="bx bx-like"></i> Lo que más me gusta es...
                </li>
                <li @mousedown.prevent="insertSuggestion('Sería útil si pudiera...')">
                  <i class="bx bx-bulb"></i> Sería útil si pudiera...
                </li>
              </ul>
            </div>
          </div>

          <div class="character-count" :class="{ 'limit-close': feedbackText.length > 450 }">
            {{ feedbackText.length }}/500 caracteres
          </div>

          <span v-if="errors.feedbackText" class="error-message">{{ errors.feedbackText }}</span>
        </div>

        <div class="form-section user-info-section">
          <h3>Información de contacto (opcional)</h3>
          <p class="section-description">Si deseas que te contactemos sobre tu feedback, puedes dejarnos tu correo.</p>

          <input
            v-model="userEmail"
            :class="{ error: errors.userEmail }"
            placeholder="tu@email.com (opcional)"
            type="email"
          />
          <span v-if="errors.userEmail" class="error-message">{{ errors.userEmail }}</span>
        </div>

        <div class="form-actions">
          <button class="submit-button" :disabled="isSubmitting" type="submit">
            <span v-if="!isSubmitting">Enviar Feedback</span>
            <span v-else class="submitting">
              <i class="bx bx-loader-alt bx-spin"></i>
              Enviando...
            </span>
          </button>
        </div>
      </form>

      <!-- Mensaje de éxito -->
      <div v-if="feedbackSubmitted" class="success-message">
        <div class="success-icon">
          <i class="bx bx-check"></i>
        </div>
        <h3>¡Gracias por tu feedback!</h3>
        <p>Tu opinión es muy valiosa para nosotros y nos ayuda a mejorar SignalSysAnalyzer.</p>
        <button class="reset-button" @click="resetForm">Enviar otro feedback</button>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { reactive, ref } from 'vue';

  const rating = ref(0);
  const hoverRating = ref(0);
  const feedbackText = ref('');
  const userEmail = ref('');
  const isSubmitting = ref(false);
  const feedbackSubmitted = ref(false);
  const isTextareaActive = ref(false);
  const textareaElement = ref(null);

  const errors = reactive({
    feedbackText: '',
    userEmail: '',
  });

  const getRatingLabel = star => {
    const labels = {
      1: 'Muy insatisfecho',
      2: 'Insatisfecho',
      3: 'Neutral',
      4: 'Satisfecho',
      5: 'Muy satisfecho',
    };
    return labels[star] || '';
  };

  const selectRating = star => {
    rating.value = star;
  };

  const setActiveTextarea = active => {
    isTextareaActive.value = active;
  };

  const insertSuggestion = suggestion => {
    feedbackText.value = suggestion;

    setTimeout(() => {
      if (textareaElement.value) {
        textareaElement.value.focus();
        const length = feedbackText.value.length;
        textareaElement.value.setSelectionRange(length, length);
      }
    }, 10);

    isTextareaActive.value = true;
  };

  const validateForm = () => {
    let isValid = true;

    Object.keys(errors).forEach(key => errors[key] = '');

    if (!feedbackText.value.trim()) {
      errors.feedbackText = 'Por favor, comparte tus comentarios o sugerencias';
      isValid = false;
    } else if (feedbackText.value.length > 500) {
      errors.feedbackText = 'El texto no puede exceder los 500 caracteres';
      isValid = false;
    }

    if (userEmail.value.trim()) {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!emailRegex.test(userEmail.value.trim())) {
        errors.userEmail = 'Por favor, introduce un email válido';
        isValid = false;
      }
    }

    return isValid;
  };

  const submitFeedback = async () => {
    if (!validateForm()) {
      return;
    }

    isSubmitting.value = true;

    try {
      const feedbackData = {
        rating: rating.value,
        feedback: feedbackText.value,
        email: userEmail.value || 'No proporcionado',
      };

      const response = await fetch('https://formspree.io/f/mdkgbjez', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(feedbackData),
      });

      if (response.ok) {
        console.log('Feedback enviado correctamente');
        feedbackSubmitted.value = true;
      } else {
        console.error('Error al enviar el feedback');
        alert('Hubo un problema al enviar tu feedback. Por favor, intenta nuevamente.');
      }
    } catch (error) {
      console.error('Error al enviar el feedback:', error);
      alert('Hubo un problema de conexión. Por favor, intenta nuevamente.');
    } finally {
      isSubmitting.value = false;
    }
  };

  const resetForm = () => {
    rating.value = 0;
    feedbackText.value = '';
    userEmail.value = '';
    feedbackSubmitted.value = false;

    Object.keys(errors).forEach(key => errors[key] = '');
  };
</script>

<style lang="scss">
@import '/src/styles/variables.scss';

.feedback-page {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  padding: 0 1rem 4rem;

  .feedback-header {
    text-align: center;
    margin-bottom: 3rem;

    .page-title {
      font-family: $font-primary;
      font-size: 2.5rem;
      font-weight: 700;
      color: $primary-color-light-mode;
      margin-bottom: 1rem;

      .dark-mode & {
        color: $primary-color-dark-mode;
      }
    }

    .subtitle {
      font-size: 1.1rem;
      color: $text-color-light-mode;
      opacity: 0.8;
      max-width: 600px;
      margin: 0 auto;

      .dark-mode & {
        color: $text-color-dark-mode;
      }
    }
  }

  .feedback-container {
    position: relative;
    background-color: white;
    border-radius: 10px;
    box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
    overflow: hidden;

    .dark-mode & {
      background-color: #2a3238;
      box-shadow: 0 5px 20px rgba(0, 0, 0, 0.2);
    }

    .feedback-form {
      padding: 2rem;

      &.is-submitting {
        opacity: 0.7;
        pointer-events: none;
      }

      .form-section {
        margin-bottom: 2.5rem;

        &:last-of-type {
          margin-bottom: 1.5rem;
        }

        h3 {
          font-size: 1.2rem;
          margin-bottom: 0.5rem;
          color: $secondary-color-light-mode;

          .dark-mode & {
            color: $secondary-color-dark-mode;
          }
        }

        .section-description {
          font-size: 0.95rem;
          color: $text-color-light-mode;
          opacity: 0.8;
          margin-bottom: 1.5rem;

          .dark-mode & {
            color: $text-color-dark-mode;
          }
        }

        &.rating-section {
          .star-rating {
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 0.25rem;
            margin-bottom: 1rem;

            .star-container {
              position: relative;
              padding: 0.5rem;
              cursor: pointer;

              i {
                font-size: 2.5rem;
                color: #ccc;
                transition: color 0.2s ease, transform 0.2s ease;

                &.filled {
                  color: #FFD700;
                }

                &.selected {
                  transform: scale(1.1);
                }

                .dark-mode & {
                  color: #555;

                  &.filled {
                    color: #FFD700;
                  }
                }
              }

              .star-tooltip {
                position: absolute;
                top: -30px;
                left: 50%;
                transform: translateX(-50%);
                background-color: rgba(0, 0, 0, 0.7);
                color: white;
                padding: 0.25rem 0.5rem;
                border-radius: 4px;
                font-size: 0.8rem;
                white-space: nowrap;
                opacity: 0;
                visibility: hidden;
                transition: all 0.2s ease;

                &.active {
                  opacity: 1;
                  visibility: visible;
                }
              }

              &:hover {
                i {
                  transform: scale(1.2);
                }
              }
            }
          }

          .rating-label {
            text-align: center;
            font-size: 1rem;
            color: $text-color-light-mode;

            .dark-mode & {
              color: $text-color-dark-mode;
            }

            strong {
              font-weight: 600;
            }
          }
        }

        .textarea-wrapper {
          position: relative;
          margin-bottom: 0.5rem;

          textarea {
            width: 100%;
            min-height: 150px;
            padding: 1rem;
            font-size: 1rem;
            border: 1px solid #ddd;
            border-radius: 6px;
            resize: vertical;
            transition: all 0.3s ease;

            &:focus {
              outline: none;
              border-color: $primary-color-light-mode;
              box-shadow: 0 0 0 2px rgba($primary-color-light-mode, 0.2);
            }

            &.error {
              border-color: #e53935;

              &:focus {
                box-shadow: 0 0 0 2px rgba(#e53935, 0.2);
              }
            }

            .dark-mode & {
              background-color: #343e46;
              border-color: #555;
              color: $text-color-dark-mode;

              &::placeholder {
                color: #aaa;
              }

              &:focus {
                border-color: $primary-color-dark-mode;
                box-shadow: 0 0 0 2px rgba($primary-color-dark-mode, 0.2);
              }
            }
          }

          .feedback-suggestions {
            position: absolute;
            top: 1rem;
            left: 1rem;
            right: 1rem;
            background-color: white;
            border-radius: 6px;
            box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
            padding: 1rem;
            z-index: 5;
            animation: fadeIn 0.3s ease;

            .dark-mode & {
              background-color: #343e46;
              box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
            }

            .suggestion-prompt {
              font-size: 0.9rem;
              font-weight: 500;
              margin-bottom: 0.75rem;
              color: $secondary-color-light-mode;

              .dark-mode & {
                color: $secondary-color-dark-mode;
              }
            }

            .suggestion-list {
              list-style: none;
              padding: 0;
              margin: 0;

              li {
                padding: 0.5rem;
                margin-bottom: 0.25rem;
                font-size: 0.9rem;
                color: $text-color-light-mode;
                border-radius: 4px;
                cursor: pointer;
                display: flex;
                align-items: center;

                &:last-child {
                  margin-bottom: 0;
                }

                i {
                  margin-right: 0.5rem;
                  font-size: 1.1rem;
                  color: $primary-color-light-mode;
                }

                &:hover {
                  background-color: rgba($primary-color-light-mode, 0.1);
                }

                .dark-mode & {
                  color: $text-color-dark-mode;

                  i {
                    color: $primary-color-dark-mode;
                  }

                  &:hover {
                    background-color: rgba($primary-color-dark-mode, 0.2);
                  }
                }
              }
            }
          }
        }

        .character-count {
          font-size: 0.85rem;
          text-align: right;
          color: #888;

          &.limit-close {
            color: #e53935;
          }

          .dark-mode & {
            color: #aaa;

            &.limit-close {
              color: #ff6b6b;
            }
          }
        }

        input[type="email"] {
          width: 100%;
          padding: 0.75rem 1rem;
          font-size: 1rem;
          border: 1px solid #ddd;
          border-radius: 6px;
          transition: all 0.3s ease;

          &:focus {
            outline: none;
            border-color: $primary-color-light-mode;
            box-shadow: 0 0 0 2px rgba($primary-color-light-mode, 0.2);
          }

          &.error {
            border-color: #e53935;

            &:focus {
              box-shadow: 0 0 0 2px rgba(#e53935, 0.2);
            }
          }

          .dark-mode & {
            background-color: #343e46;
            border-color: #555;
            color: $text-color-dark-mode;

            &::placeholder {
              color: #aaa;
            }

            &:focus {
              border-color: $primary-color-dark-mode;
              box-shadow: 0 0 0 2px rgba($primary-color-dark-mode, 0.2);
            }
          }
        }

        .error-message {
          display: block;
          margin-top: 0.5rem;
          font-size: 0.85rem;
          color: #e53935;
        }
      }

      .form-actions {
        display: flex;
        justify-content: center;
        margin-top: 2rem;

        .submit-button {
          background-color: $primary-color-light-mode;
          color: white;
          border: none;
          padding: 0.75rem 2.5rem;
          border-radius: 6px;
          font-size: 1rem;
          font-weight: 500;
          cursor: pointer;
          transition: all 0.3s ease;
          display: flex;
          align-items: center;
          justify-content: center;
          min-width: 180px;

          &:hover {
            background-color: darken($primary-color-light-mode, 10%);
            transform: translateY(-2px);
          }

          &:active {
            transform: translateY(0);
          }

          &:disabled {
            opacity: 0.7;
            cursor: not-allowed;
            transform: none;
          }

          .dark-mode & {
            background-color: $primary-color-dark-mode;
            color: #081925;

            &:hover {
              background-color: darken($primary-color-dark-mode, 10%);
            }
          }

          .submitting {
            display: flex;
            align-items: center;
            gap: 0.5rem;

            i {
              font-size: 1.2rem;
            }
          }
        }
      }
    }

    .success-message {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background-color: white;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      text-align: center;
      padding: 2rem;
      animation: fadeIn 0.5s ease;

      .dark-mode & {
        background-color: #2a3238;
      }

      .success-icon {
        width: 80px;
        height: 80px;
        border-radius: 50%;
        background-color: rgba($primary-color-light-mode, 0.1);
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 1.5rem;

        i {
          font-size: 3rem;
          color: $primary-color-light-mode;
        }

        .dark-mode & {
          background-color: rgba($primary-color-dark-mode, 0.2);

          i {
            color: $primary-color-dark-mode;
          }
        }
      }

      h3 {
        font-size: 1.5rem;
        margin-bottom: 1rem;
        color: $secondary-color-light-mode;

        .dark-mode & {
          color: $secondary-color-dark-mode;
        }
      }

      p {
        font-size: 1rem;
        color: $text-color-light-mode;
        margin-bottom: 2rem;
        max-width: 500px;

        .dark-mode & {
          color: $text-color-dark-mode;
        }
      }

      .reset-button {
        background-color: $primary-color-light-mode;
        color: white;
        border: none;
        padding: 0.75rem 2rem;
        border-radius: 6px;
        font-size: 1rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.3s ease;

        &:hover {
          background-color: darken($primary-color-light-mode, 10%);
          transform: translateY(-2px);
        }

        &:active {
          transform: translateY(0);
        }

        .dark-mode & {
          background-color: $primary-color-dark-mode;

          &:hover {
            background-color: darken($primary-color-dark-mode, 10%);
          }
        }
      }
    }
  }

  // Animaciones
  @keyframes fadeIn {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
}
</style>
