<template>
  <div class="mission-vision-page">
    <div class="hero-section">
      <div class="hero-content">
        <h1 class="page-title">Misión y Visión</h1>
        <p class="subtitle">Nuestro propósito y hacia dónde nos dirigimos</p>
      </div>
      <div class="hero-decoration">
        <div class="decoration-circle c1"></div>
        <div class="decoration-circle c2"></div>
        <div class="decoration-circle c3"></div>
      </div>
    </div>

    <div class="content-container">
      <div class="mission-vision-cards">
        <div class="mv-card mission-card">
          <div class="card-header">
            <div class="icon-container">
              <i class="bx bx-target-lock"></i>
            </div>
            <h2>Misión</h2>
          </div>
          <div class="card-content">
            <p>
              Desarrollar una plataforma inteligente basada en Machine Learning que permita a estudiantes e ingenieros
              analizar sistemas físicos de forma rápida, confiable y automatizada, facilitando la identificación
              de la función de transferencia a partir de datos experimentales.
            </p>
          </div>
          <div class="card-footer">
            <div class="key-elements">
              <div class="key-element">
                <div class="key-icon"><i class="bx bx-brain"></i></div>
                <span>Machine Learning</span>
              </div>
              <div class="key-element">
                <div class="key-icon"><i class="bx bx-line-chart"></i></div>
                <span>Análisis Rápido</span>
              </div>
              <div class="key-element">
                <div class="key-icon"><i class="bx bx-math"></i></div>
                <span>Función de Transferencia</span>
              </div>
            </div>
          </div>
        </div>

        <div class="mv-card vision-card">
          <div class="card-header">
            <div class="icon-container">
              <i class="bx bx-bulb"></i>
            </div>
            <h2>Visión</h2>
          </div>
          <div class="card-content">
            <p>
              Ser una herramienta de referencia en el área de análisis de sistemas, contribuyendo a la mejora
              del diseño de controladores y a la formación de futuros profesionales, integrando tecnologías de
              aprendizaje automático en entornos educativos e industriales.
            </p>
          </div>
          <div class="card-footer">
            <div class="key-elements">
              <div class="key-element">
                <div class="key-icon"><i class="bx bx-star"></i></div>
                <span>Referencia</span>
              </div>
              <div class="key-element">
                <div class="key-icon"><i class="bx bx-book-open"></i></div>
                <span>Formación</span>
              </div>
              <div class="key-element">
                <div class="key-icon"><i class="bx bx-cog"></i></div>
                <span>Integración</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="impact-section">
        <h2 class="section-title">Nuestro Impacto</h2>
        <div class="impact-grid">
          <div class="impact-item">
            <div class="impact-icon">
              <i class="bx bx-time-five"></i>
            </div>
            <h3>Ahorro de Tiempo</h3>
            <p>Automatizamos procesos que normalmente tomarían horas de cálculos manuales</p>
          </div>
          <div class="impact-item">
            <div class="impact-icon">
              <i class="bx bx-book-reader"></i>
            </div>
            <h3>Educación</h3>
            <p>Facilitamos el aprendizaje práctico en cursos de sistemas de control</p>
          </div>
          <div class="impact-item">
            <div class="impact-icon">
              <i class="bx bx-bar-chart-alt-2"></i>
            </div>
            <h3>Precisión</h3>
            <p>Ofrecemos resultados confiables basados en algoritmos avanzados</p>
          </div>
          <div class="impact-item">
            <div class="impact-icon">
              <i class="bx bx-network-chart"></i>
            </div>
            <h3>Conectando Industria</h3>
            <p>Acercamos la academia y la industria con herramientas profesionales</p>
          </div>
        </div>
      </div>

      <div class="future-section">
        <div class="future-content">
          <h2 class="section-title">Hacia el Futuro</h2>
          <p>
            En SignalSysAnalyzer estamos comprometidos con la constante evolución de nuestra plataforma.
            Trabajamos continuamente para implementar nuevas capacidades y refinar nuestros algoritmos de
            análisis.
          </p>
          <p>
            Nuestro objetivo es expandir las capacidades de nuestra herramienta para incluir sistemas más
            complejos y ofrecer recomendaciones automáticas para el diseño de controladores.
          </p>
          <div class="future-cta">
            <router-link class="cta-button" to="/">
              Comienza a usar SignalSysAnalyzer
            </router-link>
          </div>
        </div>
        <div class="future-decoration">
          <div class="waveform-container">
            <div class="waveform"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss">
