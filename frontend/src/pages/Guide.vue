<template>
  <div class="guide-page">
    <h1 class="page-title">Guía / Tutorial - SignalSys Analyzer</h1>

    <div class="guide-content">
      <div class="guide-section">
        <h2>Bienvenida</h2>
        <p>
          Bienvenido a SignalSysAnalyzer. Esta herramienta ha sido desarrollada para ayudarte a identificar
          el tipo de sistema y sus respectivos polos de un circuito físico, utilizando datos obtenidos
          experimentalmente con un osciloscopio.
        </p>
      </div>

      <div class="guide-section">
        <h2>Pasos para usar la herramienta</h2>

        <ol class="steps-list">
          <li>
            <h3>Preparar el archivo</h3>
            <p>
              Debes contar con un archivo en formato CSV, que contenga al menos 5.000 datos experimentales.
              El archivo debe tener los siguientes datos en una sola columna:
            </p>
            <ul>
              <li>Tiempo (t)</li>
              <li>Salida del sistema (y(t))</li>
            </ul>
            <div class="file-format-box">
              <h3>Formato CSV esperado</h3>
              <pre>
Waveform Data,
<span class="comentario">tiempo(t),salida(y(t)),</span>
-4.488000e+00,4.00e-01,
-4.487000e+00,4.00e-01,
-4.486000e+00,4.00e-01,
-4.485000e+00,8.00e-01,
-4.484000e+00,1.80e+00,
-4.483000e+00,3.00e+00,
-4.482000e+00,3.80e+00,
-4.481000e+00,4.60e+00,
...</pre>
            </div>
          </li>

          <li>
            <h3>Subir el archivo</h3>
            <p>
              Arrastra el archivo o selecciona uno desde tu equipo en la zona de carga.
            </p>
            <div class="step-image">
              <div class="image-placeholder">
                <img v-if="isDarkMode" alt="Subir archivo" src="../assets/upload-file-dark-mode.png" />
                <img v-else alt="Subir archivo" src="../assets/upload-file-light-mode.png" />
              </div>
            </div>
          </li>

          <li>
            <h3>Ingresar el voltaje de entrada</h3>
            <p>
              En la interfaz encontrarás un campo para colocar el valor del voltaje de entrada utilizado durante la toma de datos.
              Este valor no viene en el archivo, por lo tanto debe ser ingresado manualmente. El voltaje debe estar en el rango de 4 a 12V.
            </p>
          </li>

          <li>
            <h3>Presionar "Enviar"</h3>
            <p>
              Una vez cargado el archivo y el valor del voltaje, da clic en el botón "Enviar".
            </p>
          </li>

          <li>
            <h3>Ver resultados</h3>
            <p>
              Después del procesamiento, se mostrará:
            </p>
            <ul>
              <li>El tipo de sistema detectado (subamortiguado, sobreamortiguado, críticamente amortiguado) junto con una breve explicación.</li>
              <li>Los polos del sistema, según la función de transferencia estimada.</li>
              <li>Una gráfica representativa del comportamiento del sistema.</li>
            </ul>
            <div class="step-image">
              <div class="image-placeholder" style="height: 400px;">
                <img v-if="isDarkMode" alt="Resultados" src="../assets/results-dark-mode.png" style="height: 100%;" />
                <img v-else alt="Resultados" src="../assets/results-light-mode.png" style="height: 100%;" />
              </div>
            </div>
          </li>
        </ol>
      </div>

      <div class="guide-section">
        <h2>Información Adicional</h2>
        <p>
          Para obtener los mejores resultados:
        </p>
        <ul>
          <li>Asegúrate de que el archivo CSV tenga una frecuencia de muestreo adecuada</li>
          <li>El voltaje de entrada debe ser el valor exacto utilizado durante el experimento</li>
          <li>Es recomendable que los datos capturen la respuesta completa del sistema, desde el estado inicial hasta el estado estable</li>
        </ul>
      </div>

      <div class="guide-section">
        <h2>Soporte Técnico</h2>
        <p>
          Para asistencia adicional o reportes de problemas, contáctanos a través de:
        </p>
        <ul>
          <li class="contact-info">Email:
            <a href="mailto:tvega@unillanos.edu.co">
              tvega@unillanos.edu.co
            </a>
          </li>
          <li class="contact-info">Formulario web: <router-link to="/support/contact">Contacto</router-link></li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
  const isDarkMode = inject('isDarkMode');
