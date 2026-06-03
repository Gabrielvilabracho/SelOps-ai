package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Gabrielvilabracho/selops-ai/internal/agents/kimi"
	"github.com/Gabrielvilabracho/selops-ai/internal/agents/opencode"
	"github.com/Gabrielvilabracho/selops-ai/internal/backup"
	"github.com/Gabrielvilabracho/selops-ai/internal/installcmd"
	"github.com/Gabrielvilabracho/selops-ai/internal/model"
	"github.com/Gabrielvilabracho/selops-ai/internal/system"
	"github.com/Gabrielvilabracho/selops-ai/internal/versions"
)

// missingBinaryLookPath simulates all installable binaries (engram, gga) as
// missing. Go availability is no longer required for engram installation
// (pre-built binaries are downloaded directly from GitHub Releases).
func missingBinaryLookPath(name string) (string, error) {
	return "", exec.ErrNotFound
}

func assertFileContains(t *testing.T, path string, want string) {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	if !strings.Contains(string(body), want) {
		t.Fatalf("file %q missing %q; got:\n%s", path, want, string(body))
	}
}

func stringSliceContains(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}

func engramInitCommandForTest() string {
	if _, err := exec.LookPath("pnpm"); err == nil {
		return fmt.Sprintf("pnpm dlx gentle-engram@%s pi-engram init", versions.GentleEngram)
	}
	return fmt.Sprintf("npm exec --yes --package gentle-engram@%s -- pi-engram init", versions.GentleEngram)
}

func TestRunInstallAppliesFilesystemChanges(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	runCommand = func(string, ...string) error { return nil }
	cmdLookPath = missingBinaryLookPath

	result, err := RunInstall([]string{"--agent", "opencode", "--component", "permissions"}, system.DetectionResult{})
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("verification ready = false, report = %#v", result.Verify)
	}

	configPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("expected config file %q: %v", configPath, err)
	}
}

func TestRunInstallEngramForPiAndOpenCodeProvisionsBothMCPTargets(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = func(name string) (string, error) {
		return filepath.Join(home, "bin", name), nil
	}
	restorePreflightLookPath := installcmd.OverrideLookPath(func(name string) (string, error) {
		return filepath.Join(home, "bin", name), nil
	})
	t.Cleanup(restorePreflightLookPath)

	var commands []string
	runCommand = func(name string, args ...string) error {
		commands = append(commands, strings.Join(append([]string{name}, args...), " "))
		// Simulate pi-engram init writing mcp.json with the new schema.
		isNpmEngramInit := name == "npm" && len(args) >= 7 && args[5] == "pi-engram" && args[6] == "init"
		isPnpmEngramInit := name == "pnpm" && len(args) >= 4 && args[2] == "pi-engram" && args[3] == "init"
		if isNpmEngramInit || isPnpmEngramInit {
			mcpPath := filepath.Join(home, ".pi", "agent", "mcp.json")
			if err := os.MkdirAll(filepath.Dir(mcpPath), 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(mcpPath, []byte(`{"activeMCP":"engram","mcpServers":{"engram":{"command":"node","args":["--eval","require('child_process').spawn('engram',['mcp','--tools=agent'],{stdio:'inherit'})"]}}}`+"\n"), 0o644); err != nil {
				return err
			}
		}
		return nil
	}

	result, err := RunInstall([]string{
		"--agent", "pi",
		"--agent", "opencode",
		"--component", "engram",
	}, system.DetectionResult{})
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}
	if !result.Verify.Ready {
		t.Fatalf("verification ready = false, report = %#v", result.Verify)
	}

	assertFileContains(t, filepath.Join(home, ".pi", "agent", "settings.json"), "npm:pi-mcp-adapter")
	assertFileContains(t, filepath.Join(home, ".pi", "npm", "package.json"), "pi-mcp-adapter")
	assertFileContains(t, filepath.Join(home, ".config", "opencode", "opencode.json"), "engram")

	if !stringSliceContains(commands, "pi install npm:pi-mcp-adapter") {
		t.Fatalf("commands missing %q; got %v", "pi install npm:pi-mcp-adapter", commands)
	}
	if !stringSliceContains(commands, fmt.Sprintf("npm exec --yes --package gentle-engram@%s -- pi-engram init", versions.GentleEngram)) &&
		!stringSliceContains(commands, fmt.Sprintf("pnpm dlx gentle-engram@%s pi-engram init", versions.GentleEngram)) {
		t.Fatalf("commands missing Engram init command; got %v", commands)
	}
}

func TestPiAgentInstallRunsPackageCommandsWhenPiAlreadyInstalled(t *testing.T) {
	binDir := t.TempDir()
	fakePi := filepath.Join(binDir, "pi")
	if err := os.WriteFile(fakePi, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("WriteFile(fake pi) error = %v", err)
	}
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	restorePreflightLookPath := installcmd.OverrideLookPath(func(name string) (string, error) {
		if name == "pi" {
			return fakePi, nil
		}
		return "", exec.ErrNotFound
	})
	t.Cleanup(restorePreflightLookPath)

	restoreCommand := runCommand
	t.Cleanup(func() { runCommand = restoreCommand })

	var commands []string
	runCommand = func(name string, args ...string) error {
		commands = append(commands, strings.Join(append([]string{name}, args...), " "))
		return nil
	}

	step := agentInstallStep{
		id:      "agent:pi",
		agent:   model.AgentPi,
		homeDir: t.TempDir(),
	}

	if err := step.Run(); err != nil {
		t.Fatalf("agentInstallStep.Run() error = %v", err)
	}

	for _, want := range []string{
		"pi install npm:gentle-pi",
		"pi install npm:gentle-engram",
		"pi install npm:pi-mcp-adapter",
		engramInitCommandForTest(),
		"pi install npm:pi-subagents",
		"pi install npm:pi-intercom",
		"pi install npm:@juicesharp/rpiv-ask-user-question",
		"pi install npm:pi-web-access",
		"pi install npm:pi-lens",
		"pi install npm:@juicesharp/rpiv-todo",
		"pi install npm:pi-btw",
	} {
		if !stringSliceContains(commands, want) {
			t.Fatalf("commands missing %q; got %v", want, commands)
		}
	}
}

