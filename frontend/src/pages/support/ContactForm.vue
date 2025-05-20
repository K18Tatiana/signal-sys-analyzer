<template>
  <div class="contact-page">
    <div class="contact-header">
      <h1 class="page-title">Contacto</h1>
      <p class="subtitle">Estamos aquí para escucharte. Envíanos tu mensaje y te responderemos lo antes posible.</p>
    </div>

    <div class="contact-container">
      <div class="contact-info">
        <div class="info-card">
          <div class="card-icon">
            <i class="bx bx-envelope"></i>
          </div>
          <h3>Correo Electrónico</h3>
          <p>tvega@unillanos.edu.co</p>
        </div>

        <div class="info-card">
          <div class="card-icon">
            <i class="bx bx-time-five"></i>
          </div>
          <h3>Tiempo de Respuesta</h3>
          <p>Respondemos la mayoría de consultas en un plazo de 24-48 horas en días laborables.</p>
        </div>

        <div class="info-card">
          <div class="card-icon">
            <i class="bx bx-help-circle"></i>
          </div>
          <h3>Preguntas Frecuentes</h3>
          <p>Revisa nuestra <router-link to="/support/help">página de ayuda</router-link> para encontrar respuestas rápidas a preguntas comunes.</p>
        </div>
      </div>

      <div class="form-container">
        <form class="contact-form" :class="{ 'is-submitting': isSubmitting }" @submit.prevent="submitForm">
          <div class="form-header">
            <h2>Formulario de Contacto</h2>
            <p>Todos los campos marcados con <span class="required">*</span> son obligatorios</p>
          </div>

          <div class="form-group">
            <label for="name">Nombre <span class="required">*</span></label>
            <input
              id="name"
              v-model="formData.name"
              :class="{ 'error': errors.name }"
              placeholder="Ingresa tu nombre completo"
              required
              type="text"
            />
            <span v-if="errors.name" class="error-message">{{ errors.name }}</span>
          </div>

          <div class="form-group">
            <label for="email">Email <span class="required">*</span></label>
            <input
              id="email"
              v-model="formData.email"
              :class="{ 'error': errors.email }"
              placeholder="ejemplo@correo.com"
              required
              type="email"
            />
            <span v-if="errors.email" class="error-message">{{ errors.email }}</span>
          </div>

          <div class="form-group">
            <label for="subject">Asunto <span class="required">*</span></label>
            <select
              id="subject"
              v-model="formData.subject"
              :class="{ 'error': errors.subject }"
              required
            >
              <option disabled selected value="">Selecciona un asunto</option>
              <option value="info">Información general</option>
              <option value="feature">Sugerir una funcionalidad o herramienta</option>
              <option value="bug">Informar de un problema</option>
              <option value="privacy">Privacidad</option>
            </select>
            <span v-if="errors.subject" class="error-message">{{ errors.subject }}</span>
          </div>

          <div class="form-group">
            <label for="message">Mensaje <span class="required">*</span></label>
            <textarea
              id="message"
              v-model="formData.message"
              :class="{ 'error': errors.message }"
              placeholder="Escribe tu mensaje aquí..."
              required
              rows="5"
            ></textarea>
            <span v-if="errors.message" class="error-message">{{ errors.message }}</span>
          </div>

          <div class="form-group checkbox-group">
            <div class="checkbox-container">
              <input
                id="terms"
                v-model="formData.acceptTerms"
                :class="{ 'error': errors.acceptTerms }"
                required
                type="checkbox"
              />
              <label class="checkbox-label" for="terms">
                Acepto los <a href="/legal/terms" target="_blank">términos y condiciones</a>, y la <a href="/legal/privacy" target="_blank">política de privacidad</a> <span class="required">*</span>
              </label>
            </div>
            <span v-if="errors.acceptTerms" class="error-message">{{ errors.acceptTerms }}</span>
          </div>

          <div class="form-actions">
            <button class="submit-button" :disabled="isSubmitting" type="submit">
              <span v-if="!isSubmitting">Enviar mensaje</span>
              <span v-else class="submitting">
                <i class="bx bx-loader-alt bx-spin"></i>
                Enviando...
              </span>
            </button>
          </div>
        </form>

        <!-- Mensaje de éxito -->
        <div v-if="formSubmitted" class="success-message">
          <div class="success-icon">
            <i class="bx bx-check"></i>
          </div>
          <h3>¡Mensaje enviado con éxito!</h3>
          <p>Gracias por contactarnos. Hemos recibido tu mensaje y te responderemos lo antes posible.</p>
          <button class="reset-button" @click="resetForm">Enviar otro mensaje</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { reactive, ref } from 'vue';

  const formData = reactive({
    name: '',
    email: '',
    subject: '',
    message: '',
    acceptTerms: false,
  });

  const errors = reactive({
    name: '',
    email: '',
    subject: '',
    message: '',
    acceptTerms: '',
  });

  const isSubmitting = ref(false);
  const formSubmitted = ref(false);

  const validateForm = () => {
    let isValid = true;

    Object.keys(errors).forEach(key => errors[key] = '');

    if (!formData.name.trim()) {
      errors.name = 'El nombre es obligatorio';
      isValid = false;
    } else if (formData.name.trim().length < 3) {
      errors.name = 'El nombre debe tener al menos 3 caracteres';
      isValid = false;
    }

    // Validar email
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!formData.email.trim()) {
      errors.email = 'El email es obligatorio';
      isValid = false;
    } else if (!emailRegex.test(formData.email.trim())) {
      errors.email = 'Por favor, introduce un email válido';
      isValid = false;
    }

    // Validar asunto
    if (!formData.subject) {
      errors.subject = 'Por favor, selecciona un asunto';
      isValid = false;
    }

    // Validar mensaje
    if (!formData.message.trim()) {
      errors.message = 'El mensaje es obligatorio';
      isValid = false;
    } else if (formData.message.trim().length < 10) {
      errors.message = 'El mensaje debe tener al menos 10 caracteres';
      isValid = false;
    }

    // Validar términos
    if (!formData.acceptTerms) {
      errors.acceptTerms = 'Debes aceptar los términos y condiciones';
      isValid = false;
    }

    return isValid;
  };

  const submitForm = async () => {
    if (!validateForm()) return;
    isSubmitting.value = true;

    try {
      const response = await fetch('https://formspree.io/f/mzzrowdw', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: formData.name,
          email: formData.email,
          subject: formData.subject,
          message: formData.message,
        }),
      });

      if (response.ok) {
        formSubmitted.value = true;
      } else {
        console.error('Error al enviar el formulario');
      }
    } catch (error) {
      console.error('Error:', error);
    } finally {
      isSubmitting.value = false;
    }
  };

  const resetForm = () => {
    Object.keys(formData).forEach(key => {
      if (typeof formData[key] === 'boolean') {
        formData[key] = false;
      } else {
        formData[key] = '';
      }
    });

    formSubmitted.value = false;

    Object.keys(errors).forEach(key => errors[key] = '');
  };
