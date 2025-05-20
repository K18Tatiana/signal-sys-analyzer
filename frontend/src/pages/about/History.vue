<template>
  <div ref="storySection" class="our-story-section">
    <div class="story-container">
      <h2 class="section-title">Nuestra historia</h2>

      <div class="story-content">
        <div class="story-image">
          <div class="image-placeholder">
            <i class="bx bx-line-chart story-icon"></i>
          </div>
        </div>

        <div class="story-text">
          <p class="story-paragraph">
            SignalSysAnalyzer nace de la iniciativa de tres estudiantes de Ingeniería Electrónica
            de la Universidad de los Llanos, con el propósito de facilitar el análisis y diseño de
            sistemas de control, brindando una herramienta sencilla, accesible y precisa.
          </p>

          <p class="story-paragraph">
            Nos motivó la necesidad de automatizar parte del proceso de identificación de sistemas físicos,
            especialmente en entornos educativos y profesionales.
          </p>

          <div class="story-values">
            <div class="value-item">
              <div class="value-icon">
                <i class="bx bx-bulb"></i>
              </div>
              <h3>Innovación</h3>
              <p>Buscamos nuevas formas de simplificar procesos complejos</p>
            </div>

            <div class="value-item">
              <div class="value-icon">
                <i class="bx bx-book-open"></i>
              </div>
              <h3>Educación</h3>
              <p>Facilitamos el aprendizaje y aplicación de conceptos de control</p>
            </div>

            <div class="value-item">
              <div class="value-icon">
                <i class="bx bx-target-lock"></i>
              </div>
              <h3>Precisión</h3>
              <p>Nos comprometemos con resultados confiables y exactos</p>
            </div>
          </div>
        </div>
      </div>

      <div class="team-glimpse">
        <h3 class="team-title">Desarrollado por estudiantes para estudiantes</h3>
        <div class="team-members">
          <div class="team-member">
            <div class="member-avatar">
              <i class="bx bx-user"></i>
            </div>
            <div class="member-info">
              <h4>Cristian Gerardo Alvarado Ramirez</h4>
              <p>Ingeniería Electrónica</p>
            </div>
          </div>

          <div class="team-member">
            <div class="member-avatar">
              <i class="bx bx-user"></i>
            </div>
            <div class="member-info">
              <h4>William David Avila Ducon</h4>
              <p>Ingeniería Electrónica</p>
            </div>
          </div>

          <div class="team-member">
            <div class="member-avatar">
              <i class="bx bx-user"></i>
            </div>
            <div class="member-info">
              <h4>Tatiana Vega Cardenas</h4>
              <p>Ingeniería Electrónica</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


<style lang="scss">
@import '/src/styles/variables.scss';

.our-story-section {
  width: 100%;
  padding: 4rem 1rem;
  background-color: #f9f9f9;

  .dark-mode & {
    background-color: #232a30;
  }

  .story-container {
    max-width: 1200px;
    margin: 0 auto;
  }

  .section-title {
    font-family: $font-primary;
    font-size: 2.5rem;
    font-weight: 700;
    text-align: center;
    margin-bottom: 3rem;
    color: $primary-color-light-mode;
    position: relative;

    &:after {
      content: "";
      position: absolute;
      bottom: -15px;
      left: 50%;
      transform: translateX(-50%);
      width: 80px;
      height: 4px;
      background-color: $primary-color-light-mode;
      border-radius: 2px;
    }

    .dark-mode & {
      color: $primary-color-dark-mode;

      &:after {
        background-color: $primary-color-dark-mode;
      }
    }
  }

  .story-content {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    margin-bottom: 3rem;

    @media (min-width: 768px) {
      flex-direction: row;
      align-items: center;
    }
  }

  .story-image {
    flex: 1;

    .image-placeholder {
      height: 300px;
      background-color: #eee;
      border-radius: 8px;
      display: flex;
      align-items: center;
      justify-content: center;

      .dark-mode & {
        background-color: #2a3238;
      }

      .story-icon {
        font-size: 6rem;
        color: $primary-color-light-mode;

        .dark-mode & {
          color: $primary-color-dark-mode;
        }
      }
    }
  }

  .story-text {
    flex: 1;

    .story-paragraph {
      font-size: 1.1rem;
      line-height: 1.8;
      margin-bottom: 1.5rem;
      color: $text-color-light-mode;

      .dark-mode & {
        color: $text-color-dark-mode;
      }

      &:first-of-type {
        font-weight: 500;
      }
    }
  }

  .story-values {
    display: grid;
    grid-template-columns: 1fr;
    gap: 1.5rem;
    margin-top: 2rem;

    @media (min-width: 576px) {
      grid-template-columns: repeat(3, 1fr);
    }

    .value-item {
      text-align: center;
      padding: 1.5rem 1rem;
      background-color: white;
      border-radius: 8px;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
      transition: transform 0.3s ease;

      &:hover {
        transform: translateY(-5px);
      }

      .dark-mode & {
        background-color: #2a3238;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
      }

      .value-icon {
        font-size: 2.5rem;
        color: $primary-color-light-mode;
        margin-bottom: 1rem;

        .dark-mode & {
          color: $primary-color-dark-mode;
        }
      }

      h3 {
        font-size: 1.2rem;
        margin-bottom: 0.5rem;
        color: $secondary-color-light-mode;

        .dark-mode & {
          color: $secondary-color-dark-mode;
        }
      }

      p {
        font-size: 0.9rem;
        color: $text-color-light-mode;

        .dark-mode & {
          color: $text-color-dark-mode;
        }
      }
    }
  }

  .team-glimpse {
    background-color: white;
    border-radius: 8px;
    padding: 2rem;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);

    .dark-mode & {
      background-color: #2a3238;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
    }

    .team-title {
      font-size: 1.5rem;
      text-align: center;
      margin-bottom: 2rem;
      color: $secondary-color-light-mode;

      .dark-mode & {
        color: $secondary-color-dark-mode;
      }
    }

    .team-members {
      display: flex;
      flex-wrap: wrap;
      justify-content: center;
      gap: 2rem;

      .team-member {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 180px;

        .member-avatar {
          width: 80px;
          height: 80px;
          background-color: #eee;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-bottom: 1rem;

          .dark-mode & {
            background-color: #343e46;
          }

          i {
            font-size: 2rem;
            color: $primary-color-light-mode;

            .dark-mode & {
              color: $primary-color-dark-mode;
            }
          }
        }

        .member-info {
          text-align: center;

          h4 {
            font-size: 1rem;
            margin-bottom: 0.25rem;
            color: $text-color-light-mode;

            .dark-mode & {
              color: $text-color-dark-mode;
            }
          }

          p {
            font-size: 0.9rem;
            color: $text-color-light-mode;
            opacity: 0.8;

            .dark-mode & {
              color: $text-color-dark-mode;
            }
          }
        }
      }
    }
  }
}
</style>