func TestRunInstallRollsBackOnComponentFailure(t *testing.T) {
	home := t.TempDir()
	settingsPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	before := []byte("{\n  \"existing\": true\n}\n")
	if err := os.WriteFile(settingsPath, before, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})
	cmdLookPath = missingBinaryLookPath

	osUserHomeDir = func() (string, error) { return home, nil }
	runCommand = func(name string, args ...string) error {
		if name == "brew" && len(args) == 2 && args[0] == "install" && args[1] == "engram" {
			return os.ErrPermission
		}
		return nil
	}

	// Use only engram (not context7) — context7 injects MCP config into
	// the settings file and does not have a rollback step, so including it
	// makes the before/after comparison fail even when the pipeline rollback
	// works correctly. Context7 rollback is tracked separately.
	_, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "engram"},
		system.DetectionResult{},
	)
	if err == nil {
		t.Fatalf("RunInstall() expected error")
	}

	if !strings.Contains(err.Error(), "execute install pipeline") {
		t.Fatalf("RunInstall() error = %v", err)
	}

	after, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(after) != string(before) {
		t.Fatalf("settings content changed after rollback\nafter=%s\nbefore=%s", after, before)
	}
}

// --- Batch D: Linux profile runtime wiring integration tests ---

// linuxDetectionResult builds a DetectionResult with a Linux profile for integration tests.
func linuxDetectionResult(distro, pkgMgr string) system.DetectionResult {
	return system.DetectionResult{
		System: system.SystemInfo{
			OS:        "linux",
			Arch:      "amd64",
			Shell:     "/bin/bash",
			Supported: true,
			Profile: system.PlatformProfile{
				OS:             "linux",
				LinuxDistro:    distro,
				PackageManager: pkgMgr,
				Supported:      true,
			},
		},
	}
}

// commandRecorder captures all external commands invoked during a pipeline run.
type commandRecorder struct {
	mu       sync.Mutex
	commands []string
}

func (r *commandRecorder) record(name string, args ...string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.commands = append(r.commands, fmt.Sprintf("%s %s", name, strings.Join(args, " ")))
	return nil
}

func (r *commandRecorder) get() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := make([]string, len(r.commands))
	copy(cp, r.commands)
	return cp
}

func TestRunInstallLinuxUbuntuResolvesAptCommands(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = missingBinaryLookPath
	recorder := &commandRecorder{}
	runCommand = recorder.record

	detection := linuxDetectionResult(system.LinuxDistroUbuntu, "apt")
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "permissions"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("verification ready = false, report = %#v", result.Verify)
	}

	// Verify platform decision was resolved from the Linux profile.
	if result.Resolved.PlatformDecision.OS != "linux" {
		t.Fatalf("platform decision OS = %q, want linux", result.Resolved.PlatformDecision.OS)
	}
	if result.Resolved.PlatformDecision.PackageManager != "apt" {
		t.Fatalf("platform decision package manager = %q, want apt", result.Resolved.PlatformDecision.PackageManager)
	}
}

func TestRunInstallLinuxArchResolvesPacmanCommands(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = missingBinaryLookPath
	recorder := &commandRecorder{}
	runCommand = recorder.record

	detection := linuxDetectionResult(system.LinuxDistroArch, "pacman")
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "permissions"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("verification ready = false, report = %#v", result.Verify)
	}

	if result.Resolved.PlatformDecision.PackageManager != "pacman" {
		t.Fatalf("platform decision package manager = %q, want pacman", result.Resolved.PlatformDecision.PackageManager)
	}
}

func TestRunInstallLinuxUbuntuWithEngramUsesDirectDownload(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = missingBinaryLookPath
	recorder := &commandRecorder{}
	runCommand = recorder.record

	// Override engramDownloadFn to avoid real HTTP calls.
	origDownloadFn := engramDownloadFn
	engramDownloadFn = func(profile system.PlatformProfile) (string, error) {
		return "/tmp/fake-engram", nil
	}
	t.Cleanup(func() { engramDownloadFn = origDownloadFn })

	detection := linuxDetectionResult(system.LinuxDistroUbuntu, "apt")
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "engram"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("verification ready = false, report = %#v", result.Verify)
	}

	// Must NOT use go install for engram on Linux.
	for _, cmd := range recorder.get() {
		if strings.Contains(cmd, "go install") && strings.Contains(cmd, "engram") {
			t.Fatalf("Linux engram install should NOT use go install, got command: %s", cmd)
		}
	}
}

func TestRunInstallLinuxArchWithEngramUsesDirectDownload(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = missingBinaryLookPath
	recorder := &commandRecorder{}
	runCommand = recorder.record

	origDownloadFn := engramDownloadFn
	engramDownloadFn = func(profile system.PlatformProfile) (string, error) {
		return "/tmp/fake-engram", nil
	}
	t.Cleanup(func() { engramDownloadFn = origDownloadFn })

	detection := linuxDetectionResult(system.LinuxDistroArch, "pacman")
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "engram"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("verification ready = false, report = %#v", result.Verify)
	}

	// Must NOT use go install for engram on Arch Linux.
	for _, cmd := range recorder.get() {
		if strings.Contains(cmd, "go install") && strings.Contains(cmd, "engram") {
			t.Fatalf("Arch Linux engram install should NOT use go install, got command: %s", cmd)
		}
	}
}

