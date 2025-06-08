package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MLClient struct {
	BaseURL string
	Client  *http.Client
}

type MLTypeRequest struct {
	Datos []float64 `json:"datos"`
}

type MLPolosRequest struct {
	Datos []float64 `json:"datos"`
}

type MLTypeResponse struct {
	TipoSistema int `json:"tipo_sistema"`
}

type MLPolosResponse struct {
	PoloS1Real float64 `json:"polo_s1_real"`
	PoloS1Imag float64 `json:"polo_s1_imag"`
	PoloS2Real float64 `json:"polo_s2_real"`
	PoloS2Imag float64 `json:"polo_s2_imag"`
}

func NewMLClient(baseURL string) *MLClient {
	return &MLClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *MLClient) PredictType(features []float64) (*MLTypeResponse, error) {
	reqBody := MLTypeRequest{Datos: features}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := c.Client.Post(c.BaseURL+"/predecir_tipo", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned status: %d", resp.StatusCode)
	}

	var result MLTypeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

func (c *MLClient) PredictPolos(features []float64) (*MLPolosResponse, error) {
	// CAMBIO PRINCIPAL: Usar directamente las características, no como matriz
	reqBody := MLPolosRequest{Datos: features}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := c.Client.Post(c.BaseURL+"/predecir_polos", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned status: %d", resp.StatusCode)
	}

	var result MLPolosResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}
