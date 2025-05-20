<template>
  <div class="terms-page">
    <div class="terms-header">
      <h1 class="page-title">Términos y Condiciones</h1>
      <p class="subtitle">Condiciones de uso de SignalSysAnalyzer</p>
    </div>

    <div class="terms-container">
      <div class="terms-intro">
        <div class="intro-icon">
          <i class="bx bx-file"></i>
        </div>
        <div class="intro-content">
          <p>
            Al utilizar SignalSysAnalyzer, aceptas los siguientes términos de uso. Te recomendamos leer detenidamente
            este documento para entender tus derechos y responsabilidades al usar nuestra plataforma.
          </p>
        </div>
      </div>

      <!-- Secciones -->
      <div class="terms-main">
        <div class="terms-card">
          <div class="terms-sections">
            <!-- Sección 1: Uso del servicio -->
            <div id="service-usage" class="terms-section">
              <div class="section-header">
                <div class="section-number">1</div>
                <h2>Uso del servicio</h2>
              </div>
              <div class="section-content">
                <div class="service-container">
                  <div class="service-icon"><i class="bx bx-laptop"></i></div>
                  <p>
                    SignalSysAnalyzer es una herramienta educativa para facilitar el análisis de sistemas físicos.
                    No debe usarse como única fuente de diagnóstico en entornos críticos. La plataforma está diseñada
                    para complementar los conocimientos académicos y profesionales, no para reemplazarlos.
                  </p>
                </div>
              </div>
            </div>

            <!-- Sección 2: Exactitud de los resultados -->
            <div id="results-accuracy" class="terms-section">
              <div class="section-header">
                <div class="section-number">2</div>
                <h2>Exactitud de los resultados</h2>
              </div>
              <div class="section-content">
                <div class="accuracy-container">
                  <div class="accuracy-icon"><i class="bx bx-chart"></i></div>
                  <p>
                    Aunque empleamos modelos basados en Machine Learning y validaciones internas, los resultados
                    no están garantizados al 100% y deben ser evaluados por un profesional. Nuestros algoritmos
                    están en constante mejora, pero pueden existir limitaciones inherentes a la naturaleza de los
                    datos analizados.
                  </p>
                </div>
              </div>
            </div>

            <!-- Sección 3: Responsabilidad limitada -->
            <div id="limited-liability" class="terms-section">
              <div class="section-header">
                <div class="section-number">3</div>
                <h2>Responsabilidad limitada</h2>
              </div>
              <div class="section-content">
                <div class="liability-alert">
                  <div class="alert-icon">
                    <i class="bx bx-error-circle"></i>
                  </div>
                  <div class="alert-content">
                    <p>
                      No nos hacemos responsables por decisiones tomadas basadas en los resultados entregados
                      por la herramienta. El usuario asume la responsabilidad de verificar y validar cualquier
                      resultado obtenido antes de utilizarlo en aplicaciones reales.
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <!-- Sección 4: Propiedad intelectual -->
            <div id="intellectual-property" class="terms-section">
              <div class="section-header">
                <div class="section-number">4</div>
                <h2>Propiedad intelectual</h2>
              </div>
              <div class="section-content">
                <p>
                  Todo el contenido de esta plataforma (software, diseño, nombre, logotipo) pertenece a
                  SignalSysAnalyzer o a sus creadores originales. Está prohibida la reproducción, distribución,
                  modificación o uso no autorizado de cualquier elemento de nuestra plataforma.
                </p>
                <div class="ip-examples">
                  <div class="ip-item">
                    <div class="ip-icon"><i class="bx bx-code-alt"></i></div>
                    <div class="ip-text">Software y algoritmos</div>
                  </div>
                  <div class="ip-item">
                    <div class="ip-icon"><i class="bx bx-palette"></i></div>
                    <div class="ip-text">Diseño y elementos visuales</div>
                  </div>
                  <div class="ip-item">
                    <div class="ip-icon"><i class="bx bx-registered"></i></div>
                    <div class="ip-text">Nombre y logotipo</div>
                  </div>
                  <div class="ip-item">
                    <div class="ip-icon"><i class="bx bx-book-content"></i></div>
                    <div class="ip-text">Documentación y textos</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Sección 5: Modificaciones -->
            <div id="modifications" class="terms-section">
              <div class="section-header">
                <div class="section-number">5</div>
                <h2>Modificaciones</h2>
              </div>
              <div class="section-content">
                <p>
                  Nos reservamos el derecho de actualizar estos términos. Los usuarios registrados serán
                  notificados sobre cualquier cambio importante. Recomendamos revisar periódicamente esta
                  página para estar al tanto de posibles actualizaciones.
                </p>
                <div class="update-note">
                  <div class="update-icon"><i class="bx bx-bell"></i></div>
                  <p>
                    Las modificaciones importantes serán notificadas por correo electrónico con al menos
                    15 días de anticipación antes de su entrada en vigor.
                  </p>
                </div>
              </div>
            </div>

            <!-- Sección 6: Privacidad -->
            <div id="privacy" class="terms-section">
              <div class="section-header">
                <div class="section-number">6</div>
                <h2>Privacidad</h2>
              </div>
              <div class="section-content">
                <p>
                  El uso de nuestros servicios implica la aceptación de nuestras políticas de privacidad y
                  manejo de datos. Para más información sobre cómo gestionamos tus datos, consulta nuestra
                  <router-link to="/legal/privacy">Política de Privacidad</router-link>.
                </p>
              </div>
            </div>

            <!-- Sección 7: Legislación aplicable -->
            <div id="applicable-law" class="terms-section">
              <div class="section-header">
                <div class="section-number">7</div>
                <h2>Legislación aplicable</h2>
              </div>
              <div class="section-content">
                <p>
                  Estos términos se rigen por las leyes vigentes en el país donde se desarrolló
                  SignalSysAnalyzer. Cualquier disputa relacionada con estos términos estará sujeta a
                  la jurisdicción exclusiva de los tribunales competentes de dicha localidad.
                </p>
                <div class="law-container">
                  <div class="law-icon"><i class="bx bx-buildings"></i></div>
                  <p class="law-note">
                    En caso de conflicto, se intentará resolver mediante un proceso de mediación
                    antes de recurrir a medidas legales.
                  </p>
                </div>
              </div>
            </div>

            <!-- Sección 8: Aceptación de los términos -->
            <div id="terms-acceptance" class="terms-section">
              <div class="section-header">
                <div class="section-number">8</div>
                <h2>Aceptación de los términos</h2>
              </div>
              <div class="section-content">
                <div class="acceptance-container">
                  <div class="acceptance-icon"><i class="bx bx-check-circle"></i></div>
                  <p>
                    Al registrarte y utilizar nuestra plataforma, confirmas que has leído, entendido y
                    aceptado estos términos y condiciones. Si no estás de acuerdo con algún punto,
                    te recomendamos no utilizar SignalSysAnalyzer.
                  </p>
                </div>
                <div class="last-update">
                  <span class="update-label">Última actualización:</span>
                  <span class="update-date">Mayo 2025</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Navegación rápida -->
        <div class="terms-nav">
          <div class="nav-title">Navegación rápida</div>
          <div class="nav-links">
            <a class="nav-link" href="#service-usage">1. Uso del servicio</a>
            <a class="nav-link" href="#results-accuracy">2. Exactitud de los resultados</a>
            <a class="nav-link" href="#limited-liability">3. Responsabilidad limitada</a>
            <a class="nav-link" href="#intellectual-property">4. Propiedad intelectual</a>
            <a class="nav-link" href="#modifications">5. Modificaciones</a>
            <a class="nav-link" href="#privacy">6. Privacidad</a>
            <a class="nav-link" href="#applicable-law">7. Legislación aplicable</a>
            <a class="nav-link" href="#terms-acceptance">8. Aceptación de los términos</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss">