func TestRunInstallLinuxRollsBackOnComponentFailure(t *testing.T) {
	home := t.TempDir()
	settingsPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	before := []byte("{\n  \"linux-original\": true\n}\n")
	if err := os.WriteFile(settingsPath, before, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})
	cmdLookPath = missingBinaryLookPath

	osUserHomeDir = func() (string, error) { return home, nil }
	runCommand = func(name string, args ...string) error { return nil }

	// Fail the engram download to trigger rollback.
	origDownloadFn := engramDownloadFn
	engramDownloadFn = func(profile system.PlatformProfile) (string, error) {
		return "", os.ErrPermission
	}
	t.Cleanup(func() { engramDownloadFn = origDownloadFn })

	detection := linuxDetectionResult(system.LinuxDistroUbuntu, "apt")
	// Exclude context7 — it has no rollback and taints the settings file.
	_, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "engram"},
		detection,
	)
	if err == nil {
		t.Fatalf("RunInstall() expected error")
	}

	if !strings.Contains(err.Error(), "execute install pipeline") {
		t.Fatalf("RunInstall() error = %v", err)
	}

	// Verify rollback restored the original file.
	after, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(after) != string(before) {
		t.Fatalf("settings content changed after rollback on Linux\nafter=%s\nbefore=%s", after, before)
	}
}

func TestRunInstallFedoraQwenEngramSkipsUnsupportedSetupAndWritesSettings(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = missingBinaryLookPath
	recorder := &commandRecorder{}
	runCommand = recorder.record

	origDownloadFn := engramDownloadFn
	engramDownloadFn = func(profile system.PlatformProfile) (string, error) {
		return filepath.Join(home, "bin", "engram"), nil
	}
	t.Cleanup(func() { engramDownloadFn = origDownloadFn })

	detection := linuxDetectionResult(system.LinuxDistroFedora, "dnf")
	result, err := RunInstall(
		[]string{"--agent", "qwen-code", "--component", "engram"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}
	if !result.Verify.Ready {
		t.Fatalf("verification ready = false, report = %#v", result.Verify)
	}

	settingsPath := filepath.Join(home, ".qwen", "settings.json")
	if _, err := os.Stat(settingsPath); err != nil {
		t.Fatalf("expected qwen settings at %q: %v", settingsPath, err)
	}

	for _, cmd := range recorder.get() {
		if strings.Contains(cmd, "engram setup qwen-code") {
			t.Fatalf("unexpected unsupported setup command: %s", cmd)
		}
	}
}

func TestRunInstallLinuxAgentInstallResolvesGoInstallCommand(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = missingBinaryLookPath
	recorder := &commandRecorder{}
	runCommand = recorder.record

	// Set the agent adapter's lookPath to simulate missing opencode
	opencodeAdapterLookPath := opencode.LookPathOverride
	opencode.LookPathOverride = missingBinaryLookPath
	t.Cleanup(func() {
		opencode.LookPathOverride = opencodeAdapterLookPath
	})

	detection := linuxDetectionResult(system.LinuxDistroUbuntu, "apt")
	_, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "permissions"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	// OpenCode on Ubuntu should resolve via npm install (official method from opencode.ai).
	commands := recorder.get()
	foundNpmInstall := false
	for _, cmd := range commands {
		if strings.Contains(cmd, "sudo npm install -g --ignore-scripts opencode-ai@"+versions.OpenCode) {
			foundNpmInstall = true
			break
		}
	}
	if !foundNpmInstall {
		t.Fatalf("expected npm install command for opencode agent, got commands: %v", commands)
	}
}

// --- Batch E: Linux verification and macOS parity matrix ---

func TestRunInstallLinuxVerificationReportsReadyOnSuccess(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	runCommand = func(string, ...string) error { return nil }
	cmdLookPath = missingBinaryLookPath

	detection := linuxDetectionResult(system.LinuxDistroUbuntu, "apt")
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "permissions"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("Verify.Ready = false, want true for successful Linux install")
	}
	if result.Verify.Failed != 0 {
		t.Fatalf("Verify.Failed = %d, want 0", result.Verify.Failed)
	}
}

func TestRunInstallLinuxArchVerificationReportsReadyOnSuccess(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	runCommand = func(string, ...string) error { return nil }
	cmdLookPath = missingBinaryLookPath

	detection := linuxDetectionResult(system.LinuxDistroArch, "pacman")
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "permissions"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("Verify.Ready = false, want true for successful Arch install")
	}
}

func TestRunInstallLinuxDryRunSkipsVerification(t *testing.T) {
	detection := linuxDetectionResult(system.LinuxDistroUbuntu, "apt")
	result, err := RunInstall([]string{"--dry-run"}, detection)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.DryRun {
		t.Fatalf("DryRun = false, want true")
	}
	// Verify report should be zero-value (no checks run in dry-run)
	if result.Verify.Passed != 0 || result.Verify.Failed != 0 {
		t.Fatalf("expected zero verify counters in dry-run, got passed=%d failed=%d", result.Verify.Passed, result.Verify.Failed)
	}
}

