package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
)

type Params struct {
	JWTSecret string `json:"jwtSecret"`
}

func loadOrCreateParams() (*Params, error) {
	paramsFilePath := "/opt/zen/data/params.json"

	if data, err := os.ReadFile(paramsFilePath); err == nil {
		var params Params
		if err := json.Unmarshal(data, &params); err == nil {
			return &params, nil
		}
	}

	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, err
	}

	params := &Params{
		JWTSecret: base64.StdEncoding.EncodeToString(secret),
	}

	paramsDir := filepath.Dir(paramsFilePath)
	if err := os.MkdirAll(paramsDir, 0755); err != nil {
		return nil, err
	}

	jsonData, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(paramsFilePath, jsonData, 0600); err != nil {
		return nil, err
	}

	return params, nil
}
