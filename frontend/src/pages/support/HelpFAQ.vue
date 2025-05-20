<template>
  <div class="help-faq-page">
    <div class="help-header">
      <h1 class="page-title">Ayuda / Preguntas Frecuentes</h1>
      <p class="subtitle">Encuentra respuestas a las preguntas más comunes sobre SignalSysAnalyzer</p>

      <div class="search-container">
        <div class="search-box">
          <i class="bx bx-search search-icon"></i>
          <input
            v-model="searchQuery"
            class="search-input"
            placeholder="Buscar en las preguntas frecuentes..."
            type="text"
            @input="filterFAQs"
          />
          <button v-if="searchQuery" class="clear-search" @click="clearSearch">
            <i class="bx bx-x"></i>
          </button>
        </div>
      </div>
    </div>

    <div class="faq-container">
      <div class="faq-categories">
        <button
          v-for="category in categories"
          :key="category.id"
          :class="['category-button', { active: selectedCategory === category.id }]"
          @click="selectCategory(category.id)"
        >
          <i :class="category.icon"></i>
          <span>{{ category.name }}</span>
        </button>
      </div>

      <div class="faq-list">
        <div v-if="filteredFAQs.length === 0" class="no-results">
          <i class="bx bx-search-alt no-results-icon"></i>
          <p>No se encontraron resultados para "{{ searchQuery }}"</p>
          <button class="reset-search-btn" @click="clearSearch">Mostrar todas las preguntas</button>
        </div>

        <transition-group class="faq-accordion" name="faq-item" tag="div">
          <div v-for="(faq, index) in filteredFAQs" :key="faq.id" class="faq-item">
            <div
              class="faq-question"
              :class="{ 'open': openFAQs.includes(faq.id) }"
              @click="toggleFAQ(faq.id)"
            >
              <div class="question-content">
                <span class="question-number">{{ index + 1 }}</span>
                <h3>{{ faq.question }}</h3>
              </div>
              <div class="toggle-icon">
                <i :class="openFAQs.includes(faq.id) ? 'bx bx-minus' : 'bx bx-plus'"></i>
              </div>
            </div>
            <div
              class="faq-answer"
              :class="{ 'open': openFAQs.includes(faq.id) }"
              :style="{ maxHeight: openFAQs.includes(faq.id) ? faq.answerHeight : '0px' }"
            >
              <div ref="answerContent" class="answer-content">
                <p>{{ faq.answer }}</p>
                <div v-if="faq.additionalInfo" class="additional-info">
                  <i class="bx bx-info-circle"></i>
                  <span>{{ faq.additionalInfo }}</span>
                </div>
              </div>
            </div>
          </div>
        </transition-group>
      </div>
    </div>

    <div class="help-contact">
      <h2>¿No encontraste lo que buscabas?</h2>
      <p>Si tienes alguna otra pregunta o necesitas asistencia personalizada, no dudes en contactarnos.</p>
      <div class="contact-buttons">
        <router-link class="contact-button primary" to="/support/contact">
          <i class="bx bx-envelope"></i>
          Contáctanos
        </router-link>
        <a  class="contact-button secondary" href="mailto:tvega@unillanos.edu.co">
          <i class="bx bx-at"></i>
          Correo de soporte
        </a>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { computed, onMounted, ref } from 'vue';

  const searchQuery = ref('');
  const selectedCategory = ref('all');
  const openFAQs = ref([]);
  const answerContent = ref([]);

  const categories = [
    { id: 'all', name: 'Todas', icon: 'bx bx-list-ul' },
    { id: 'privacy', name: 'Privacidad y Seguridad', icon: 'bx bx-shield-quarter' },
    { id: 'usage', name: 'Uso de la plataforma', icon: 'bx bx-laptop' },
    { id: 'account', name: 'Cuenta y perfil', icon: 'bx bx-user' },
  ];

  const faqs = [
    {
      id: 1,
      question: '¿Se mantiene una copia de mis archivos procesados?',
      answer: 'Sí, conservamos los archivos para fines de mejora de nuestros modelos de Machine Learning. Sin embargo, el archivo original no es accesible ni visible por otros usuarios.',
      category: 'privacy',
      answerHeight: 'auto',
    },
    {
      id: 2,
      question: '¿Los archivos de mi empresa están seguros con SignalSysAnalyzer?',
      answer: 'La seguridad y privacidad de tus datos es nuestra prioridad. Los archivos se almacenan de forma encriptada y cumplen con estándares de protección de datos.',
      category: 'privacy',
      answerHeight: 'auto',
    },
    {
      id: 3,
      question: '¿Cuáles son los requerimientos mínimos del sistema?',
      answer: 'Navegador actualizado (Chrome, Firefox, Edge), conexión a internet estable, archivo CSV con al menos 10,000 datos.',
      category: 'usage',
      answerHeight: 'auto',
    },
    {
      id: 4,
      question: '¿Cómo puedo subir mis archivos?',
      answer: 'Puedes arrastrar el archivo a la zona indicada o seleccionarlo manualmente desde tu dispositivo.',
      category: 'usage',
      answerHeight: 'auto',
    },
    {
      id: 5,
      question: '¿Por qué mis resultados están tardando tanto?',
      answer: 'Dependiendo del tamaño del archivo o la carga del servidor, los resultados pueden demorar algunos segundos. Si el tiempo excede los 2 minutos, intenta nuevamente o contáctanos.',
      category: 'usage',
      additionalInfo: 'Los tiempos de procesamiento típicos son de 15-30 segundos.',
      answerHeight: 'auto',
    },
    {
      id: 6,
      question: '¿Puedo cambiar mi contraseña?',
      answer: 'Sí, desde tu perfil puedes modificarla en cualquier momento.',
      category: 'account',
      answerHeight: 'auto',
    },
    {
      id: 7,
      question: '¿Puedo descargar mis resultados?',
      answer: 'Sí, en la sección de perfil puedes descargar un informe PDF de cada análisis realizado.',
      category: 'usage',
      answerHeight: 'auto',
    },
    {
      id: 8,
      question: '¿Puedo usar la herramienta desde el celular?',
      answer: 'Sí, aunque recomendamos usar un computador para una mejor visualización y manejo de archivos.',
      category: 'usage',
      answerHeight: 'auto',
    },
    {
      id: 9,
      question: '¿Necesito una cuenta para usar SignalSysAnalyzer?',
      answer: 'No es obligatorio, pero tener una cuenta te permite guardar tu historial, descargar informes y configurar tu perfil.',
      category: 'account',
      answerHeight: 'auto',
    },
    {
      id: 10,
      question: '¿Se guardan mis resultados automáticamente?',
      answer: 'Sí, si estás registrado. Los resultados se almacenan en tu perfil con fecha y nombre de archivo.',
      category: 'account',
      answerHeight: 'auto',
    },
    {
      id: 11,
      question: '¿Usan mis datos con fines comerciales?',
      answer: 'No. Tus datos se utilizan exclusivamente para fines de análisis técnico y mejora del sistema, nunca con fines comerciales.',
      category: 'privacy',
      answerHeight: 'auto',
    },
  ];

  const filteredFAQs = computed(() => {
    let filtered = [...faqs];

    if (selectedCategory.value !== 'all') {
      filtered = filtered.filter(faq => faq.category === selectedCategory.value);
    }

    if (searchQuery.value.trim() !== '') {
      const query = searchQuery.value.toLowerCase();
      filtered = filtered.filter(faq =>
        faq.question.toLowerCase().includes(query) ||
        faq.answer.toLowerCase().includes(query)
      );
    }

    return filtered;
  });

  const selectCategory = categoryId => {
    selectedCategory.value = categoryId;
  };

  const toggleFAQ = faqId => {
    if (openFAQs.value.includes(faqId)) {
      openFAQs.value = openFAQs.value.filter(id => id !== faqId);
    } else {
      openFAQs.value.push(faqId);
    }
  };

  const filterFAQs = () => {
    if (searchQuery.value.trim() !== '') {
      openFAQs.value = filteredFAQs.value.map(faq => faq.id);
    }
  };

  const clearSearch = () => {
    searchQuery.value = '';
    openFAQs.value = [];
  };

  onMounted(() => {
    if (faqs.length > 0) {
      openFAQs.value = [faqs[0].id];
    }

    nextTick(() => {
      const elements = document.querySelectorAll('.answer-content');
      elements.forEach((el, index) => {
        if (index < faqs.length) {
          faqs[index].answerHeight = `${el.offsetHeight + 40}px`; // Padding extra
        }
      });
    });
  });