func TestRunInstallLinuxDryRunPlatformDecisionRendersCorrectly(t *testing.T) {
	detection := linuxDetectionResult(system.LinuxDistroArch, "pacman")
	result, err := RunInstall([]string{"--dry-run"}, detection)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	output := RenderDryRun(result)
	want := "os=linux distro=arch package-manager=pacman status=supported"
	if !strings.Contains(output, want) {
		t.Fatalf("RenderDryRun() missing platform decision\noutput=%s\nwant contains=%s", output, want)
	}
}

// --- macOS parity regression checks ---

func macOSDetectionResult() system.DetectionResult {
	return system.DetectionResult{
		System: system.SystemInfo{
			OS:        "darwin",
			Arch:      "arm64",
			Shell:     "/bin/zsh",
			Supported: true,
			Profile: system.PlatformProfile{
				OS:             "darwin",
				PackageManager: "brew",
				Supported:      true,
			},
		},
	}
}

func TestRunInstallMacOSStillResolvesBrewCommands(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = missingBinaryLookPath
	recorder := &commandRecorder{}
	runCommand = recorder.record

	detection := macOSDetectionResult()
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "engram"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("macOS verification ready = false")
	}

	// Verify brew install command was used, not apt or pacman.
	commands := recorder.get()
	foundBrew := false
	for _, cmd := range commands {
		if strings.Contains(cmd, "brew install engram") {
			foundBrew = true
			break
		}
	}
	if !foundBrew {
		t.Fatalf("expected brew install for macOS engram, got commands: %v", commands)
	}
}

func TestRunInstallMacOSDryRunPlatformDecision(t *testing.T) {
	detection := macOSDetectionResult()
	result, err := RunInstall([]string{"--dry-run"}, detection)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if result.Resolved.PlatformDecision.OS != "darwin" {
		t.Fatalf("macOS platform decision OS = %q, want darwin", result.Resolved.PlatformDecision.OS)
	}
	if result.Resolved.PlatformDecision.PackageManager != "brew" {
		t.Fatalf("macOS platform decision PM = %q, want brew", result.Resolved.PlatformDecision.PackageManager)
	}
	if !result.Resolved.PlatformDecision.Supported {
		t.Fatalf("macOS platform decision Supported = false, want true")
	}
}

func TestRunInstallMacOSVerificationMatchesPreLinuxBehavior(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	runCommand = func(string, ...string) error { return nil }
	cmdLookPath = missingBinaryLookPath

	detection := macOSDetectionResult()
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "permissions"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("macOS verify ready = false, want true")
	}
	if result.Verify.Failed != 0 {
		t.Fatalf("macOS verify failed = %d, want 0", result.Verify.Failed)
	}
}

func TestRunInstallMacOSRollbackStillWorks(t *testing.T) {
	home := t.TempDir()
	settingsPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	before := []byte("{\n  \"macos-original\": true\n}\n")
	if err := os.WriteFile(settingsPath, before, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})
	cmdLookPath = missingBinaryLookPath

	osUserHomeDir = func() (string, error) { return home, nil }
	runCommand = func(name string, args ...string) error {
		if name == "brew" && len(args) == 2 && args[0] == "install" && args[1] == "engram" {
			return os.ErrPermission
		}
		return nil
	}

	detection := macOSDetectionResult()
	// Exclude context7 — it has no rollback and taints the settings file.
	_, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "engram"},
		detection,
	)
	if err == nil {
		t.Fatalf("RunInstall() expected error")
	}

	if !strings.Contains(err.Error(), "execute install pipeline") {
		t.Fatalf("RunInstall() error = %v", err)
	}

	after, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(after) != string(before) {
		t.Fatalf("macOS settings changed after rollback\nafter=%s\nbefore=%s", after, before)
	}
}

// --- Skip-when-installed and Go auto-install tests ---

func TestRunInstallEngramSkipsInstallWhenAlreadyOnPath(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	// Simulate engram already installed on PATH.
	cmdLookPath = func(name string) (string, error) {
		return "/usr/local/bin/" + name, nil
	}
	recorder := &commandRecorder{}
	runCommand = recorder.record

	detection := macOSDetectionResult()
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "engram"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("verification ready = false")
	}

	// No brew/go install commands should have been recorded — only agent install.
	for _, cmd := range recorder.get() {
		if strings.Contains(cmd, "brew install engram") || (strings.Contains(cmd, "go install") && strings.Contains(cmd, "engram")) {
			t.Fatalf("expected engram install to be skipped, but got command: %s", cmd)
		}
	}
}

func TestRunInstallEngramAttemptsOpenCodeSetupWhenBinaryPresent(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = func(name string) (string, error) {
		return "/usr/local/bin/" + name, nil
	}
	recorder := &commandRecorder{}
	runCommand = recorder.record

	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "engram"},
		macOSDetectionResult(),
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}
	if !result.Verify.Ready {
		t.Fatalf("verification ready = false")
	}

	commands := recorder.get()
	foundSetup := false
	for _, cmd := range commands {
		if strings.Contains(cmd, "engram setup opencode") {
			foundSetup = true
			break
		}
	}
	if !foundSetup {
		t.Fatalf("expected engram setup command, got commands: %v", commands)
	}
}

func TestRunInstallEngramFallsBackToInjectWhenSetupFails(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = func(name string) (string, error) {
		return "/usr/local/bin/" + name, nil
	}
	runCommand = func(name string, args ...string) error {
		if name == "engram" && len(args) == 2 && args[0] == "setup" && args[1] == "opencode" {
			return errors.New("setup failed")
		}
		return nil
	}

	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "engram"},
		macOSDetectionResult(),
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}
	if !result.Verify.Ready {
		t.Fatalf("verification ready = false")
	}

	configPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("expected fallback inject to create %q: %v", configPath, err)
	}
}