</script>

<style lang="scss">
@import '/src/styles/variables.scss';

.contact-page {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem 4rem;

  .contact-header {
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
      max-width: 700px;
      margin: 0 auto;

      .dark-mode & {
        color: $text-color-dark-mode;
      }
    }
  }

  .contact-container {
    display: flex;
    flex-direction: column;
    gap: 2rem;

    @media (min-width: 992px) {
      flex-direction: row;
      align-items: flex-start;
    }

    .contact-info {
      width: 100%;
      display: flex;
      flex-direction: column;
      gap: 1.5rem;

      @media (min-width: 992px) {
        width: 30%;
      }

      .info-card {
        background-color: white;
        border-radius: 10px;
        padding: 1.5rem;
        box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);

        .dark-mode & {
          background-color: #2a3238;
          box-shadow: 0 5px 20px rgba(0, 0, 0, 0.2);
        }

        .card-icon {
          width: 50px;
          height: 50px;
          border-radius: 50%;
          background-color: rgba($primary-color-light-mode, 0.1);
          display: flex;
          align-items: center;
          justify-content: center;
          margin-bottom: 1rem;

          i {
            font-size: 1.5rem;
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
          font-size: 1.1rem;
          margin-bottom: 0.75rem;
          color: $secondary-color-light-mode;

          .dark-mode & {
            color: $secondary-color-dark-mode;
          }
        }

        p {
          font-size: 0.95rem;
          color: $text-color-light-mode;
          margin-bottom: 0.5rem;

          &:last-child {
            margin-bottom: 0;
          }

          .dark-mode & {
            color: $text-color-dark-mode;
          }

          a {
            color: $primary-color-light-mode;
            text-decoration: none;

            &:hover {
              text-decoration: underline;
            }

            .dark-mode & {
              color: $primary-color-dark-mode;
            }
          }
        }
      }
    }

    .form-container {
      width: 100%;
      position: relative;

      @media (min-width: 992px) {
        width: 70%;
      }

      .contact-form {
        background-color: white;
        border-radius: 10px;
        padding: 2rem;
        box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);

        .dark-mode & {
          background-color: #2a3238;
          box-shadow: 0 5px 20px rgba(0, 0, 0, 0.2);
        }

        &.is-submitting {
          opacity: 0.7;
          pointer-events: none;
        }

        .form-header {
          margin-bottom: 2rem;

          h2 {
            font-size: 1.5rem;
            color: $secondary-color-light-mode;
            margin-bottom: 0.5rem;

            .dark-mode & {
              color: $secondary-color-dark-mode;
            }
          }

          p {
            font-size: 0.9rem;
            color: $text-color-light-mode;
            opacity: 0.8;

            .dark-mode & {
              color: $text-color-dark-mode;
            }

            .required {
              color: #e53935;
            }
          }
        }

        .form-group {
          margin-bottom: 1.5rem;

          label {
            display: block;
            margin-bottom: 0.5rem;
            font-size: 0.95rem;
            font-weight: 500;
            color: $text-color-light-mode;

            .dark-mode & {
              color: $text-color-dark-mode;
            }

            .required {
              color: #e53935;
            }
          }

          input[type="text"],
          input[type="email"],
          select,
          textarea {
            width: 100%;
            padding: 0.75rem 1rem;
            border: 1px solid #ddd;
            border-radius: 6px;
            font-size: 1rem;
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

              &.error {
                border-color: #e53935;

                &:focus {
                  box-shadow: 0 0 0 2px rgba(#e53935, 0.2);
                }
              }
            }
          }

          select {
            appearance: none;
            background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='%23555555' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
            background-repeat: no-repeat;
            background-position: right 1rem center;
            background-size: 1rem;
            padding-right: 2.5rem;

            .dark-mode & {
              background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='24' height='24' viewBox='0 0 24 24' fill='none' stroke='%23aaaaaa' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
            }
          }

          textarea {
            resize: vertical;
            min-height: 120px;
          }

          .error-message {
            display: block;
            margin-top: 0.5rem;
            font-size: 0.85rem;
            color: #e53935;
          }

          &.checkbox-group {
            .checkbox-container {
              display: flex;
              align-items: flex-start;
              gap: 0.75rem;

              input[type="checkbox"] {
                margin-top: 0.25rem;
                width: 20px;

                &.error {
                  outline: 2px solid #e53935;
                  border-radius: 2px;
                }
              }

              .checkbox-label {
                margin-bottom: 0;
                font-size: 0.9rem;
                line-height: 1.5;

                a {
                  color: $primary-color-light-mode;
                  text-decoration: none;

                  &:hover {
                    text-decoration: underline;
                  }

                  .dark-mode & {
                    color: $primary-color-dark-mode;
                  }
                }
              }
            }
          }
        }

        .form-actions {
          display: flex;
          justify-content: flex-end;

          .submit-button {
            background-color: $primary-color-light-mode;
            color: white;
            border: none;
            padding: 0.75rem 2rem;
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
            }

            &:disabled {
              opacity: 0.7;
              cursor: not-allowed;
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
        border-radius: 10px;
        padding: 3rem 2rem;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        text-align: center;
        box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
        animation: fadeIn 0.5s ease;

        .dark-mode & {
          background-color: #2a3238;
          box-shadow: 0 5px 20px rgba(0, 0, 0, 0.2);
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