</script>

<style lang="scss">
@import '../styles/variables.scss';

.guide-page {
  width: 100%;
  max-width: 900px;
  margin: 0 auto;
  padding: 0 1rem;

  .page-title {
    font-family: $font-primary;
    font-size: 2rem;
    font-weight: 700;
    margin-bottom: 2rem;
    text-align: center;
    color: $primary-color-light-mode;

    .dark-mode & {
      color: $primary-color-dark-mode;
    }
  }

  .guide-content {
    background-color: white;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
    padding: 2rem;

    .dark-mode & {
      background-color: #2a3238;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
    }
  }

  .guide-section {
    margin-bottom: 3rem;

    p {
      text-align: justify;
    }

    &:last-child {
      margin-bottom: 0;
    }

    h2 {
      font-size: 1.5rem;
      color: $secondary-color-light-mode;
      margin-bottom: 1rem;
      padding-bottom: 0.5rem;
      border-bottom: 1px solid #eee;

      .dark-mode & {
        color: $secondary-color-dark-mode;
        border-bottom-color: #444;
      }
    }

    h3 {
      font-size: 1.2rem;
      margin-top: 1.5rem;
      margin-bottom: 0.75rem;
      color: $text-color-light-mode;

      .dark-mode & {
        color: $text-color-dark-mode;
      }
    }

    p {
      margin-bottom: 1rem;
      line-height: 1.6;
    }

    ul, ol {
      margin-bottom: 1rem;
      padding-left: 1.5rem;

      li {
        margin-bottom: 0.5rem;
        line-height: 1.6;
      }
    }
  }

  .file-format-box {
    background-color: #f8f9fa;
    border-radius: 6px;
    padding: 1.5rem;
    margin-bottom: 1.5rem;

    .dark-mode & {
      background-color: #343e46;
    }

    h3 {
      color: $primary-color-light-mode;
      margin-top: 0;

      .dark-mode & {
        color: $primary-color-dark-mode;
      }
    }

    pre {
      background-color: #f1f1f1;
      border-radius: 4px;
      padding: 1rem;
      overflow-x: auto;
      font-family: monospace;
      margin: 1rem 0;

      .dark-mode & {
        background-color: #272d33;
        color: #e0e0e0;
      }
    }
  }

  .steps-list {
    counter-reset: steps-counter;
    list-style-type: none;
    padding-left: 0;

    .comentario {
      color: gray;
      font-style: italic;
    }

    li {
      position: relative;
      padding-left: 3rem;
      margin-bottom: 2rem;

      &:before {
        content: counter(steps-counter);
        counter-increment: steps-counter;
        position: absolute;
        left: 0;
        top: 0;
        width: 2rem;
        height: 2rem;
        background-color: $primary-color-light-mode;
        color: white;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: bold;

        .dark-mode & {
          background-color: $primary-color-dark-mode;
          color: $background-color-dark-mode;
        }
      }

      h3 {
        margin-top: 0;
      }

      .step-image {
        margin: 1rem 0;
        border: 1px dashed #ccc;
        border-radius: 4px;
        overflow: hidden;

        .dark-mode & {
          border-color: #555;
        }

        .image-placeholder {
          height: 300px;
          background-color: $background-color-light-mode;
          display: flex;
          align-items: center;
          justify-content: center;
          color: #aaa;

          .dark-mode & {
            background-color: $background-color-dark-mode;
            color: #888;
          }

          @media (max-width: 900px) {
            img {
              width: 100%;
              height: auto !important;
            }
          }
        }
      }
    }
    ul {
      li {
        padding-left: 0;

        &:before {
          content: none;
        }
      }
    }
  }

  .contact-info {
    padding: 0.5rem 0;
  }

  @media (max-width: 500px) {
  .steps-list li {
    padding-left: 0 !important;

    &:before {
      left: -40px !important;
    }
  }
}
}
</style>