func TestRunInstallEngramSetupStrictFailsWhenSetupFails(t *testing.T) {
	t.Setenv("GENTLE_AI_ENGRAM_SETUP_STRICT", "1")

	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	origUserHomeDirFn := backup.UserHomeDirFn
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
		backup.UserHomeDirFn = origUserHomeDirFn
	})
	// Override restore path validation to accept test temp dirs.
	backup.UserHomeDirFn = func() (string, error) { return home, nil }

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = func(name string) (string, error) {
		return "/usr/local/bin/" + name, nil
	}
	runCommand = func(name string, args ...string) error {
		if name == "engram" && len(args) == 2 && args[0] == "setup" && args[1] == "opencode" {
			return errors.New("setup failed")
		}
		return nil
	}

	_, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "engram"},
		macOSDetectionResult(),
	)
	if err == nil {
		t.Fatalf("RunInstall() expected error in strict setup mode")
	}
	if !strings.Contains(err.Error(), "engram setup for \"opencode\"") {
		t.Fatalf("RunInstall() error = %v, want setup error", err)
	}
}

func TestRunInstallEngramDefaultModeAttemptsClaudeSetup(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = func(name string) (string, error) {
		return "/usr/local/bin/" + name, nil
	}
	recorder := &commandRecorder{}
	runCommand = recorder.record

	result, err := RunInstall(
		[]string{"--agent", "claude-code", "--component", "engram"},
		macOSDetectionResult(),
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}
	if !result.Verify.Ready {
		t.Fatalf("verification ready = false")
	}

	commands := recorder.get()
	foundSetup := false
	for _, cmd := range commands {
		if strings.Contains(cmd, "engram setup claude-code") {
			foundSetup = true
			break
		}
	}
	if !foundSetup {
		t.Fatalf("expected default setup mode to attempt claude-code setup, got commands: %v", commands)
	}
}

func TestRunInstallAntigravityInitializesCLISettingsAfterEngramSetup(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = func(name string) (string, error) {
		return "/usr/local/bin/" + name, nil
	}
	runCommand = func(name string, args ...string) error {
		if name == "engram" && len(args) == 2 && args[0] == "setup" && args[1] == "gemini-cli" {
			settingsPath := filepath.Join(home, ".gemini", "settings.json")
			if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
				return err
			}
			return os.WriteFile(settingsPath, []byte("{\"theme\":\"dark\"}\n"), 0o644)
		}
		return nil
	}

	result, err := RunInstall(
		[]string{"--agent", "antigravity", "--component", "engram", "--component", "context7", "--component", "permissions"},
		macOSDetectionResult(),
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}
	if !result.Verify.Ready {
		t.Fatalf("verification ready = false")
	}

	settingsPath := filepath.Join(home, ".gemini", "antigravity-cli", "settings.json")
	got, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", settingsPath, err)
	}
	if string(got) != "{}\n" {
		t.Fatalf("antigravity settings = %q, want initialized empty settings", got)
	}
}

func TestRunInstallDeduplicatesSharedEngramSetupSlugs(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = func(name string) (string, error) {
		return "/usr/local/bin/" + name, nil
	}

	recorder := &commandRecorder{}
	runCommand = func(name string, args ...string) error {
		if err := recorder.record(name, args...); err != nil {
			return err
		}
		if name == "engram" && len(args) == 2 && args[0] == "setup" && args[1] == "gemini-cli" {
			settingsPath := filepath.Join(home, ".gemini", "settings.json")
			if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
				return err
			}
			return os.WriteFile(settingsPath, []byte("{\"theme\":\"dark\"}\n"), 0o644)
		}
		return nil
	}

	result, err := RunInstall(
		[]string{"--agent", "gemini-cli", "--agent", "antigravity", "--component", "engram", "--component", "context7", "--component", "permissions"},
		macOSDetectionResult(),
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}
	if !result.Verify.Ready {
		t.Fatalf("verification ready = false")
	}

	var setupCount int
	for _, cmd := range recorder.get() {
		if strings.Contains(cmd, "engram setup gemini-cli") {
			setupCount++
		}
	}
	if setupCount != 1 {
		t.Fatalf("engram setup gemini-cli count = %d, want 1", setupCount)
	}
}

func TestRunInstallGGASkipsInstallWhenAlreadyOnPath(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = func(name string) (string, error) {
		return "/usr/local/bin/" + name, nil
	}
	recorder := &commandRecorder{}
	runCommand = recorder.record

	detection := macOSDetectionResult()
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "gga"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("verification ready = false")
	}

	// No brew/git clone commands for GGA should have been recorded.
	for _, cmd := range recorder.get() {
		if strings.Contains(cmd, "gga") || strings.Contains(cmd, "gentleman-guardian-angel") {
			t.Fatalf("expected gga install to be skipped, but got command: %s", cmd)
		}
	}

	prModePath := filepath.Join(home, ".local", "share", "gga", "lib", "pr_mode.sh")
	content, err := os.ReadFile(prModePath)
	if err != nil {
		t.Fatalf("expected gga runtime asset at %q: %v", prModePath, err)
	}
	if !strings.Contains(string(content), "detect_base_branch") {
		t.Fatalf("expected pr_mode.sh to contain detect_base_branch")
	}
}

