<template>
  <div class="system-chart">
    <canvas ref="chartCanvas"></canvas>
  </div>
</template>

<script setup>
  import { onMounted, onUnmounted, ref, watch } from 'vue';
  import Chart from 'chart.js/auto';

  const props = defineProps({
    graphData: {
      type: Object,
      required: true,
    },
    isDarkMode: {
      type: Boolean,
      default: false,
    },
  });

  const chartCanvas = ref(null);
  let chartInstance = null;

  const renderChart = () => {
    if (!chartCanvas.value) return;

    const ctx = chartCanvas.value.getContext('2d');

    if (chartInstance) {
      chartInstance.destroy();
    }

    const gridColor = props.isDarkMode ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)';
    const textColor = props.isDarkMode ? '#ffffff' : '#333333';

    // Crear el gráfico
    chartInstance = new Chart(ctx, {
      type: 'line',
      data: props.graphData,
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            position: 'top',
            labels: {
              color: textColor,
              font: {
                family: 'Libre Baskerville, serif',
                size: 12,
              },
            },
          },
          tooltip: {
            mode: 'index',
            intersect: false,
            backgroundColor: props.isDarkMode ? '#2a3238' : '#ffffff',
            titleColor: props.isDarkMode ? '#ffffff' : '#333333',
            bodyColor: props.isDarkMode ? '#dddddd' : '#666666',
            borderColor: props.isDarkMode ? '#444444' : '#eeeeee',
            borderWidth: 1,
          },
        },
        scales: {
          x: {
            title: {
              display: true,
              text: 'Tiempo (s)',
              color: textColor,
            },
            grid: {
              color: gridColor,
            },
            ticks: {
              color: textColor,
            },
          },
          y: {
            title: {
              display: true,
              text: 'Amplitud',
              color: textColor,
            },
            grid: {
              color: gridColor,
            },
            ticks: {
              color: textColor,
            },
          },
        },
        animation: {
          duration: 1000,
          easing: 'easeOutQuart',
        },
        elements: {
          point: {
            radius: 3,
            hoverRadius: 5,
          },
          line: {
            tension: 0.4,
          },
        },
      },
    });
  };

  onMounted(() => {
    renderChart();
  });

  watch(() => props.isDarkMode, () => {
    renderChart();
  });

  watch(() => props.graphData, () => {
    renderChart();
  }, { deep: true });

  onUnmounted(() => {
    if (chartInstance) {
      chartInstance.destroy();
      chartInstance = null;
    }
  });
</script>

<style lang="scss">
.system-chart {
  width: 100%;
  height: 100%;
  min-height: 300px;
}
</style>
