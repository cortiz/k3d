package containerd

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/platforms"
	l "github.com/rancher/k3d/v5/pkg/logger"
	runtimeTypes "github.com/rancher/k3d/v5/pkg/runtimes/types"
	"strings"
)

// Due to some containerd client limitations we can't get too much information
// we are assuming that, k3d is running in the same arch/os that the containerd
// client
// Todo check how to extract more information

func (c Containerd) Info() (*runtimeTypes.RuntimeInfo, error) {
	client, err := getContainerdClient()
	ctx := context.Background()
	if err != nil {
		return nil, fmt.Errorf("containerd failed to provide info output: %w", err)
	}
	defer client.Close()
	containerdVersion, err := client.Version(ctx)
	if err != nil {
		return nil, fmt.Errorf("containerd failed to provide info output: %w", err)
	}
	platform := platforms.DefaultSpec()
	runtimeInfo := runtimeTypes.RuntimeInfo{
		Name:          c.ID(),
		Endpoint:      c.GetRuntimePath(),
		Version:       containerdVersion.Version,
		OS:            platform.OS,
		OSType:        platform.Variant,
		Arch:          platform.Architecture,
		CgroupDriver:  "UNKNOWN",
		Filesystem:    getAllowedFS(client),
		CgroupVersion: "UNKNOWN",
	}
	return &runtimeInfo, nil
}

func getAllowedFS(client *containerd.Client) string {
	result := "UNKNOWN"
	plugins, err := client.IntrospectionService().Plugins(context.Background(), nil)
	if err != nil {
		l.Log().Infof("Error getting Filesystem '%s'", err)
	}
	var okPlugins []string
	// Inspire from https://github.com/containerd/nerdctl/blob/master/pkg/infoutil/infoutil.go@7cfc87bedefeba60f4b574a2ca2bf73fb1ff6d5d
	for _, plugin := range plugins.Plugins {
		if plugin.InitErr == nil && strings.HasPrefix(plugin.Type, "io.containerd.snapshotter.") {
			okPlugins = append(okPlugins, plugin.ID)
		}
	}
	result = strings.Join(okPlugins, ",")
	return result
}

func (d Containerd) GetRuntimePath() string {
	return "/run/containerd/containerd.sock"
}
func (d Containerd) ID() string {
	return "containerd"
}