func TestRunInstallGGALinuxIncludesTempCleanupBeforeClone(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = func(name string) (string, error) {
		if name == "gga" {
			return "", exec.ErrNotFound
		}
		return "/usr/local/bin/" + name, nil
	}
	recorder := &commandRecorder{}
	runCommand = recorder.record

	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "gga"},
		linuxDetectionResult(system.LinuxDistroUbuntu, "apt"),
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}
	if !result.Verify.Ready {
		t.Fatalf("verification ready = false")
	}

	commands := recorder.get()
	cleanupIdx := -1
	cloneIdx := -1
	for i, cmd := range commands {
		if strings.Contains(cmd, "rm -rf /tmp/gentleman-guardian-angel") {
			cleanupIdx = i
		}
		if strings.Contains(cmd, "git clone https://github.com/Gentleman-Programming/gentleman-guardian-angel.git /tmp/gentleman-guardian-angel") {
			cloneIdx = i
		}
	}

	for _, cmd := range commands {
		if strings.Contains(cmd, "gga install") || strings.Contains(cmd, "gga init") {
			t.Fatalf("expected global gga provisioning only, got repo-level command: %s", cmd)
		}
	}

	if cleanupIdx == -1 {
		t.Fatalf("expected cleanup command before clone, got commands: %v", commands)
	}
	if cloneIdx == -1 {
		t.Fatalf("expected clone command, got commands: %v", commands)
	}
	if cleanupIdx >= cloneIdx {
		t.Fatalf("cleanup should run before clone (cleanup=%d clone=%d)", cleanupIdx, cloneIdx)
	}
}

// TestRunInstallEngramLinuxUsesDirectDownloadNoGoRequired verifies that on Linux,
// engram is now installed via pre-built binary download — Go is NOT required.
func TestRunInstallEngramLinuxUsesDirectDownloadNoGoRequired(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	// Simulate: engram missing, Go also NOT available — should still succeed.
	cmdLookPath = func(string) (string, error) {
		return "", exec.ErrNotFound
	}
	recorder := &commandRecorder{}
	runCommand = recorder.record

	// Override download to succeed without hitting GitHub.
	origDownloadFn := engramDownloadFn
	engramDownloadFn = func(profile system.PlatformProfile) (string, error) {
		return "/tmp/fake-engram", nil
	}
	t.Cleanup(func() { engramDownloadFn = origDownloadFn })

	detection := linuxDetectionResult(system.LinuxDistroUbuntu, "apt")
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "engram"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("verification ready = false")
	}

	// Neither "go install" nor "apt-get install golang" should appear.
	for _, cmd := range recorder.get() {
		if strings.Contains(cmd, "apt-get install -y golang") {
			t.Fatalf("Go should NOT be auto-installed (no longer needed for engram), got command: %s", cmd)
		}
		if strings.Contains(cmd, "go install") && strings.Contains(cmd, "engram") {
			t.Fatalf("engram should NOT be installed via go install, got command: %s", cmd)
		}
	}
}

// TestRunInstallEngramLinuxNeverInstallsGo verifies that even if Go is present,
// we never install Go as a prerequisite for engram (direct download path).
func TestRunInstallEngramLinuxNeverInstallsGo(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = missingBinaryLookPath
	recorder := &commandRecorder{}
	runCommand = recorder.record

	origDownloadFn := engramDownloadFn
	engramDownloadFn = func(profile system.PlatformProfile) (string, error) {
		return "/tmp/fake-engram", nil
	}
	t.Cleanup(func() { engramDownloadFn = origDownloadFn })

	detection := linuxDetectionResult(system.LinuxDistroUbuntu, "apt")
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "engram"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("verification ready = false")
	}

	// No Go installation commands should appear.
	for _, cmd := range recorder.get() {
		if strings.Contains(cmd, "apt-get install -y golang") || strings.Contains(cmd, "apt-get install -y go") {
			t.Fatalf("Go should never be installed as engram dependency, got command: %s", cmd)
		}
	}
}

func TestRunInstallEngramBrewSkipsGoCheck(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	// Simulate: engram missing — brew platform, no Go or download needed.
	cmdLookPath = func(string) (string, error) {
		return "", exec.ErrNotFound
	}
	recorder := &commandRecorder{}
	runCommand = recorder.record

	detection := macOSDetectionResult()
	result, err := RunInstall(
		[]string{"--agent", "opencode", "--component", "engram"},
		detection,
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("verification ready = false")
	}

	// Should use brew install, NOT go install, and no Go auto-install.
	commands := recorder.get()
	for _, cmd := range commands {
		if strings.Contains(cmd, "golang") || strings.Contains(cmd, "apt-get") {
			t.Fatalf("brew platform should not install Go, got command: %s", cmd)
		}
		if strings.Contains(cmd, "go install") {
			t.Fatalf("brew platform should not use go install, got command: %s", cmd)
		}
	}

	foundBrew := false
	for _, cmd := range commands {
		if strings.Contains(cmd, "brew install engram") {
			foundBrew = true
		}
	}
	if !foundBrew {
		t.Fatalf("expected brew install engram, got commands: %v", commands)
	}
}