@import '/src/styles/variables.scss';

.mission-vision-page {
  width: 100%;

  .hero-section {
    position: relative;
    background: linear-gradient(135deg, $primary-color-light-mode 0%, darken($primary-color-light-mode, 20%) 100%);
    padding: 4rem 2rem;
    text-align: center;
    overflow: hidden;

    .dark-mode & {
      background: linear-gradient(135deg, $secondary-color-light-mode 0%, darken($primary-color-dark-mode, 20%) 100%);
    }

    .hero-content {
      position: relative;
      z-index: 2;
      max-width: 800px;
      margin: 0 auto;

      .page-title {
        font-family: $font-primary;
        font-size: 3rem;
        font-weight: 700;
        color: white;
        margin-bottom: 1rem;
        text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
      }

      .subtitle {
        font-size: 1.2rem;
        color: rgba(255, 255, 255, 0.9);
        max-width: 600px;
        margin: 0 auto;
      }
    }

    .hero-decoration {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      z-index: 1;

      .decoration-circle {
        position: absolute;
        border-radius: 50%;
        opacity: 0.1;
        background-color: white;

        &.c1 {
          width: 300px;
          height: 300px;
          top: -100px;
          left: -100px;
        }

        &.c2 {
          width: 200px;
          height: 200px;
          bottom: -50px;
          right: 10%;
        }

        &.c3 {
          width: 150px;
          height: 150px;
          top: 20%;
          right: -50px;
        }
      }
    }
  }

  .content-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 3rem 1rem;
  }

  .mission-vision-cards {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    margin-top: -4rem;
    margin-bottom: 4rem;

    @media (min-width: 992px) {
      flex-direction: row;
    }

    .mv-card {
      flex: 1;
      background-color: white;
      border-radius: 10px;
      box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
      overflow: hidden;
      display: flex;
      flex-direction: column;

      .dark-mode & {
        background-color: #2a3238;
        box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
      }

      .card-header {
        padding: 2rem;
        display: flex;
        align-items: center;
        gap: 1rem;

        .icon-container {
          width: 60px;
          height: 60px;
          background-color: $primary-color-light-mode;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;

          i {
            font-size: 2rem;
            color: white;
          }

          .dark-mode & {
            background-color: $primary-color-dark-mode;
          }
        }

        h2 {
          font-size: 1.8rem;
          color: $secondary-color-light-mode;
          margin: 0;

          .dark-mode & {
            color: $secondary-color-dark-mode;
          }
        }
      }

      .card-content {
        padding: 0 2rem 2rem;
        flex-grow: 1;

        p {
          font-size: 1.1rem;
          line-height: 1.7;
          color: $text-color-light-mode;

          .dark-mode & {
            color: $text-color-dark-mode;
          }
        }
      }

      .card-footer {
        padding: 1.5rem 2rem;
        background-color: #f9f9f9;

        .dark-mode & {
          background-color: #343e46;
        }

        .key-elements {
          display: flex;
          flex-wrap: wrap;
          gap: 1rem;

          .key-element {
            display: flex;
            align-items: center;
            gap: 0.5rem;

            .key-icon {
              width: 30px;
              height: 30px;
              background-color: $primary-color-light-mode;
              border-radius: 50%;
              display: flex;
              align-items: center;
              justify-content: center;

              i {
                font-size: 0.9rem;
                color: white;
              }

              .dark-mode & {
                background-color: $primary-color-dark-mode;
              }
            }

            span {
              font-size: 0.9rem;
              font-weight: 500;
              color: $text-color-light-mode;

              .dark-mode & {
                color: $text-color-dark-mode;
              }
            }
          }
        }
      }

      &.mission-card {
        .icon-container {
          background-color: $primary-color-light-mode;

          .dark-mode & {
            background-color: $primary-color-dark-mode;
          }
        }
      }

      &.vision-card {
        .icon-container {
          background-color: $secondary-color-light-mode;

          .dark-mode & {
            background-color: $secondary-color-dark-mode;
          }
        }
      }
    }
  }

  .impact-section {
    margin: 4rem 0;

    .section-title {
      font-size: 2rem;
      text-align: center;
      margin-bottom: 3rem;
      color: $secondary-color-light-mode;
      position: relative;

      &:after {
        content: "";
        position: absolute;
        bottom: -15px;
        left: 50%;
        transform: translateX(-50%);
        width: 60px;
        height: 3px;
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

    .impact-grid {
      display: grid;
      grid-template-columns: repeat(1, 1fr);
      gap: 2rem;

      @media (min-width: 576px) {
        grid-template-columns: repeat(2, 1fr);
      }

      @media (min-width: 992px) {
        grid-template-columns: repeat(4, 1fr);
      }

      .impact-item {
        text-align: center;
        padding: 2rem;
        background-color: white;
        border-radius: 8px;
        box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
        transition: transform 0.3s ease;

        &:hover {
          transform: translateY(-10px);
        }

        .dark-mode & {
          background-color: #2a3238;
          box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
        }

        .impact-icon {
          width: 70px;
          height: 70px;
          background-color: #f1f7ff;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          margin: 0 auto 1.5rem;

          i {
            font-size: 2.5rem;
            color: $primary-color-light-mode;
          }

          .dark-mode & {
            background-color: #343e46;

            i {
              color: $primary-color-dark-mode;
            }
          }
        }

        h3 {
          font-size: 1.2rem;
          margin-bottom: 1rem;
          color: $secondary-color-light-mode;

          .dark-mode & {
            color: $secondary-color-dark-mode;
          }
        }

        p {
          font-size: 0.95rem;
          color: $text-color-light-mode;
          line-height: 1.6;

          .dark-mode & {
            color: $text-color-dark-mode;
          }
        }
      }
    }
  }

  .future-section {
    margin: 5rem 0 2rem;
    display: flex;
    flex-direction: column;
    background-color: #f9f9f9;
    border-radius: 10px;
    overflow: hidden;

    @media (min-width: 992px) {
      flex-direction: row;
    }

    .dark-mode & {
      background-color: #343e46;
    }

    .future-content {
      flex: 1;
      padding: 3rem;

      .section-title {
        font-size: 2rem;
        margin-bottom: 1.5rem;
        color: $secondary-color-light-mode;

        .dark-mode & {
          color: $secondary-color-dark-mode;
        }
      }

      p {
        font-size: 1.1rem;
        line-height: 1.7;
        margin-bottom: 1.5rem;
        color: $text-color-light-mode;

        .dark-mode & {
          color: $text-color-dark-mode;
        }

        &:last-of-type {
          margin-bottom: 2rem;
        }
      }

      .future-cta {
        margin-top: 2rem;

        .cta-button {
          display: inline-block;
          padding: 1rem 2rem;
          background-color: $primary-color-light-mode;
          color: white;
          font-weight: 600;
          border-radius: 30px;
          text-decoration: none;
          transition: all 0.3s ease;

          &:hover {
            background-color: darken($primary-color-light-mode, 10%);
            transform: translateY(-3px);
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
          }

          .dark-mode & {
            background-color: $secondary-color-dark-mode;
            color: $secondary-color-light-mode;

            &:hover {
              background-color: darken($primary-color-dark-mode, 10%);
              box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
            }
          }
        }
      }
    }

    .future-decoration {
      flex: 1;
      min-height: 300px;
      background-color: #f1f7ff;
      display: flex;
      align-items: center;
      justify-content: center;
      position: relative;

      .dark-mode & {
        background-color: #272d33;
      }

      .waveform-container {
        width: 80%;
        height: 150px;
        position: relative;
        overflow: hidden;

        .waveform {
          position: absolute;
          top: 50%;
          left: 0;
          width: 100%;
          height: 100px;
          transform: translateY(-50%);
          background: repeating-linear-gradient(
            to right,
            transparent,
            transparent 5px,
            $primary-color-light-mode 5px,
            $primary-color-light-mode 6px
          );
          mask-image: url("data:image/svg+xml,%3Csvg width='100%25' height='100%25' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M0,50 C100,0 200,100 300,50 C400,0 500,100 600,50 C700,0 800,100 900,50 C1000,0 1100,100 1200,50' stroke='black' fill='transparent' stroke-width='3'/%3E%3C/svg%3E");
          mask-size: 100% 100%;
          mask-repeat: no-repeat;

          .dark-mode & {
            background: repeating-linear-gradient(
              to right,
              transparent,
              transparent 5px,
              $primary-color-dark-mode 5px,
              $primary-color-dark-mode 6px
            );
          }
        }
      }
    }
  }
}
</style>
