/*
Copyright © 2020-2021 The k3d Author(s)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package containerd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	runtimeTypes "github.com/rancher/k3d/v5/pkg/runtimes/types"
	k3d "github.com/rancher/k3d/v5/pkg/types"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Containerd struct{}

func (c Containerd) GetHost() string {
	panic("implement me")
}

func (c Containerd) CreateNode(ctx context.Context, node *k3d.Node) error {
	if checkNamespace(node) != nil {
		return fmt.Errorf("unable to create namespace")
	}
	namespacedCtx := namespaces.WithNamespace(ctx, NAMESPACE)
	client, err := getContainerdClient()

	if err != nil {
		return fmt.Errorf("unable to list containers: %w", err)
	}
	defer client.Close()
	img, err := client.Pull(namespacedCtx, node.Image, containerd.WithPullUnpack)
	if err != nil {
		return fmt.Errorf("unable to pull image %s: %w", node.Image, err)
	}
	_, err = k3dNodeToContainerD(client, namespacedCtx, img, node)
	if err != nil {
		return fmt.Errorf("unable to create container: %w", err)
	}
	return nil
}

func checkNamespace(node *k3d.Node) error {
	ctx := context.Background()
	client, _ := getContainerdClient()
	ns, err := client.NamespaceService().List(ctx)
	if err != nil {
		return err
	}
	for i := 0; i < len(ns); i++ {
		if NAMESPACE == ns[i] {
			return nil
		}
	}
	err = client.NamespaceService().Create(ctx, NAMESPACE, node.K3sNodeLabels)
	defer client.Close()
	return err
}

func (c Containerd) DeleteNode(ctx context.Context, node *k3d.Node) error {
	panic("implement me")
}

func (c Containerd) RenameNode(ctx context.Context, node *k3d.Node, s string) error {
	panic("implement me")
}

func (c Containerd) GetNodesByLabel(ctx context.Context, m map[string]string) ([]*k3d.Node, error) {
	var nodes []*k3d.Node
	filters := toContainerdFilter(m)
	namespacedCtx := namespaces.WithNamespace(ctx, NAMESPACE)
	client, err := getContainerdClient()
	defer client.Close()
	if err != nil {
		return nil, fmt.Errorf("containerd failed to provide info output: %w", err)
	}
	containers, err := client.Containers(namespacedCtx, filters)
	if err != nil {
		return nil, fmt.Errorf("unable to list containers: %w", err)
	}

	nodes = containersToNodes(containers, namespacedCtx)
	return nodes, nil
}

func (c Containerd) GetNode(ctx context.Context, node *k3d.Node) (*k3d.Node, error) {
	return nil, nil
}

func (c Containerd) GetNodeStatus(ctx context.Context, node *k3d.Node) (bool, string, error) {
	panic("implement me")
}

func (c Containerd) GetNodesInNetwork(ctx context.Context, s string) ([]*k3d.Node, error) {
	panic("implement me")
}

func (c Containerd) GetKubeconfig(ctx context.Context, node *k3d.Node) (io.ReadCloser, error) {
	panic("implement me")
}

func (c Containerd) StartNode(ctx context.Context, node *k3d.Node) error {

	panic("implement me")
}

func (c Containerd) StopNode(ctx context.Context, node *k3d.Node) error {
	panic("implement me")
}

func (c Containerd) CreateVolume(ctx context.Context, s string, m map[string]string) error {
	cwd, _ := os.Getwd()
	volPath := filepath.Join(cwd, ".volumes", s)
	_, err := os.Stat(volPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(volPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Containerd) DeleteVolume(ctx context.Context, s string) error {
	panic("implement me")
}

func (c Containerd) GetVolume(s string) (string, error) {
	panic("implement me")
}

func (c Containerd) ExecInNode(ctx context.Context, node *k3d.Node, strings []string) error {
	panic("implement me")
}

func (c Containerd) ExecInNodeGetLogs(ctx context.Context, node *k3d.Node, strings []string) (*bufio.Reader, error) {
	panic("implement me")
}

func (c Containerd) GetNodeLogs(ctx context.Context, node *k3d.Node, time time.Time, opts *runtimeTypes.NodeLogsOpts) (io.ReadCloser, error) {
	panic("implement me")
}

func (c Containerd) GetImages(ctx context.Context) ([]string, error) {
	panic("implement me")
}

func (c Containerd) CopyToNode(ctx context.Context, s string, s2 string, node *k3d.Node) error {
	panic("implement me")
}

func (c Containerd) WriteToNode(ctx context.Context, bytes []byte, s string, mode os.FileMode, node *k3d.Node) error {
	panic("implement me")
}

func (c Containerd) ReadFromNode(ctx context.Context, s string, node *k3d.Node) (io.ReadCloser, error) {
	panic("implement me")
}

func (c Containerd) ConnectNodeToNetwork(ctx context.Context, node *k3d.Node, s string) error {
	panic("implement me")
}

func (c Containerd) DisconnectNodeFromNetwork(ctx context.Context, node *k3d.Node, s string) error {
	panic("implement me")
}

func getContainerdClient() (*containerd.Client, error) {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return nil, fmt.Errorf("unable to connect to containerd stock: %w", err)
	}
	return client, nil
}