// TestRunInstallDryRunMatchesActualInstall verifies parity: every file path
// reported by the dry-run plan is actually created by the real install.
//
// Strategy:
//  1. Run with DryRun=true to obtain the resolved plan (agents + ordered components).
//  2. Derive the expected file paths from the plan using componentPaths() — the
//     same function the runtime uses for backup targets and post-apply verification.
//  3. Run the real install (same flags, same mocks, fresh temp dir).
//  4. Assert that every expected file exists on disk — no missing files.
func TestRunInstallDryRunMatchesActualInstall(t *testing.T) {
	// ── Phase 1: dry-run — resolve the plan ───────────────────────────────────
	// We do NOT need temp dir or mocks for dry-run; it never touches the FS.
	installArgs := []string{"--agent", "opencode", "--component", "permissions"}
	dryRunArgs := append([]string{"--dry-run"}, installArgs...)
	dryResult, err := RunInstall(dryRunArgs, system.DetectionResult{})
	if err != nil {
		t.Fatalf("dry-run RunInstall() error = %v", err)
	}
	if !dryResult.DryRun {
		t.Fatalf("expected DryRun=true in result, got false")
	}

	// Use a synthetic home dir for path computation — the paths are derived
	// from the resolved plan (agents + components) and will use this root.
	// We reuse the same dir for the real install so the paths are identical.
	home := t.TempDir()

	// Derive expected file paths from the dry-run plan.  componentPaths() is
	// the single source of truth that both backup and verification use.
	adapters := resolveAdapters(dryResult.Resolved.Agents)
	var expectedPaths []string
	for _, component := range dryResult.Resolved.OrderedComponents {
		expectedPaths = append(expectedPaths, componentPaths(home, dryResult.Selection, adapters, component)...)
	}
	if len(expectedPaths) == 0 {
		t.Fatal("dry-run resolved zero file paths — test is misconfigured")
	}

	// ── Phase 2: real install — apply the plan ────────────────────────────────
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	runCommand = func(string, ...string) error { return nil }
	cmdLookPath = missingBinaryLookPath

	realResult, err := RunInstall(installArgs, system.DetectionResult{})
	if err != nil {
		t.Fatalf("real RunInstall() error = %v", err)
	}
	if !realResult.Verify.Ready {
		t.Fatalf("post-apply verification not ready: %#v", realResult.Verify)
	}

	// ── Phase 3: parity assertion ─────────────────────────────────────────────
	// Every file the dry-run said would be touched must exist on disk.
	var missing []string
	for _, path := range expectedPaths {
		if _, statErr := os.Stat(path); statErr != nil {
			missing = append(missing, path)
		}
	}
	if len(missing) > 0 {
		t.Errorf("dry-run planned %d file(s) that were NOT created by the real install:", len(missing))
		for _, p := range missing {
			t.Errorf("  missing: %s", p)
		}
	}
}


func TestEnsureGoAvailableAfterInstallWindowsRefreshesPath(t *testing.T) {
	restoreLookPath := cmdLookPath
	restoreStat := osStat
	restoreSetenv := osSetenv
	oldPath := os.Getenv("PATH")
	oldProgramFiles := os.Getenv("ProgramFiles")
	t.Cleanup(func() {
		cmdLookPath = restoreLookPath
		osStat = restoreStat
		osSetenv = restoreSetenv
		_ = os.Setenv("PATH", oldPath)
		_ = os.Setenv("ProgramFiles", oldProgramFiles)
	})

	programFiles := `C:\Program Files`
	if err := os.Setenv("ProgramFiles", programFiles); err != nil {
		t.Fatalf("Setenv(ProgramFiles) error = %v", err)
	}
	if err := os.Setenv("PATH", `C:\Windows\System32`); err != nil {
		t.Fatalf("Setenv(PATH) error = %v", err)
	}

	cmdLookPath = func(name string) (string, error) {
		if name == "go" {
			return "", exec.ErrNotFound
		}
		return name, nil
	}
	osStat = func(name string) (os.FileInfo, error) {
		want := filepath.Join(programFiles, "Go", "bin", "go.exe")
		if name == want {
			return fakeFileInfo{name: "go.exe"}, nil
		}
		return nil, os.ErrNotExist
	}
	osSetenv = os.Setenv

	if err := ensureGoAvailableAfterInstall(system.PlatformProfile{OS: "windows", PackageManager: "winget"}); err != nil {
		t.Fatalf("ensureGoAvailableAfterInstall() error = %v", err)
	}

	updatedPath := os.Getenv("PATH")
	expectedPrefix := filepath.Join(programFiles, "Go", "bin") + string(os.PathListSeparator)
	if !strings.HasPrefix(updatedPath, expectedPrefix) {
		t.Fatalf("PATH = %q, want prefix %q", updatedPath, expectedPrefix)
	}
}

type fakeFileInfo struct{ name string }

func (f fakeFileInfo) Name() string     { return f.name }
func (fakeFileInfo) Size() int64        { return 0 }
func (fakeFileInfo) Mode() os.FileMode  { return 0 }
func (fakeFileInfo) ModTime() time.Time { return time.Time{} }
func (fakeFileInfo) IsDir() bool        { return false }
func (fakeFileInfo) Sys() any           { return nil }

// TestRunInstallUpgradeIdempotency verifies that running install twice with the
// same configuration does NOT duplicate any content.  The second run must be a
// no-op or a clean update — never an append of already-present sections or MCP
// entries.

// --- Custom preset integration tests ---

