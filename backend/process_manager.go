package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type ProcessInfo struct {
	PID         int    `json:"pid"`
	AppKey      string `json:"appKey"`
	Version     string `json:"version"`
	InstallPath string `json:"installPath"`
}

type ProcessManager interface {
	Start(appKey, version, command, workDir string) error
	Stop(appKey string) error
	StopAll()
	IsRunning(appKey string) bool
	GetProcess(appKey string) (*ProcessInfo, error)
}

type processManager struct {
	processes map[string]*ProcessInfo
}

func NewProcessManager() ProcessManager {
	return &processManager{
		processes: make(map[string]*ProcessInfo),
	}
}

func (pm *processManager) isProcessAlive(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func (pm *processManager) Start(appKey, version, command, workDir string) error {
	if existing, exists := pm.processes[appKey]; exists {
		if pm.isProcessAlive(existing.PID) {
			if err := pm.Stop(appKey); err != nil {
				return fmt.Errorf("failed to stop existing process: %w", err)
			}
		}
	}

	logFile := filepath.Join(workDir, "log.txt")
	logF, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logF.Close()

	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = workDir
	cmd.Stdout = logF
	cmd.Stderr = logF

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	pm.processes[appKey] = &ProcessInfo{
		PID:         cmd.Process.Pid,
		AppKey:      appKey,
		Version:     version,
		InstallPath: workDir,
	}

	return nil
}

func (pm *processManager) Stop(appKey string) error {
	info, exists := pm.processes[appKey]
	if !exists {
		return nil
	}

	process, err := os.FindProcess(info.PID)
	if err != nil {
		delete(pm.processes, appKey)
		return nil
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		process.Kill()
	}

	delete(pm.processes, appKey)
	return nil
}

func (pm *processManager) StopAll() {
	for appKey := range pm.processes {
		pm.Stop(appKey)
	}
}

func (pm *processManager) IsRunning(appKey string) bool {
	info, exists := pm.processes[appKey]
	if !exists {
		return false
	}
	return pm.isProcessAlive(info.PID)
}

func (pm *processManager) GetProcess(appKey string) (*ProcessInfo, error) {
	info, exists := pm.processes[appKey]
	if !exists {
		return nil, fmt.Errorf("process not found")
	}
	return info, nil
}
