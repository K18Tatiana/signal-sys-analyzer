from flask import Flask, request, jsonify
import joblib
import numpy as np
import os

app = Flask(__name__)

class ModeloPolosEspecializado:
    """Clase para manejar el modelo especializado en la API"""
    def __init__(self):
        self.modelo_tipo = None
        self.scaler_X = None
        self.modelo_sobre_s1 = None
        self.modelo_sobre_s2 = None
        self.modelo_sub_s1 = None
        self.modelo_sub_s2 = None
        self.esta_entrenado = False
    
    def cargar_modelos(self, carpeta_modelos="models"):
        """Carga modelos previamente entrenados"""
        try:
            self.modelo_tipo = joblib.load(f"{carpeta_modelos}/modelo_tipo.pkl")
            self.scaler_X = joblib.load(f"{carpeta_modelos}/scaler_X.pkl")
            self.modelo_sobre_s1 = joblib.load(f"{carpeta_modelos}/modelo_sobre_s1.pkl")
            self.modelo_sobre_s2 = joblib.load(f"{carpeta_modelos}/modelo_sobre_s2.pkl")
            self.modelo_sub_s1 = joblib.load(f"{carpeta_modelos}/modelo_sub_s1.pkl")
            self.modelo_sub_s2 = joblib.load(f"{carpeta_modelos}/modelo_sub_s2.pkl")
            
            self.esta_entrenado = True
            print("Modelos especializados cargados exitosamente")
            return True
        except Exception as e:
            print(f"Error al cargar modelos: {e}")
            return False
    
    def predecir(self, X):
        """Predice tanto el tipo como los polos"""
        if not self.esta_entrenado:
            raise ValueError("El modelo no ha sido entrenado aún")
        
        # Asegurar que X sea 2D
        if X.ndim == 1:
            X = X.reshape(1, -1)
        
        # Escalar características
        X_scaled = self.scaler_X.transform(X)
        
        # Predecir tipo
        tipo_pred = self.modelo_tipo.predict(X_scaled)
        
        resultados = []
        
        for i in range(len(X_scaled)):
            if tipo_pred[i] == 1:  # Sobreamortiguada
                polo_s1 = self.modelo_sobre_s1.predict(X_scaled[i:i+1])[0]
                polo_s2 = self.modelo_sobre_s2.predict(X_scaled[i:i+1])[0]
                tipo_str = "sobre"
            else:  # Subamortiguada
                polo_s1 = self.modelo_sub_s1.predict(X_scaled[i:i+1])[0]
                polo_s2 = self.modelo_sub_s2.predict(X_scaled[i:i+1])[0]
                tipo_str = "sub"
            
            resultados.append({
                'tipo': tipo_str,
                'tipo_int': int(tipo_pred[i]),
                'polo_s1_real': float(polo_s1[0]),
                'polo_s1_imag': float(polo_s1[1]),
                'polo_s2_real': float(polo_s2[0]),
                'polo_s2_imag': float(polo_s2[1])
            })
        
        return resultados[0] if len(resultados) == 1 else resultados

# Cargar el modelo especializado al iniciar
modelo_especializado = ModeloPolosEspecializado()
if not modelo_especializado.cargar_modelos("models"):
    print("FALLA CRÍTICA: No se pudieron cargar los modelos")
    exit(1)

@app.route("/predecir_polos", methods=["POST"])
def predecir_polos():
    try:
        datos = request.json["datos"]  # Lista de características
        X = np.array([datos])  # Convertir a array 2D
        
        resultado = modelo_especializado.predecir(X)
        
        # Mantener EXACTAMENTE el mismo formato que la API anterior
        return jsonify({
            "polo_s1_real": resultado['polo_s1_real'], 
            "polo_s1_imag": resultado['polo_s1_imag'],
            "polo_s2_real": resultado['polo_s2_real'], 
            "polo_s2_imag": resultado['polo_s2_imag']
        })
            
    except Exception as e:
        return jsonify({"error": str(e)}), 400

@app.route("/predecir_tipo", methods=["POST"])
def predecir_tipo():
    try:
        datos = request.json["datos"]
        X = np.array([datos])
        
        resultado = modelo_especializado.predecir(X)
        
        # Mantener EXACTAMENTE el mismo formato que la API anterior
        return jsonify({
            "tipo_sistema": resultado['tipo_int']
        })
            
    except Exception as e:
        return jsonify({"error": str(e)}), 400

@app.route("/health", methods=["GET"])
def health_check():
    """Endpoint para verificar el estado de la API"""
    return jsonify({
        "status": "ok",
        "modelo_cargado": modelo_especializado.esta_entrenado
    })

@app.route("/", methods=["GET"])
def home():
    return jsonify({
        "mensaje": "API para predicción de polos en sistemas RLC",
        "version": "2.0 - Modelo Especializado Optimizado",
        "endpoints": [
            "/predecir_polos - POST: Predice polos de una señal",
            "/predecir_tipo - POST: Predice tipo de sistema",
            "/health - GET: Estado de la API"
        ]
    })

if __name__ == "__main__":
    app.run(debug=True, host="0.0.0.0", port=5001)