func TestRunInstallCustomPresetNoComponentsIsNoop(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	runCommand = func(string, ...string) error { return nil }
	cmdLookPath = missingBinaryLookPath

	result, err := RunInstall(
		[]string{"--agent", "claude-code", "--preset", "custom"},
		system.DetectionResult{},
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	// Custom preset with no components should resolve to zero ordered components.
	if len(result.Resolved.OrderedComponents) != 0 {
		t.Fatalf("expected 0 ordered components for custom preset, got %d: %v",
			len(result.Resolved.OrderedComponents), result.Resolved.OrderedComponents)
	}
}



func TestRunInstallCustomPresetDryRunShowsCustomPreset(t *testing.T) {
	result, err := RunInstall(
		[]string{"--agent", "claude-code", "--preset", "custom", "--dry-run"},
		system.DetectionResult{},
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.DryRun {
		t.Fatalf("expected DryRun=true")
	}

	if result.Selection.Preset != model.PresetCustom {
		t.Fatalf("preset = %q, want %q", result.Selection.Preset, model.PresetCustom)
	}

	// Zero components when no --component flags provided.
	if len(result.Resolved.OrderedComponents) != 0 {
		t.Fatalf("expected 0 ordered components, got %d", len(result.Resolved.OrderedComponents))
	}

	output := RenderDryRun(result)
	if !strings.Contains(output, "custom") {
		t.Fatalf("dry-run output missing 'custom' preset name:\n%s", output)
	}
}


// TestOpenCodePersonaBeforeSDDPreservesAllSections is the regression test for
// issue #121: on StrategyFileReplace agents, if Persona ran after SDD it would
// overwrite the entire AGENTS.md, destroying the SDD orchestrator section.
//
// This test exercises the full install pipeline for OpenCode with Persona +
// Engram + SDD selected together and verifies that the final AGENTS.md
// contains all three sections with no duplicates.
func TestRunInstallKimiBootstrapsHub(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})
	osUserHomeDir = func() (string, error) { return home, nil }
	runCommand = func(string, ...string) error { return nil }
	cmdLookPath = missingBinaryLookPath
	restoreInstallcmdLookPath := installcmd.OverrideLookPath(func(name string) (string, error) {
		if name == "uv" {
			return "/usr/bin/uv", nil
		}
		return "", exec.ErrNotFound
	})
	t.Cleanup(restoreInstallcmdLookPath)

	// Install Kimi with minimalist component (e.g., permissions only, NO persona).
	_, err := RunInstall(
		[]string{"--agent", "kimi", "--component", "permissions"},
		system.DetectionResult{},
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	// Verify that KIMI.md was created in the agent's config dir.
	hubPath := filepath.Join(home, ".kimi", "KIMI.md")
	if _, err := os.Stat(hubPath); err != nil {
		t.Fatalf("expected Kimi prompt hub file %q to be bootstrapped: %v", hubPath, err)
	}

	// Verify content includes sub-modules (basic check).
	content, err := os.ReadFile(hubPath)
	if err != nil {
		t.Fatalf("failed to read bootstrapped hub: %v", err)
	}
	if !strings.Contains(string(content), "{% include \"persona.md\" ignore missing %}") {
		t.Errorf("bootstrapped hub missing modular include: %s", string(content))
	}
}

func TestRunInstallKimiMissingUVFailsBeforeExecutingInstallCommands(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = missingBinaryLookPath

	recorder := &commandRecorder{}
	runCommand = recorder.record

	restoreInstallcmdLookPath := installcmd.OverrideLookPath(func(name string) (string, error) {
		if name == "uv" {
			return "", exec.ErrNotFound
		}
		return "/usr/bin/" + name, nil
	})
	t.Cleanup(restoreInstallcmdLookPath)

	_, err := RunInstall(
		[]string{"--agent", "kimi", "--component", "permissions"},
		macOSDetectionResult(),
	)
	if err == nil {
		t.Fatal("RunInstall() expected error when Kimi uv preflight fails")
	}

	if !strings.Contains(err.Error(), "preflight for agent \"kimi\"") || !strings.Contains(err.Error(), "uv") {
		t.Fatalf("RunInstall() error = %q, expected Kimi uv preflight error", err.Error())
	}

	if got := recorder.get(); len(got) != 0 {
		t.Fatalf("expected no install commands to execute before Kimi preflight failure, got: %v", got)
	}
}

func TestRunInstallKimiAlreadyInstalledDoesNotRequireUV(t *testing.T) {
	home := t.TempDir()
	restoreHome := osUserHomeDir
	restoreCommand := runCommand
	restoreLookPath := cmdLookPath
	t.Cleanup(func() {
		osUserHomeDir = restoreHome
		runCommand = restoreCommand
		cmdLookPath = restoreLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = missingBinaryLookPath
	recorder := &commandRecorder{}
	runCommand = recorder.record

	originalKimiLookPath := kimi.LookPathOverride
	kimi.LookPathOverride = func(name string) (string, error) {
		if name == "kimi" {
			return "/usr/local/bin/kimi", nil
		}
		return "", exec.ErrNotFound
	}
	t.Cleanup(func() { kimi.LookPathOverride = originalKimiLookPath })

	restoreInstallcmdLookPath := installcmd.OverrideLookPath(func(name string) (string, error) {
		if name == "uv" {
			return "", exec.ErrNotFound
		}
		return "/usr/bin/" + name, nil
	})
	t.Cleanup(restoreInstallcmdLookPath)

	result, err := RunInstall(
		[]string{"--agent", "kimi", "--component", "permissions"},
		macOSDetectionResult(),
	)
	if err != nil {
		t.Fatalf("RunInstall() error = %v", err)
	}

	if !result.Verify.Ready {
		t.Fatalf("verification ready = false, report = %#v", result.Verify)
	}

	hubPath := filepath.Join(home, ".kimi", "KIMI.md")
	if _, err := os.Stat(hubPath); err != nil {
		t.Fatalf("expected Kimi prompt hub file %q to be bootstrapped: %v", hubPath, err)
	}

	if got := recorder.get(); len(got) != 0 {
		t.Fatalf("expected no install commands when Kimi is already installed, got: %v", got)
	}
}
