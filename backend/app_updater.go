package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type CommandExecutor interface {
	Run(command, workDir string) error
}

type ArchiveExtractor interface {
	ExtractTarGz(archivePath, destPath string) error
	ExtractZip(archivePath, destPath string) error
}

type GitHubDownloader interface {
	GetLatestRelease(repo, token string) (*GitHubRelease, error)
	DownloadAsset(url, token string) (io.ReadCloser, error)
}

type shellExecutor struct{}

func (e *shellExecutor) Run(command, workDir string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type archiveExtractorImpl struct {
	fs FileSystemOps
}

func (e *archiveExtractorImpl) ExtractTarGz(archivePath, destPath string) error {
	file, err := e.fs.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := e.fs.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := e.fs.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			outFile, err := e.fs.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
			if err := os.Chmod(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *archiveExtractorImpl) ExtractZip(archivePath, destPath string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destPath, f.Name)

		if f.FileInfo().IsDir() {
			e.fs.MkdirAll(fpath, 0755)
			continue
		}

		if err := e.fs.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		outFile, err := e.fs.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

type githubDownloader struct {
	client HTTPClient
}

func (gd *githubDownloader) GetLatestRelease(repo, token string) (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := gd.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func (gd *githubDownloader) DownloadAsset(url, token string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := gd.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	return resp.Body, nil
}

type AppUpdater struct {
	setupFilePath  string
	fs             FileSystemOps
	extractor      ArchiveExtractor
	downloader     GitHubDownloader
	ProcessManager ProcessManager
}

func NewAppUpdater(
	setupFilePath string,
	fs FileSystemOps,
	extractor ArchiveExtractor,
	downloader GitHubDownloader,
	processManager ProcessManager,
) *AppUpdater {
	return &AppUpdater{
		setupFilePath:  setupFilePath,
		fs:             fs,
		extractor:      extractor,
		downloader:     downloader,
		ProcessManager: processManager,
	}
}

func NewDefaultAppUpdater(setupFilePath string) *AppUpdater {
	fs := &osFileSystem{}
	httpClient := &http.Client{Timeout: 30 * time.Second}
	return NewAppUpdater(
		setupFilePath,
		fs,
		&archiveExtractorImpl{fs: fs},
		&githubDownloader{client: httpClient},
		NewProcessManager(),
	)
}

func (au *AppUpdater) Start() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	au.checkAndUpdateApps()

	for range ticker.C {
		au.checkAndUpdateApps()
	}
}

func (au *AppUpdater) checkAndUpdateApps() {
	setupData, err := au.loadSetupData()
	if err != nil {
		log.Printf("Failed to load setup data: %v", err)
		return
	}

	if setupData.GithubToken == "" {
		log.Println("No GitHub token configured, skipping app updates")
		return
	}

	for _, app := range setupData.Apps {
		if err := au.updateApp(app, setupData.GithubToken); err != nil {
			log.Printf("Failed to update app %s: %v", app.Key, err)
		}
	}
}

func (au *AppUpdater) loadSetupData() (*SetupData, error) {
	data, err := au.fs.ReadFile(au.setupFilePath)
	if err != nil {
		return nil, err
	}

	var setupData SetupData
	if err := json.Unmarshal(data, &setupData); err != nil {
		return nil, err
	}

	return &setupData, nil
}

func (au *AppUpdater) updateApp(app App, githubToken string) error {
	if app.Provider != "github" {
		return fmt.Errorf("unsupported provider: %s", app.Provider)
	}

	release, err := au.downloader.GetLatestRelease(app.Key, githubToken)
	if err != nil {
		return fmt.Errorf("failed to get latest release: %w", err)
	}

	slug := toSlug(app.Key)
	releaseID := sanitizeReleaseID(release.TagName)
	installPath := filepath.Join("/opt/zen/apps", fmt.Sprintf("%s-%s", slug, releaseID))

	existingProcess, _ := au.ProcessManager.GetProcess(app.Key)
	if existingProcess != nil && existingProcess.Version == releaseID {
		if au.ProcessManager.IsRunning(app.Key) {
			log.Printf("App %s version %s already running", app.Key, releaseID)
			return nil
		}
	}

	if _, err := au.fs.Stat(installPath); err != nil {
		log.Printf("Installing app %s version %s", app.Key, releaseID)

		if len(release.Assets) == 0 {
			return fmt.Errorf("no assets found in release")
		}

		asset := release.Assets[0]
		if err := au.downloadAndExtract(asset.BrowserDownloadURL, asset.Name, installPath, githubToken); err != nil {
			return fmt.Errorf("failed to download and extract: %w", err)
		}

		log.Printf("Successfully installed app %s version %s", app.Key, releaseID)
	}

	if app.Command != "" {
		if existingProcess != nil && existingProcess.Version != releaseID {
			log.Printf("Stopping old version of %s (version %s)", app.Key, existingProcess.Version)
			if err := au.ProcessManager.Stop(app.Key); err != nil {
				log.Printf("Failed to stop old version: %v", err)
			}
		}

		log.Printf("Starting app %s version %s", app.Key, releaseID)
		if err := au.ProcessManager.Start(app.Key, releaseID, app.Command, installPath); err != nil {
			return fmt.Errorf("failed to start app: %w", err)
		}
		log.Printf("App %s version %s started successfully", app.Key, releaseID)
	}

	return nil
}

func (au *AppUpdater) downloadAndExtract(url, filename, installPath, token string) error {
	if err := au.fs.MkdirAll(installPath, 0755); err != nil {
		return err
	}

	body, err := au.downloader.DownloadAsset(url, token)
	if err != nil {
		return err
	}
	defer body.Close()

	tmpFile := filepath.Join(installPath, filename)
	out, err := au.fs.Create(tmpFile)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, body); err != nil {
		return err
	}
	out.Close()

	if err := au.extractArchive(tmpFile, installPath); err != nil {
		return err
	}

	au.fs.Remove(tmpFile)
	return nil
}

func (au *AppUpdater) extractArchive(archivePath, destPath string) error {
	if strings.HasSuffix(archivePath, ".tar.gz") || strings.HasSuffix(archivePath, ".tgz") {
		return au.extractor.ExtractTarGz(archivePath, destPath)
	} else if strings.HasSuffix(archivePath, ".zip") {
		return au.extractor.ExtractZip(archivePath, destPath)
	}
	return fmt.Errorf("unsupported archive format: %s", archivePath)
}

func toSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "/", "-")
	reg := regexp.MustCompile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func sanitizeReleaseID(s string) string {
	s = strings.TrimPrefix(s, "v")
	reg := regexp.MustCompile("[^a-zA-Z0-9._-]+")
	return reg.ReplaceAllString(s, "-")
}
