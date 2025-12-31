package main

import (
	"errors"
	"os"
	"testing"
)

type mockFileSystemSetup struct {
	statFunc      func(name string) (os.FileInfo, error)
	writeFileFunc func(name string, data []byte, perm os.FileMode) error
	mkdirAllFunc  func(path string, perm os.FileMode) error
}

func (m *mockFileSystemSetup) Stat(name string) (os.FileInfo, error) {
	if m.statFunc != nil {
		return m.statFunc(name)
	}
	return nil, nil
}

func (m *mockFileSystemSetup) WriteFile(name string, data []byte, perm os.FileMode) error {
	if m.writeFileFunc != nil {
		return m.writeFileFunc(name, data, perm)
	}
	return nil
}

func (m *mockFileSystemSetup) MkdirAll(path string, perm os.FileMode) error {
	if m.mkdirAllFunc != nil {
		return m.mkdirAllFunc(path, perm)
	}
	return nil
}

func (m *mockFileSystemSetup) ReadFile(filename string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (m *mockFileSystemSetup) Create(name string) (*os.File, error) {
	return nil, errors.New("not implemented")
}

func (m *mockFileSystemSetup) Open(name string) (*os.File, error) {
	return nil, errors.New("not implemented")
}

func (m *mockFileSystemSetup) Remove(name string) error {
	return errors.New("not implemented")
}

func (m *mockFileSystemSetup) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return nil, errors.New("not implemented")
}

type mockPasswordHasher struct {
	hashFunc func(password string) (string, error)
}

func (m *mockPasswordHasher) Hash(password string) (string, error) {
	if m.hashFunc != nil {
		return m.hashFunc(password)
	}
	return "hashed_" + password, nil
}

type mockJWTValidator struct {
	validateFunc func(token string) (*Claims, error)
}

func (m *mockJWTValidator) Validate(token string) (*Claims, error) {
	if m.validateFunc != nil {
		return m.validateFunc(token)
	}
	return nil, nil
}

func TestCheckSetupStatus_PendingSetup(t *testing.T) {
	mockFS := &mockFileSystemSetup{
		statFunc: func(name string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		},
	}

	service := NewSetupService(mockFS, &mockPasswordHasher{}, &mockJWTValidator{}, "/test/setup.json")
	status, err := service.CheckSetupStatus("")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if status != StatusPendingSetup {
		t.Errorf("expected status %s, got %s", StatusPendingSetup, status)
	}
}

func TestCheckSetupStatus_Authenticated(t *testing.T) {
	mockFS := &mockFileSystemSetup{
		statFunc: func(name string) (os.FileInfo, error) {
			return nil, nil
		},
	}

	mockValidator := &mockJWTValidator{
		validateFunc: func(token string) (*Claims, error) {
			if token == "valid_token" {
				return &Claims{Username: "testuser"}, nil
			}
			return nil, errors.New("invalid token")
		},
	}

	service := NewSetupService(mockFS, &mockPasswordHasher{}, mockValidator, "/test/setup.json")
	status, err := service.CheckSetupStatus("valid_token")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if status != StatusAuthenticated {
		t.Errorf("expected status %s, got %s", StatusAuthenticated, status)
	}
}

func TestCheckSetupStatus_Ready(t *testing.T) {
	mockFS := &mockFileSystemSetup{
		statFunc: func(name string) (os.FileInfo, error) {
			return nil, nil
		},
	}

	mockValidator := &mockJWTValidator{
		validateFunc: func(token string) (*Claims, error) {
			return nil, errors.New("invalid token")
		},
	}

	service := NewSetupService(mockFS, &mockPasswordHasher{}, mockValidator, "/test/setup.json")
	status, err := service.CheckSetupStatus("invalid_token")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if status != StatusReady {
		t.Errorf("expected status %s, got %s", StatusReady, status)
	}
}

func TestCheckSetupStatus_ReadyNoToken(t *testing.T) {
	mockFS := &mockFileSystemSetup{
		statFunc: func(name string) (os.FileInfo, error) {
			return nil, nil
		},
	}

	service := NewSetupService(mockFS, &mockPasswordHasher{}, &mockJWTValidator{}, "/test/setup.json")
	status, err := service.CheckSetupStatus("")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if status != StatusReady {
		t.Errorf("expected status %s, got %s", StatusReady, status)
	}
}

func TestPerformSetup_Success(t *testing.T) {
	var writtenData []byte
	var writtenPath string

	mockFS := &mockFileSystemSetup{
		mkdirAllFunc: func(path string, perm os.FileMode) error {
			if perm != 0755 {
				t.Errorf("expected MkdirAll perm 0755, got %v", perm)
			}
			return nil
		},
		writeFileFunc: func(name string, data []byte, perm os.FileMode) error {
			writtenPath = name
			writtenData = data
			if perm != 0600 {
				t.Errorf("expected WriteFile perm 0600, got %v", perm)
			}
			return nil
		},
	}

	mockHasher := &mockPasswordHasher{
		hashFunc: func(password string) (string, error) {
			return "hashed_" + password, nil
		},
	}

	service := NewSetupService(mockFS, mockHasher, &mockJWTValidator{}, "/test/data/setup.json")

	setupData := SetupData{
		Username:    "testuser",
		Password:    "testpass",
		GithubToken: "ghtoken",
		Apps:        []App{{Provider: "openai", Key: "key1", Command: "cmd1"}},
	}

	err := service.PerformSetup(setupData)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if writtenPath != "/test/data/setup.json" {
		t.Errorf("expected file written to /test/data/setup.json, got %s", writtenPath)
	}

	if len(writtenData) == 0 {
		t.Error("expected data to be written, got empty data")
	}

	expectedPasswordSubstring := "hashed_testpass"
	if !contains(string(writtenData), expectedPasswordSubstring) {
		t.Errorf("expected written data to contain %s", expectedPasswordSubstring)
	}
}

func TestPerformSetup_HashError(t *testing.T) {
	mockHasher := &mockPasswordHasher{
		hashFunc: func(password string) (string, error) {
			return "", errors.New("hash failed")
		},
	}

	service := NewSetupService(&mockFileSystemSetup{}, mockHasher, &mockJWTValidator{}, "/test/setup.json")

	setupData := SetupData{
		Username: "testuser",
		Password: "testpass",
	}

	err := service.PerformSetup(setupData)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "hash failed" {
		t.Errorf("expected error 'hash failed', got %v", err)
	}
}

func TestPerformSetup_MkdirError(t *testing.T) {
	mockFS := &mockFileSystemSetup{
		mkdirAllFunc: func(path string, perm os.FileMode) error {
			return errors.New("mkdir failed")
		},
	}

	service := NewSetupService(mockFS, &mockPasswordHasher{}, &mockJWTValidator{}, "/test/setup.json")

	setupData := SetupData{
		Username: "testuser",
		Password: "testpass",
	}

	err := service.PerformSetup(setupData)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "mkdir failed" {
		t.Errorf("expected error 'mkdir failed', got %v", err)
	}
}

func TestPerformSetup_WriteFileError(t *testing.T) {
	mockFS := &mockFileSystemSetup{
		writeFileFunc: func(name string, data []byte, perm os.FileMode) error {
			return errors.New("write failed")
		},
	}

	service := NewSetupService(mockFS, &mockPasswordHasher{}, &mockJWTValidator{}, "/test/setup.json")

	setupData := SetupData{
		Username: "testuser",
		Password: "testpass",
	}

	err := service.PerformSetup(setupData)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "write failed" {
		t.Errorf("expected error 'write failed', got %v", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
