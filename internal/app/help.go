package app

import (
	"fmt"
	"io"
)

func printHelp(w io.Writer, version string) {
	fmt.Fprintf(w, `selops — SelOps: AI operations CLI (%s)

USAGE
  selops                     Launch interactive TUI
  selops <command> [flags]

COMMANDS
  install      Configure AI coding agents on this machine
  uninstall    Remove SelOps managed files from this machine
  sync         Sync agent configs and skills to current version
  skill-registry refresh
               Refresh .atl/skill-registry.md with cache-hit fast path
  update       Check for available updates
  upgrade      Apply updates to managed tools
  restore      Restore a config backup
  doctor       Run ecosystem health diagnostics
  version      Print version

FLAGS
  --help, -h    Show this help

Run 'selops help' for this message.
Documentation: https://github.com/Gabrielvilabracho/selops-ai
`, version)
}
