package containerd

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	l "github.com/rancher/k3d/v5/pkg/logger"
	k3d "github.com/rancher/k3d/v5/pkg/types"
	"strconv"
	"strings"
)

func containerToK3DNode(c containerd.Container, ctx context.Context) (*k3d.Node, error) {
	info, err := c.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get container information: %w", err)
	}

	spec, err := c.Spec(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get container information: %w", err)
	}
	mem, _ := getResourceLimits(spec)
	mounts := mountsToString(spec.Mounts)
	var node = &k3d.Node{
		Name:    info.ID,
		Image:   info.Image,
		Env:     spec.Process.Env,
		Args:    spec.Process.Args,
		Volumes: mounts,
		Cmd:     []string{spec.Process.CommandLine},
		//Restart
		Created:       info.CreatedAt.String(),
		RuntimeLabels: info.Labels,
		K3sNodeLabels: info.Labels,
		//Networks
		//ExtraHosts
		//ServerOpts
		//AgentOpts
		//GPURequest
		Memory: mem,
		//State
		//IP
		//HookActions,
		//Role:  nil,
	}
	return node, nil
}

func getResourceLimits(spec *oci.Spec) (string, string) {
	var men, cpu string
	if spec.Linux != nil {
		men = byteToStr(*spec.Linux.Resources.Memory.Limit)
		cpu = strconv.FormatInt(*spec.Linux.Resources.CPU.Quota, 10)
	} else if spec.Windows != nil {
		men = byteToStr(int64(*spec.Windows.Resources.Memory.Limit))
		cpu = strconv.FormatInt(int64(*spec.Windows.Resources.CPU.Maximum), 10)
	} else if spec.Solaris != nil {
		// ¯\_(ツ)_/¯
		men = spec.Solaris.MaxShmMemory
		cpu = spec.Solaris.CappedCPU.Ncpus
	} else {
		//Quack
		l.Log().Infof("Unknow Container Spec not Linux,Windows or Solaris")
	}
	return men, cpu
}

func mountsToString(mounts []specs.Mount) []string {
	var mount []string
	for _, mnt := range mounts {
		mount = append(mount, fmt.Sprintf("%s %s %s %s", mnt.Type, mnt.Source, mnt.Destination, mnt.Options))
	}
	return mount
}

func byteToStr(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGT"[exp])
}

func toContainerdFilter(m map[string]string) string {
	var filters []string
	for key, value := range m {
		filters = append(filters, fmt.Sprintf("labels.%s==%s", key, value))
	}
	return strings.Join(filters, ",")
}

func containersToNodes(containers []containerd.Container, ctx context.Context) []*k3d.Node {
	var nodes []*k3d.Node
	for _, container := range containers {
		node, err := containerToK3DNode(container, ctx)
		if err != nil {
			return nil
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func k3dNodeToContainerD(client *containerd.Client, ctx context.Context, img containerd.Image, node *k3d.Node) (containerd.Container, error) {
	container, err := client.NewContainer(ctx, node.Name,
		containerd.WithImage(img),
		containerd.WithContainerLabels(node.K3sNodeLabels),
		containerd.WithNewSnapshot(fmt.Sprintf("%s-snapshot", node.Name), img),
		containerd.WithNewSpec(
			oci.WithEnv(node.Env),
			oci.WithImageConfig(img)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create container: %w", err)
	}

	return container, nil
}
