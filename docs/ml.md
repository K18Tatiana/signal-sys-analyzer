# Componente de Machine Learning

Este documento describe cómo se identifica el tipo de sistema y se estiman sus polos. Solo recoge lo que se puede comprobar en el código del repositorio (`backend/handlers/analysis.go`, `ml-api/main.py`) y en el notebook de entrenamiento (`Modelos_de_Go.ipynb`, ejecutado en Google Colab y **no incluido en este repositorio**).

## Alcance: qué reconoce y qué no

| Reconoce | No reconoce |
|---|---|
| Sistema **sobreamortiguado** (polos reales) | Sistema críticamente amortiguado |
| Sistema **subamortiguado** (polos complejos conjugados) | Sistemas de otro orden u otra topología |

Los datos de entrenamiento corresponden a la respuesta al escalón de **circuitos RLC serie** medidos con osciloscopio. El clasificador es binario: cualquier señal se asigna a una de las dos clases, aunque no sea un RLC serie.

## Flujo de datos

```
CSV del osciloscopio
   │  backend (Go): lectura del CSV, selección de la región útil, extractFeatures()
   ▼
30 características numéricas
   │  POST /predecir_tipo   {"datos": [30 valores]}  →  {"tipo_sistema": 0 | 1}
   │  POST /predecir_polos  {"datos": [30 valores]}  →  {"polo_s1_real", "polo_s1_imag", "polo_s2_real", "polo_s2_imag"}
   ▼
ml-api (Flask + scikit-learn): escala las características y aplica los modelos
```

`tipo_sistema`: `0` = subamortiguado, `1` = sobreamortiguado.

### Región útil de la señal

Antes de extraer las características, `detectUsefulRegion` recorta la señal: empieza unos 10 puntos antes de que la salida se separe de su valor inicial más de un 1 % del voltaje de entrada (mínimo 0.01 V) y toma como máximo 1500 puntos desde ahí.

### Las 30 características

Se calculan dos veces el mismo conjunto de 15 medidas: sobre el **voltaje de salida** (características 1–15) y sobre el **vector de tiempo** (16–30) de la región útil.

| Grupo | Medidas |
|---|---|
| Estadísticas básicas (7) | media, desviación estándar, mínimo, máximo, mediana, percentil 25, percentil 75 |
| Forma, sobre las diferencias entre muestras consecutivas (5) | suma de valores absolutos (variación total), máximo valor absoluto, media de valores absolutos, número de picos, desviación estándar |
| Adicionales (3) | rango (máximo − mínimo), energía promedio, media del valor absoluto |

## Modelos

Entrenados con scikit-learn (`Modelos_de_Go.ipynb`) y guardados con `joblib` en `ml-api/models/`.

| Modelo | Archivo | Algoritmo |
|---|---|---|
| Escalado | `scaler_X.pkl` | `StandardScaler` |
| Tipo de sistema | `modelo_tipo.pkl` | `RandomForestClassifier` (200 árboles, `max_depth=15`, `min_samples_split=5`, `min_samples_leaf=2`, `max_features='sqrt'`) |
| Polos sobreamortiguado (s1, s2) | `modelo_sobre_s1.pkl`, `modelo_sobre_s2.pkl` | Pipeline: `PolynomialFeatures(grado 2)` → `StandardScaler` → `RandomForestRegressor` (150 árboles, `max_depth=12`) |
| Polos subamortiguado (s1, s2) | `modelo_sub_s1.pkl`, `modelo_sub_s2.pkl` | `RandomForestRegressor` (200 árboles, `max_depth=15`) |

Cada regresor predice la parte real e imaginaria del polo. El servicio usa primero el clasificador y, según su resultado, aplica el par de regresores correspondiente.

## Datos de entrenamiento

- Archivo `datos_mejorados.csv`: **200 muestras** (100 sobreamortiguadas y 100 subamortiguadas) y 33 columnas (incluye `tipo`, `polo_s1`, `polo_s2` y las columnas `feature_*`).
- División 80/20 estratificada por tipo (`random_state=42`): 160 muestras de entrenamiento y 40 de prueba.
- Las grabaciones de osciloscopio disponibles del proyecto son de circuitos RLC serie con entradas de 4 V a 12 V (14 archivos, según sus nombres). No se conserva el procedimiento exacto con el que se generaron las 200 muestras a partir de ellas.

> `datos_mejorados.csv` **no está en este repositorio**, por lo que el entrenamiento no se puede reproducir desde aquí. Solo se incluyen los modelos ya entrenados.

## Resultados y limitaciones

Cifras que imprime el notebook, junto con su alcance real:

| Métrica reportada | Valor | Alcance |
|---|---|---|
| Accuracy del clasificador sobre el 20 % de prueba | 1.0000 (40 muestras) | Es la única medida con datos apartados del entrenamiento, pero el conjunto es pequeño |
| "Precisión en evaluación" | 1.0000 | Se calculó con 20 muestras tomadas de **todo** el dataset, que incluye datos de entrenamiento; **no** mide generalización |
| Error cuadrático medio de los polos, sobreamortiguado | 70 014 (s1) y 71 543 (s2) | Calculado sobre **todo** el dataset, incluidos los datos de entrenamiento |
| Error cuadrático medio de los polos, subamortiguado | 34 887 (s1 y s2) | Ídem |

Conclusiones prudentes:

- El 1.0 de accuracy se obtuvo con pocas muestras de un mismo tipo de circuito; no debe interpretarse como un rendimiento general.
- El error de los polos es grande respecto a su magnitud (del orden de centenas), y no se midió con datos independientes.
- **Ajuste manual de los polos:** tras recibir la predicción, el backend suma o resta constantes fijas antes de guardarla (subamortiguado: +127 a la parte real y ±335.25 a la imaginaria; sobreamortiguado: −425.8 a s1 y +455 a s2; ver `processAnalysisRequest` en `backend/handlers/analysis.go`). Es una corrección empírica añadida a mano, no parte del modelo entrenado.
- No hay conjunto de prueba independiente ni validación cruzada documentados.

Para mejorar el componente haría falta más grabaciones de circuitos distintos, una validación con datos no vistos (por ejemplo, por grabación completa) y, si se mantiene, justificar o eliminar el ajuste manual.
