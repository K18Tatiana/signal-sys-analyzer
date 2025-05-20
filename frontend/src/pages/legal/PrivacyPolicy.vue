<template>
  <div class="privacy-page">
    <div class="privacy-header">
      <h1 class="page-title">Política de Privacidad</h1>
      <p class="subtitle">Protección y tratamiento de tus datos en SignalSysAnalyzer</p>
    </div>

    <div class="privacy-container">
      <!-- Introducción -->
      <div class="privacy-intro">
        <div class="intro-icon">
          <i class="bx bx-shield-quarter"></i>
        </div>
        <div class="intro-content">
          <p>
            En SignalSysAnalyzer, valoramos tu privacidad y nos comprometemos a proteger tus datos personales.
            Esta Política de Privacidad explica cómo recopilamos, utilizamos y resguardamos tu información cuando utilizas nuestra plataforma.
          </p>
        </div>
      </div>

      <div class="privacy-layout">
        <!-- Navegación rápida para escritorio -->
        <div class="quick-nav desktop-nav">
          <div class="nav-title">Navegación rápida</div>
          <div class="nav-links">
            <a class="nav-link" href="#data-collection">1. Qué datos recopilamos</a>
            <a class="nav-link" href="#data-usage">2. Cómo usamos tus datos</a>
            <a class="nav-link" href="#data-protection">3. Cómo protegemos tus datos</a>
            <a class="nav-link" href="#data-retention">4. Conservación de datos</a>
            <a class="nav-link" href="#user-rights">5. Tus derechos</a>
            <a class="nav-link" href="#cookies">6. Cookies</a>
            <a class="nav-link" href="#policy-changes">7. Cambios en esta política</a>
            <a class="nav-link" href="#contact">8. Contacto</a>
          </div>
        </div>

        <!-- Secciones principales -->
        <div class="privacy-sections">
          <!-- Sección 1: Qué datos recopilamos -->
          <div id="data-collection" class="policy-section">
            <div class="section-header">
              <div class="section-number">1</div>
              <h2>Qué datos recopilamos</h2>
            </div>
            <div class="section-content">
              <p>Recopilamos la siguiente información:</p>
              <ul class="data-list">
                <li><strong>Datos de usuario:</strong> nombre, correo electrónico y contraseña (encriptada).</li>
                <li><strong>Archivos cargados:</strong> CSV con datos experimentales (tiempo y salida del sistema).</li>
                <li><strong>Historial de análisis:</strong> tipo de sistema detectado, fecha y hora del análisis, resultados y gráficas.</li>
                <li><strong>Feedback:</strong> cualquier comentario o sugerencia que envíes voluntariamente.</li>
                <li><strong>Datos técnicos:</strong> tipo de navegador, sistema operativo, dirección IP y actividad dentro de la plataforma.</li>
              </ul>
            </div>
          </div>

          <!-- Sección 2: Cómo usamos tus datos -->
          <div id="data-usage" class="policy-section">
            <div class="section-header">
              <div class="section-number">2</div>
              <h2>Cómo usamos tus datos</h2>
            </div>
            <div class="section-content">
              <p>Utilizamos tu información para:</p>
              <ul class="data-list">
                <li>Procesar los análisis solicitados.</li>
                <li>Mejorar nuestros modelos de Machine Learning.</li>
                <li>Enviarte información relevante sobre tu cuenta o actualizaciones del sistema.</li>
                <li>Optimizar la experiencia del usuario dentro de la plataforma.</li>
                <li>Generar estadísticas internas de uso (nunca compartidas con terceros).</li>
              </ul>
            </div>
          </div>

          <!-- Sección 3: Cómo protegemos tus datos -->
          <div id="data-protection" class="policy-section">
            <div class="section-header">
              <div class="section-number">3</div>
              <h2>Cómo protegemos tus datos</h2>
            </div>
            <div class="section-content">
              <div class="protection-container">
                <div class="protection-icon"><i class="bx bx-lock-alt"></i></div>
                <ul class="data-list">
                  <li>Almacenamos los archivos y datos personales de forma <strong>encriptada</strong> en servidores seguros.</li>
                  <li>Solo el usuario que subió un archivo puede acceder a sus resultados.</li>
                  <li>No vendemos, alquilamos ni compartimos tus datos con terceros.</li>
                  <li>Limitamos el acceso interno a los datos únicamente al equipo técnico autorizado.</li>
                </ul>
              </div>
            </div>
          </div>

          <!-- Sección 4: Conservación de datos -->
          <div id="data-retention" class="policy-section">
            <div class="section-header">
              <div class="section-number">4</div>
              <h2>Conservación de datos</h2>
            </div>
            <div class="section-content">
              <ul class="data-list">
                <li>Conservamos los archivos y resultados analizados durante el tiempo que tu cuenta permanezca activa, salvo que solicites su eliminación.</li>
                <li>Los datos pueden ser utilizados de forma anónima para mejorar los algoritmos, sin asociarlos directamente a tu identidad.</li>
              </ul>
            </div>
          </div>

          <!-- Sección 5: Tus derechos -->
          <div id="user-rights" class="policy-section">
            <div class="section-header">
              <div class="section-number">5</div>
              <h2>Tus derechos</h2>
            </div>
            <div class="section-content">
              <p>Tienes derecho a:</p>
              <div class="rights-grid">
                <div class="right-item">
                  <div class="right-icon"><i class="bx bx-search-alt"></i></div>
                  <div class="right-text">Acceder a tus datos personales.</div>
                </div>
                <div class="right-item">
                  <div class="right-icon"><i class="bx bx-edit-alt"></i></div>
                  <div class="right-text">Solicitar la corrección o eliminación de tu cuenta y archivos.</div>
                </div>
                <div class="right-item">
                  <div class="right-icon"><i class="bx bx-x-circle"></i></div>
                  <div class="right-text">Retirar tu consentimiento en cualquier momento.</div>
                </div>
                <div class="right-item">
                  <div class="right-icon"><i class="bx bx-message-square-detail"></i></div>
                  <div class="right-text">Presentar una queja si consideras que no hemos respetado tu privacidad.</div>
                </div>
              </div>
              <p class="note">
                Puedes ejercer estos derechos escribiéndonos a nuestro correo de contacto o mediante el
                <router-link to="/support/contact">formulario en la sección "Contacto"</router-link>.
              </p>
            </div>
          </div>

          <!-- Sección 6: Cookies -->
          <div id="cookies" class="policy-section">
            <div class="section-header">
              <div class="section-number">6</div>
              <h2>Cookies</h2>
            </div>
            <div class="section-content">
              <div class="cookies-container">
                <div class="cookies-icon"><i class="bx bx-cookie"></i></div>
                <p>
                  SignalSysAnalyzer puede usar cookies esenciales para el funcionamiento del sitio. No utilizamos cookies de seguimiento ni de publicidad.
                  Puedes configurar tu navegador para bloquearlas, pero algunas funciones podrían dejar de estar disponibles.
                </p>
              </div>
            </div>
          </div>

          <!-- Sección 7: Cambios en esta política -->
          <div id="policy-changes" class="policy-section">
            <div class="section-header">
              <div class="section-number">7</div>
              <h2>Cambios en esta política</h2>
            </div>
            <div class="section-content">
              <p>
                Nos reservamos el derecho de actualizar esta Política de Privacidad en cualquier momento.
                Si realizamos cambios significativos, te informaremos por correo electrónico o mediante un aviso dentro de la plataforma.
              </p>
              <div class="last-update">
                <span class="update-label">Última actualización:</span>
                <span class="update-date">Mayo 2025</span>
              </div>
            </div>
          </div>

          <!-- Sección 8: Contacto -->
          <div id="contact" class="policy-section">
            <div class="section-header">
              <div class="section-number">8</div>
              <h2>Contacto</h2>
            </div>
            <div class="section-content">
              <div class="contact-container">
                <div class="contact-icon"><i class="bx bx-envelope"></i></div>
                <p>
                  Si tienes preguntas sobre esta política o deseas ejercer alguno de tus derechos, contáctanos a través del
                  <router-link to="/support/contact">formulario en la sección "Contacto"</router-link> o escríbenos a:
                  <a class="email-link" href="mailto:tvega@unillanos.edu.co">tvega@unillanos.edu.co</a>
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Navegación rápida para móviles -->
      <div class="quick-nav mobile-nav">
        <div class="nav-title">Navegación rápida</div>
        <div class="nav-links">
          <a class="nav-link" href="#data-collection">1. Qué datos recopilamos</a>
          <a class="nav-link" href="#data-usage">2. Cómo usamos tus datos</a>
          <a class="nav-link" href="#data-protection">3. Cómo protegemos tus datos</a>
          <a class="nav-link" href="#data-retention">4. Conservación de datos</a>
          <a class="nav-link" href="#user-rights">5. Tus derechos</a>
          <a class="nav-link" href="#cookies">6. Cookies</a>
          <a class="nav-link" href="#policy-changes">7. Cambios en esta política</a>
          <a class="nav-link" href="#contact">8. Contacto</a>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss">
