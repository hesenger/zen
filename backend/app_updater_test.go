package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
)

type mockFileSystemUpdater struct {
	files       map[string][]byte
	directories map[string]bool
	statError   error
}

func newMockFileSystem() *mockFileSystemUpdater {
	return &mockFileSystemUpdater{
		files:       make(map[string][]byte),
		directories: make(map[string]bool),
	}
}

func (m *mockFileSystemUpdater) ReadFile(filename string) ([]byte, error) {
	if data, ok := m.files[filename]; ok {
		return data, nil
	}
	return nil, os.ErrNotExist
}

func (m *mockFileSystemUpdater) WriteFile(filename string, data []byte, perm os.FileMode) error {
	m.files[filename] = data
	return nil
}

func (m *mockFileSystemUpdater) MkdirAll(path string, perm os.FileMode) error {
	m.directories[path] = true
	return nil
}

func (m *mockFileSystemUpdater) Stat(name string) (os.FileInfo, error) {
	if m.statError != nil {
		return nil, m.statError
	}
	if _, ok := m.directories[name]; ok {
		return nil, nil
	}
	return nil, os.ErrNotExist
}

func (m *mockFileSystemUpdater) Create(name string) (*os.File, error) {
	return nil, errors.New("not implemented in mock")
}

func (m *mockFileSystemUpdater) Open(name string) (*os.File, error) {
	return nil, errors.New("not implemented in mock")
}

func (m *mockFileSystemUpdater) Remove(name string) error {
	delete(m.files, name)
	return nil
}

func (m *mockFileSystemUpdater) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return nil, errors.New("not implemented in mock")
}

type mockCommandExecutor struct {
	executed []struct {
		command string
		workDir string
	}
}

func (m *mockCommandExecutor) Run(command, workDir string) error {
	m.executed = append(m.executed, struct {
		command string
		workDir string
	}{command, workDir})
	return nil
}

type mockArchiveExtractor struct{}

func (m *mockArchiveExtractor) ExtractTarGz(archivePath, destPath string) error {
	return nil
}

func (m *mockArchiveExtractor) ExtractZip(archivePath, destPath string) error {
	return nil
}

type mockGitHubDownloader struct {
	release      *GitHubRelease
	releaseError error
	downloadData []byte
}

func (m *mockGitHubDownloader) GetLatestRelease(repo, token string) (*GitHubRelease, error) {
	if m.releaseError != nil {
		return nil, m.releaseError
	}
	return m.release, nil
}

func (m *mockGitHubDownloader) DownloadAsset(url, token string) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(m.downloadData)), nil
}

type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestLoadSetupData(t *testing.T) {
	fs := newMockFileSystem()
	setupData := SetupData{
		Username:    "test",
		Password:    "hashed",
		GithubToken: "token123",
		Apps: []App{
			{Provider: "github", Key: "test/repo", Command: "echo test"},
		},
	}

	data, _ := json.Marshal(setupData)
	fs.files["/opt/zen/data/setup.json"] = data

	updater := NewAppUpdater(
		"/opt/zen/data/setup.json",
		fs,
		&mockCommandExecutor{},
		&mockArchiveExtractor{},
		&mockGitHubDownloader{},
	)

	result, err := updater.loadSetupData()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.GithubToken != "token123" {
		t.Errorf("Expected token 'token123', got '%s'", result.GithubToken)
	}

	if len(result.Apps) != 1 {
		t.Errorf("Expected 1 app, got %d", len(result.Apps))
	}
}

func TestGitHubDownloaderGetLatestRelease(t *testing.T) {
	release := GitHubRelease{
		TagName: "v1.0.0",
		Assets: []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		}{
			{Name: "app.tar.gz", BrowserDownloadURL: "https://example.com/app.tar.gz"},
		},
	}

	responseBody, _ := json.Marshal(release)
	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		},
	}

	downloader := &githubDownloader{client: mockClient}
	result, err := downloader.GetLatestRelease("test/repo", "token")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.TagName != "v1.0.0" {
		t.Errorf("Expected tag 'v1.0.0', got '%s'", result.TagName)
	}
}

func TestToSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user/repo", "user-repo"},
		{"User/Repo", "user-repo"},
		{"user/repo-name", "user-repo-name"},
		{"user/repo_name", "user-repo-name"},
	}

	for _, tt := range tests {
		result := toSlug(tt.input)
		if result != tt.expected {
			t.Errorf("toSlug(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestSanitizeReleaseID(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"v1.0.0", "1.0.0"},
		{"1.2.3", "1.2.3"},
		{"v2.0.0-beta", "2.0.0-beta"},
	}

	for _, tt := range tests {
		result := sanitizeReleaseID(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeReleaseID(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}
