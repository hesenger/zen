package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"
)

type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type JWTValidator interface {
	Validate(token string) (*Claims, error)
}

type osFileSystem struct{}

func (fs *osFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (fs *osFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (fs *osFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
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
	fs            FileSystem
	hasher        PasswordHasher
	validator     JWTValidator
	setupFilePath string
}

func NewSetupService(fs FileSystem, hasher PasswordHasher, validator JWTValidator, setupFilePath string) *SetupService {
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
