# Formato del CSV de entrada

La aplicación analiza archivos `.csv` exportados de un osciloscopio con la **respuesta al escalón** de un circuito (los de este proyecto son RLC serie). Tamaño máximo: 10 MB.

## Archivo de ejemplo

[`ejemplo.csv`](ejemplo.csv) es una captura real de un **circuito RLC serie sobreamortiguado** con **8 V** de entrada (~2 MB, 70 001 filas de datos: de 0 a 7 s cada 100 µs, más una última fila con marca de 8 s). Para analizarlo, súbelo en la página de inicio e indica **8** en el campo "voltaje de entrada" (el voltaje no viene en el archivo, se escribe en el formulario).

## Estructura del archivo

```
Format,1.0B,
Memory Length,10000,
...                      ← 27 líneas de metadatos del osciloscopio (clave,valor,)
Source,CH2,              ← canal medido
Vertical Units,V,
...
Sampling Period,5.000E-04,
...
Firmware,V1.19,
Mode,Detail,
Waveform Data,           ← última línea de metadatos (línea 27)
0.000000e+00,1.014670e-03,    ← datos: tiempo (s), voltaje de salida (V)
1.000000e-04,6.792557e-03,
...
```

| Parte | Contenido |
|---|---|
| Líneas 1 a 27 | Metadatos del osciloscopio (formato, canal, escalas, periodo de muestreo, firmware...). |
| Desde la línea 28 | Dos columnas numéricas separadas por coma, en notación científica: **columna 1 = tiempo en segundos**, **columna 2 = voltaje de salida en voltios**. Cada línea termina con una coma. |

## Cómo lo lee el backend

Según `processAnalysisRequest` (`backend/handlers/analysis.go`):

- Se leen todas las filas con `encoding/csv`; las filas cuya primera y segunda columna **no son números** (los metadatos) simplemente se descartan.
- La columna 1 se toma como tiempo y la columna 2 como voltaje de salida.
- El periodo de muestreo se busca en las primeras 10 líneas (`Sampling Period`). En este formato aparece en la línea 22, así que **no se usa**: se calcula como la diferencia entre las dos primeras marcas de tiempo (en el ejemplo, 1e-4 s = 100 µs). El valor `Sampling Period` del encabezado (5e-4) es un dato de configuración del osciloscopio y no coincide con el paso real de los datos.
- Si no se encuentra ninguna fila numérica válida, el análisis no genera resultados.

## Si usas tus propios datos

- Exporta la señal a CSV con dos columnas `tiempo,voltaje` (los metadatos son opcionales).
- Usa el punto como separador decimal y la coma como separador de columnas.
- Registra el escalón completo: desde el reposo hasta que la salida se estabiliza.
- Los modelos solo conocen circuitos RLC serie sobre y subamortiguados; consulta [ml.md](ml.md) para ver las limitaciones.