@import '/src/styles/variables.scss';

.privacy-page {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 1rem 4rem;
  color: $text-color-light-mode;

  .dark-mode & {
    color: $text-color-dark-mode;
  }

  .privacy-header {
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

  .privacy-container {
    position: relative;

    .privacy-intro {
      background-color: white;
      border-radius: 12px;
      box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
      padding: 2rem;
      margin-bottom: 3rem;
      display: flex;
      align-items: center;
      gap: 1.5rem;
      border-left: 4px solid $primary-color-light-mode;

      .dark-mode & {
        background-color: #2a3238;
        box-shadow: 0 5px 20px rgba(0, 0, 0, 0.15);
        border-left: 4px solid $primary-color-dark-mode;
      }

      .intro-icon {
        flex-shrink: 0;
        width: 64px;
        height: 64px;
        border-radius: 50%;
        background-color: rgba($primary-color-light-mode, 0.1);
        display: flex;
        align-items: center;
        justify-content: center;

        i {
          font-size: 2.2rem;
          color: $primary-color-light-mode;
        }

        .dark-mode & {
          background-color: rgba($primary-color-dark-mode, 0.2);

          i {
            color: $primary-color-dark-mode;
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

    .privacy-layout {
      display: grid;
      grid-template-columns: 1fr 3fr;
      gap: 2rem;

      @media (max-width: 992px) {
        grid-template-columns: 1fr;
      }
    }

    .privacy-sections {
      margin-bottom: 3rem;

      .policy-section {
        background-color: white;
        border-radius: 12px;
        box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
        margin-bottom: 2rem;
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
          padding: 1.5rem 2rem;
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
            background-color: $primary-color-light-mode;
            color: white;
            font-weight: 600;
            font-size: 1.2rem;
            display: flex;
            align-items: center;
            justify-content: center;
            box-shadow: 0 3px 8px rgba($primary-color-light-mode, 0.3);

            .dark-mode & {
              background-color: $primary-color-dark-mode;
              box-shadow: 0 3px 8px rgb(0 0 0 / 30%);
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
          padding: 2rem;

          p, li {
            color: $text-color-light-mode;
            line-height: 1.7;

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

          .email-link {
            padding: 0.2rem 0.5rem;
            background-color: rgba($primary-color-light-mode, 0.1);
            border-radius: 4px;

            .dark-mode & {
              background-color: rgba($primary-color-dark-mode, 0.2);
            }
          }

          .note {
            font-style: italic;
            font-size: 0.95rem;
            margin-top: 1.2rem;
            padding: 0.8rem;
            background-color: rgba(0, 0, 0, 0.02);
            border-radius: 6px;

            .dark-mode & {
              background-color: rgba(255, 255, 255, 0.05);
            }
          }

          ul {
            padding-left: 1.5rem;

            li {
              margin-bottom: 0.6rem;

              &:last-child {
                margin-bottom: 0;
              }
            }
          }

          .data-list {
            list-style: none;
            padding-left: 0;

            li {
              margin-bottom: 0.8rem;
              padding-left: 1.8rem;
              position: relative;

              &:before {
                content: "•";
                position: absolute;
                left: 0;
                color: $primary-color-light-mode;
                font-weight: bold;
                font-size: 1.2rem;

                .dark-mode & {
                  color: $primary-color-dark-mode;
                }
              }

              &:last-child {
                margin-bottom: 0;
              }
            }
          }

          .protection-container,
          .cookies-container,
          .contact-container {
            display: flex;
            align-items: flex-start;
            gap: 1.5rem;

            .protection-icon,
            .cookies-icon,
            .contact-icon {
              flex-shrink: 0;
              width: 48px;
              height: 48px;
              border-radius: 50%;
              background-color: rgba($primary-color-light-mode, 0.1);
              display: flex;
              align-items: center;
              justify-content: center;

              i {
                font-size: 1.8rem;
                color: $primary-color-light-mode;
              }

              .dark-mode & {
                background-color: rgba($primary-color-dark-mode, 0.2);

                i {
                  color: $primary-color-dark-mode;
                }
              }
            }

            @media (max-width: 768px) {
              flex-direction: column;
              align-items: center;
              text-align: center;

              ul.data-list {
                text-align: left;
              }
            }
          }

          .rights-grid {
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 1.2rem;
            margin: 1.5rem 0;

            .right-item {
              display: flex;
              align-items: flex-start;
              gap: 0.8rem;
              padding: 1rem;
              background-color: rgba(0, 0, 0, 0.02);
              border-radius: 8px;
              transition: background-color 0.2s ease;

              &:hover {
                background-color: rgba($primary-color-light-mode, 0.05);
              }

              .dark-mode & {
                background-color: rgba(255, 255, 255, 0.05);

                &:hover {
                  background-color: rgba($primary-color-dark-mode, 0.1);
                }
              }

              .right-icon {
                flex-shrink: 0;
                width: 36px;
                height: 36px;
                border-radius: 50%;
                background-color: rgba($primary-color-light-mode, 0.1);
                display: flex;
                align-items: center;
                justify-content: center;

                i {
                  font-size: 1.4rem;
                  color: $primary-color-light-mode;
                }

                .dark-mode & {
                  background-color: rgba($primary-color-dark-mode, 0.2);

                  i {
                    color: $primary-color-dark-mode;
                  }
                }
              }

              .right-text {
                font-size: 0.95rem;
                line-height: 1.5;
              }
            }

            @media (max-width: 768px) {
              grid-template-columns: 1fr;
            }
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
              color: $primary-color-light-mode;

              .dark-mode & {
                color: $primary-color-dark-mode;
              }
            }
          }
        }
      }
    }

    .quick-nav {
      background-color: white;
      border-radius: 12px;
      box-shadow: 0 5px 20px rgba(0, 0, 0, 0.05);
      padding: 1.5rem;

      &.desktop-nav {
        position: sticky;
        top: 1rem;
        height: fit-content;

        @media (max-width: 992px) {
          display: none;
        }
      }

      &.mobile-nav {
        display: none;
        margin-bottom: 3rem;

        @media (max-width: 992px) {
          display: block;
        }
      }

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
          background-color: $primary-color-light-mode;
          border-radius: 2px;
        }

        .dark-mode & {
          color: $secondary-color-dark-mode;

          &:after {
            background-color: $primary-color-dark-mode;
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
            background-color: rgba($primary-color-light-mode, 0.05);
            color: $primary-color-light-mode;
            border-left: 2px solid $primary-color-light-mode;
          }

          .dark-mode & {
            color: $text-color-dark-mode;

            &:hover {
              background-color: rgba($primary-color-dark-mode, 0.1);
              color: $primary-color-dark-mode;
              border-left: 2px solid $primary-color-dark-mode;
            }
          }
        }
      }
    }
  }
}
</style>