</script>

<style lang="scss">
@import '/src/styles/variables.scss';

.help-faq-page {
  width: 100%;
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 1rem 4rem;

  .help-header {
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
      margin: 0 auto 2rem;

      .dark-mode & {
        color: $text-color-dark-mode;
      }
    }

    .search-container {
      max-width: 600px;
      margin: 0 auto;

      .search-box {
        position: relative;
        display: flex;
        align-items: center;

        .search-icon {
          position: absolute;
          left: 1rem;
          font-size: 1.2rem;
          color: #888;
          pointer-events: none;

          .dark-mode & {
            color: #aaa;
          }
        }

        .search-input {
          width: 100%;
          padding: 1rem 1rem 1rem 3rem;
          border: 1px solid #e0e0e0;
          border-radius: 30px;
          font-size: 1rem;
          transition: all 0.3s ease;

          &:focus {
            outline: none;
            border-color: $primary-color-light-mode;
            box-shadow: 0 0 0 2px rgba($primary-color-light-mode, 0.2);
          }

          .dark-mode & {
            background-color: #343e46;
            color: $text-color-dark-mode;
            border-color: #555;

            &:focus {
              border-color: $primary-color-dark-mode;
              box-shadow: 0 0 0 2px rgba($primary-color-dark-mode, 0.2);
            }

            &::placeholder {
              color: #aaa;
            }
          }
        }

        .clear-search {
          position: absolute;
          right: 1rem;
          background: none;
          border: none;
          cursor: pointer;
          color: #888;
          font-size: 1.2rem;
          display: flex;
          align-items: center;
          justify-content: center;

          &:hover {
            color: $primary-color-light-mode;
          }

          .dark-mode & {
            color: #aaa;

            &:hover {
              color: $primary-color-dark-mode;
            }
          }
        }
      }
    }
  }

  .faq-container {
    background-color: white;
    border-radius: 10px;
    box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
    overflow: hidden;
    margin-bottom: 4rem;

    .dark-mode & {
      background-color: #2a3238;
      box-shadow: 0 5px 20px rgba(0, 0, 0, 0.2);
    }

    .faq-categories {
      display: flex;
      overflow-x: auto;
      padding: 1rem;
      gap: 0.5rem;
      background-color: #f9f9f9;

      @media (max-width: 768px) {
        flex-wrap: nowrap;
        justify-content: flex-start;

        &::-webkit-scrollbar {
          height: 3px;
        }

        &::-webkit-scrollbar-thumb {
          background-color: rgba(0, 0, 0, 0.2);
          border-radius: 3px;
        }
      }

      .dark-mode & {
        background-color: #343e46;
      }

      .category-button {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1.25rem;
        border-radius: 30px;
        border: none;
        background-color: transparent;
        color: $text-color-light-mode;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.3s ease;
        white-space: nowrap;

        i {
          font-size: 1.1rem;
        }

        &:hover {
          background-color: rgba($primary-color-light-mode, 0.1);
          color: $primary-color-light-mode;
        }

        &.active {
          background-color: $primary-color-light-mode;
          color: white;
        }

        .dark-mode & {
          color: $text-color-dark-mode;

          &:hover {
            background-color: rgba($primary-color-dark-mode, 0.2);
            color: $primary-color-dark-mode;
          }

          &.active {
            background-color: $primary-color-dark-mode;
            color: $background-color-dark-mode;
          }
        }
      }
    }

    .faq-list {
      padding: 0;

      .no-results {
        padding: 3rem;
        text-align: center;

        .no-results-icon {
          font-size: 3rem;
          color: #ccc;
          margin-bottom: 1rem;

          .dark-mode & {
            color: #555;
          }
        }

        p {
          margin-bottom: 1.5rem;
          color: $text-color-light-mode;

          .dark-mode & {
            color: $text-color-dark-mode;
          }
        }

        .reset-search-btn {
          background-color: $primary-color-light-mode;
          color: white;
          border: none;
          padding: 0.75rem 1.5rem;
          border-radius: 30px;
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

      .faq-accordion {
        .faq-item {
          border-bottom: 1px solid #eee;

          .dark-mode & {
            border-bottom-color: #444;
          }

          &:last-child {
            border-bottom: none;
          }

          .faq-question {
            padding: 1.5rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
            cursor: pointer;

            &:hover {
              background-color: rgba($primary-color-light-mode, 0.05);

              .dark-mode & {
                background-color: rgba($primary-color-dark-mode, 0.1);
              }
            }

            &.open {
              background-color: rgba($primary-color-light-mode, 0.1);

              .dark-mode & {
                background-color: rgba($primary-color-dark-mode, 0.2);
              }
            }

            .question-content {
              display: flex;
              align-items: center;
              gap: 1rem;
              flex: 1;

              .question-number {
                display: flex;
                align-items: center;
                justify-content: center;
                width: 28px;
                height: 28px;
                border-radius: 50%;
                background-color: $primary-color-light-mode;
                color: white;
                font-size: 0.85rem;
                font-weight: 600;
                flex-shrink: 0;

                .dark-mode & {
                  background-color: $primary-color-dark-mode;
                  color: #081925;
                }
              }

              h3 {
                margin: 0;
                font-size: 1.1rem;
                font-weight: 600;
                color: $text-color-light-mode;

                .dark-mode & {
                  color: $text-color-dark-mode;
                }
              }
            }

            .toggle-icon {
              font-size: 1.1rem;
              color: $primary-color-light-mode;
              display: flex;
              align-items: center;
              justify-content: center;
              width: 30px;
              height: 30px;
              border-radius: 50%;
              transition: all 0.3s ease;

              .dark-mode & {
                color: $primary-color-dark-mode;
              }
            }
          }

          .faq-answer {
            max-height: 0;
            overflow: hidden;
            transition: max-height 0.3s ease;

            .answer-content {
              padding: 0 1.5rem 1.5rem 4rem;

              p {
                margin-top: 10px;
                line-height: 1.6;
                color: $text-color-light-mode;
                opacity: 0.9;

                .dark-mode & {
                  color: $text-color-dark-mode;
                }
              }

              .additional-info {
                margin-top: 1rem;
                padding: 0.75rem 1rem;
                background-color: rgba($primary-color-light-mode, 0.1);
                border-radius: 6px;
                display: flex;
                align-items: flex-start;
                gap: 0.75rem;
                font-size: 0.9rem;

                i {
                  font-size: 1.1rem;
                  color: $primary-color-light-mode;
                  margin-top: 2px;
                }

                .dark-mode & {
                  background-color: rgba($primary-color-dark-mode, 0.2);

                  i {
                    color: $primary-color-dark-mode;
                  }
                }
              }
            }
          }
        }
      }
    }
  }

  .help-contact {
    text-align: center;
    padding: 2rem;
    background-color: white;
    border-radius: 10px;
    box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);

    .dark-mode & {
      background-color: #2a3238;
      box-shadow: 0 5px 20px rgba(0, 0, 0, 0.2);
    }

    h2 {
      font-size: 1.5rem;
      margin-bottom: 1rem;
      color: $secondary-color-light-mode;

      .dark-mode & {
        color: $secondary-color-dark-mode;
      }
    }

    p {
      margin-bottom: 1.5rem;
      color: $text-color-light-mode;

      .dark-mode & {
        color: $text-color-dark-mode;
      }
    }

    .contact-buttons {
      display: flex;
      flex-wrap: wrap;
      justify-content: center;
      gap: 1rem;

      .contact-button {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1.5rem;
        border-radius: 30px;
        text-decoration: none;
        font-weight: 500;
        transition: all 0.3s ease;

        i {
          font-size: 1.1rem;
        }

        &.primary {
          background-color: $primary-color-light-mode;
          color: white;

          &:hover {
            background-color: darken($primary-color-light-mode, 10%);
            transform: translateY(-2px);
          }

          .dark-mode & {
            background-color: $primary-color-dark-mode;
            color: #081925;

            &:hover {
              background-color: darken($primary-color-dark-mode, 10%);
            }
          }
        }

        &.secondary {
          background-color: #f0f0f0;
          color: $text-color-light-mode;

          &:hover {
            background-color: #e5e5e5;
            transform: translateY(-2px);
          }

          .dark-mode & {
            background-color: #343e46;
            color: $text-color-dark-mode;

            &:hover {
              background-color: #3c4750;
            }
          }
        }
      }
    }
  }

  // Animaciones
  .faq-item-enter-active,
  .faq-item-leave-active {
    transition: all 0.3s ease;
  }

  .faq-item-enter-from,
  .faq-item-leave-to {
    opacity: 0;
    transform: translateY(20px);
  }
}
</style>
