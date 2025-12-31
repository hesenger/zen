package main

import (
	"encoding/json"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type JWTValidator interface {
	Validate(token string) (*Claims, error)
}

type bcryptHasher struct{}

func (h *bcryptHasher) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

type jwtValidatorImpl struct{}

func (v *jwtValidatorImpl) Validate(token string) (*Claims, error) {
	return validateJWT(token)
}

type SetupService struct {
	fs            FileSystemOps
	hasher        PasswordHasher
	validator     JWTValidator
	setupFilePath string
}

func NewSetupService(fs FileSystemOps, hasher PasswordHasher, validator JWTValidator, setupFilePath string) *SetupService {
	return &SetupService{
		fs:            fs,
		hasher:        hasher,
		validator:     validator,
		setupFilePath: setupFilePath,
	}
}

func NewDefaultSetupService() *SetupService {
	return NewSetupService(
		&osFileSystem{},
		&bcryptHasher{},
		&jwtValidatorImpl{},
		"/opt/zen/data/setup.json",
	)
}

type CheckStatus string

const (
	StatusPendingSetup  CheckStatus = "pending-setup"
	StatusAuthenticated CheckStatus = "authenticated"
	StatusReady         CheckStatus = "ready"
)

func (s *SetupService) CheckSetupStatus(authToken string) (CheckStatus, error) {
	if _, err := s.fs.Stat(s.setupFilePath); err != nil {
		return StatusPendingSetup, nil
	}

	if authToken != "" {
		if _, err := s.validator.Validate(authToken); err == nil {
			return StatusAuthenticated, nil
		}
	}

	return StatusReady, nil
}

func (s *SetupService) PerformSetup(setupData SetupData) error {
	hashedPassword, err := s.hasher.Hash(setupData.Password)
	if err != nil {
		return err
	}
	setupData.Password = hashedPassword

	setupDir := filepath.Dir(s.setupFilePath)
	if err := s.fs.MkdirAll(setupDir, 0755); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(setupData, "", "  ")
	if err != nil {
		return err
	}

	if err := s.fs.WriteFile(s.setupFilePath, jsonData, 0600); err != nil {
		return err
	}

	return nil
}
