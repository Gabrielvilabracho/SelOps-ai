package gga

import (
	"github.com/Gabrielvilabracho/selops-ai/internal/installcmd"
	"github.com/Gabrielvilabracho/selops-ai/internal/model"
	"github.com/Gabrielvilabracho/selops-ai/internal/system"
)

func InstallCommand(profile system.PlatformProfile) ([][]string, error) {
	return installcmd.NewResolver().ResolveComponentInstall(profile, model.ComponentGGA)
}

func ShouldInstall(enabled bool) bool {
	return enabled
}