@import '/src/styles/variables.scss';

.terms-page {
  max-width: 1000px;
  margin: 0 auto;
  padding: 2rem 1rem 4rem;
  color: $text-color-light-mode;

  .dark-mode & {
    color: $text-color-dark-mode;
  }

  .terms-header {
    text-align: center;
    margin-bottom: 3rem;

    .page-title {
      font-family: $font-primary;
      font-size: 2.5rem;
      font-weight: 700;
      color: $primary-color-light-mode;
      margin-bottom: 0.5rem;
      position: relative;
      display: inline-block;

      &:after {
        content: "";
        position: absolute;
        bottom: -10px;
        left: 50%;
        transform: translateX(-50%);
        width: 80px;
        height: 3px;
        background-color: $primary-color-light-mode;
        border-radius: 3px;
      }

      .dark-mode & {
        color: $primary-color-dark-mode;

        &:after {
          background-color: $primary-color-dark-mode;
        }
      }
    }

    .subtitle {
      font-size: 1.1rem;
      color: $text-color-light-mode;
      opacity: 0.85;
      margin-top: 1.2rem;

      .dark-mode & {
        color: $text-color-dark-mode;
      }
    }
  }

  .terms-container {
    position: relative;

    .terms-intro {
      background-color: white;
      border-radius: 12px;
      box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
      padding: 2rem;
      margin-bottom: 2rem;
      display: flex;
      align-items: center;
      gap: 1.5rem;
      border-left: 4px solid $secondary-color-light-mode;

      .dark-mode & {
        background-color: #2a3238;
        box-shadow: 0 5px 20px rgba(0, 0, 0, 0.15);
        border-left: 4px solid $secondary-color-dark-mode;
      }

      .intro-icon {
        flex-shrink: 0;
        width: 64px;
        height: 64px;
        border-radius: 50%;
        background-color: rgba($secondary-color-light-mode, 0.1);
        display: flex;
        align-items: center;
        justify-content: center;

        i {
          font-size: 2.2rem;
          color: $secondary-color-light-mode;
        }

        .dark-mode & {
          background-color: rgba($secondary-color-dark-mode, 0.2);

          i {
            color: $secondary-color-dark-mode;
          }
        }
      }

      .intro-content {
        p {
          font-size: 1.1rem;
          line-height: 1.7;
          color: $text-color-light-mode;
          margin: 0;

          .dark-mode & {
            color: $text-color-dark-mode;
          }
        }
      }

      @media (max-width: 768px) {
        flex-direction: column;
        text-align: center;
        padding: 1.5rem;

        .intro-icon {
          margin-bottom: 1rem;
        }
      }
    }

    .terms-main {
      display: grid;
      grid-template-columns: 3fr 1fr;
      gap: 1.5rem;
      margin-bottom: 2rem;

      @media (max-width: 992px) {
        grid-template-columns: 1fr;
      }

      .terms-card {
        background-color: transparent;

        .terms-sections {
          .terms-section {
            background-color: white;
            border-radius: 12px;
            box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
            margin-bottom: 1.5rem;
            overflow: hidden;
            transition: transform 0.3s ease, box-shadow 0.3s ease;

            &:hover {
              transform: translateY(-3px);
              box-shadow: 0 8px 25px rgba(0, 0, 0, 0.08);
            }

            .dark-mode & {
              background-color: #2a3238;
              box-shadow: 0 5px 20px rgba(0, 0, 0, 0.15);

              &:hover {
                box-shadow: 0 8px 25px rgba(0, 0, 0, 0.25);
              }
            }

            .section-header {
              padding: 1.3rem 1.5rem;
              display: flex;
              align-items: center;
              gap: 1rem;
              background: linear-gradient(to right, #f8f9fa, white);
              border-bottom: 1px solid rgba(0, 0, 0, 0.05);

              .dark-mode & {
                background: linear-gradient(to right, #343e46, #2a3238);
                border-bottom: 1px solid rgba(255, 255, 255, 0.05);
              }

              .section-number {
                width: 38px;
                height: 38px;
                border-radius: 50%;
                background-color: $secondary-color-light-mode;
                color: white;
                font-weight: 600;
                font-size: 1.2rem;
                display: flex;
                align-items: center;
                justify-content: center;
                box-shadow: 0 3px 8px rgba($secondary-color-light-mode, 0.3);

                .dark-mode & {
                  background-color: $secondary-color-dark-mode;
                  box-shadow: 0 3px 8px rgba(0, 0, 0, 0.3);
                  color: #081925;
                }
              }

              h2 {
                margin: 0;
                font-size: 1.3rem;
                color: $secondary-color-light-mode;
                font-weight: 600;

                .dark-mode & {
                  color: $secondary-color-dark-mode;
                }
              }
            }

            .section-content {
              padding: 1.5rem;

              p {
                color: $text-color-light-mode;
                line-height: 1.6;
                margin-bottom: 1rem;

                &:last-child {
                  margin-bottom: 0;
                }

                .dark-mode & {
                  color: $text-color-dark-mode;
                }
              }

              a {
                color: $primary-color-light-mode;
                text-decoration: none;
                font-weight: 500;
                transition: color 0.2s ease;

                &:hover {
                  color: darken($primary-color-light-mode, 10%);
                  text-decoration: underline;
                }

                .dark-mode & {
                  color: $primary-color-dark-mode;

                  &:hover {
                    color: lighten($primary-color-dark-mode, 10%);
                  }
                }
              }

              .service-container,
              .accuracy-container,
              .law-container,
              .acceptance-container {
                display: flex;
                align-items: flex-start;
                gap: 1.2rem;

                .service-icon,
                .accuracy-icon,
                .law-icon,
                .acceptance-icon {
                  flex-shrink: 0;
                  width: 40px;
                  height: 40px;
                  border-radius: 50%;
                  background-color: rgba($secondary-color-light-mode, 0.1);
                  display: flex;
                  align-items: center;
                  justify-content: center;

                  i {
                    font-size: 1.5rem;
                    color: $secondary-color-light-mode;
                  }

                  .dark-mode & {
                    background-color: rgba($secondary-color-dark-mode, 0.2);

                    i {
                      color: $secondary-color-dark-mode;
                    }
                  }
                }

                @media (max-width: 768px) {
                  flex-direction: column;
                  align-items: center;
                  text-align: center;
                }
              }

              .liability-alert {
                display: flex;
                align-items: flex-start;
                gap: 1.2rem;
                background-color: rgba(#ff6b6b, 0.05);
                padding: 1.2rem;
                border-radius: 8px;
                border-left: 3px solid #ff6b6b;

                .dark-mode & {
                  background-color: rgba(#ff6b6b, 0.1);
                }

                .alert-icon {
                  flex-shrink: 0;
                  width: 40px;
                  height: 40px;
                  border-radius: 50%;
                  background-color: rgba(#ff6b6b, 0.15);
                  display: flex;
                  align-items: center;
                  justify-content: center;

                  i {
                    font-size: 1.5rem;
                    color: #ff6b6b;
                  }
                }

                .alert-content {
                  p {
                    margin: 0;
                  }
                }

                @media (max-width: 768px) {
                  flex-direction: column;
                  align-items: center;
                  text-align: center;
                }
              }

              .update-note {
                display: flex;
                align-items: flex-start;
                gap: 1.2rem;
                background-color: rgba($primary-color-light-mode, 0.05);
                padding: 1.2rem;
                border-radius: 8px;

                .dark-mode & {
                  background-color: rgba($primary-color-dark-mode, 0.1);
                }

                .update-icon {
                  flex-shrink: 0;
                  width: 40px;
                  height: 40px;
                  border-radius: 50%;
                  background-color: rgba($primary-color-light-mode, 0.15);
                  display: flex;
                  align-items: center;
                  justify-content: center;

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

                p {
                  margin: 0;
                }

                @media (max-width: 768px) {
                  flex-direction: column;
                  align-items: center;
                  text-align: center;
                }
              }

              .ip-examples {
                display: grid;
                grid-template-columns: repeat(2, 1fr);
                gap: 1rem;
                margin-top: 1.2rem;

                .ip-item {
                  display: flex;
                  align-items: center;
                  gap: 0.8rem;
                  padding: 0.8rem;
                  background-color: rgba(0, 0, 0, 0.02);
                  border-radius: 8px;

                  .dark-mode & {
                    background-color: rgba(255, 255, 255, 0.05);
                  }

                  .ip-icon {
                    flex-shrink: 0;
                    width: 32px;
                    height: 32px;
                    border-radius: 50%;
                    background-color: rgba($secondary-color-light-mode, 0.1);
                    display: flex;
                    align-items: center;
                    justify-content: center;

                    i {
                      font-size: 1.2rem;
                      color: $secondary-color-light-mode;
                    }

                    .dark-mode & {
                      background-color: rgba($secondary-color-dark-mode, 0.2);

                      i {
                        color: $secondary-color-dark-mode;
                      }
                    }
                  }

                  .ip-text {
                    font-size: 0.95rem;
                  }
                }

                @media (max-width: 576px) {
                  grid-template-columns: 1fr;
                }
              }

              .law-note {
                font-style: italic;
                margin: 0;
              }

              .last-update {
                margin-top: 1.5rem;
                padding: 0.8rem 1.2rem;
                background-color: rgba(0, 0, 0, 0.02);
                border-radius: 6px;
                display: inline-flex;
                align-items: center;

                .dark-mode & {
                  background-color: rgba(255, 255, 255, 0.05);
                }

                .update-label {
                  font-weight: 600;
                  margin-right: 0.5rem;
                }

                .update-date {
                  color: $secondary-color-light-mode;

                  .dark-mode & {
                    color: $secondary-color-dark-mode;
                  }
                }
              }
            }
          }
        }
      }

      .terms-nav {
        background-color: white;
        border-radius: 12px;
        box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
        padding: 1.5rem;
        position: sticky;
        top: 1rem;
        height: fit-content;

        .dark-mode & {
          background-color: #2a3238;
          box-shadow: 0 5px 20px rgba(0, 0, 0, 0.15);
        }

        .nav-title {
          font-weight: 600;
          font-size: 1.1rem;
          margin-bottom: 1.2rem;
          color: $secondary-color-light-mode;
          position: relative;
          padding-bottom: 0.8rem;

          &:after {
            content: "";
            position: absolute;
            bottom: 0;
            left: 0;
            width: 40px;
            height: 2px;
            background-color: $secondary-color-light-mode;
            border-radius: 2px;
          }

          .dark-mode & {
            color: $secondary-color-dark-mode;

            &:after {
              background-color: $secondary-color-dark-mode;
            }
          }
        }

        .nav-links {
          display: flex;
          flex-direction: column;
          gap: 0.3rem;

          .nav-link {
            padding: 0.6rem 0.8rem;
            border-radius: 6px;
            text-decoration: none;
            color: $text-color-light-mode;
            font-size: 0.95rem;
            transition: all 0.2s ease;
            border-left: 2px solid transparent;

            &:hover {
              background-color: rgba($secondary-color-light-mode, 0.05);
              color: $secondary-color-light-mode;
              border-left: 2px solid $secondary-color-light-mode;
            }

            .dark-mode & {
              color: $text-color-dark-mode;

              &:hover {
                background-color: rgba($secondary-color-dark-mode, 0.1);
                color: $secondary-color-dark-mode;
                border-left: 2px solid $secondary-color-dark-mode;
              }
            }
          }
        }

        @media (max-width: 992px) {
          display: none;
        }
      }
    }
  }
}
</style>
