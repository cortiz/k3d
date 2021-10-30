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
	"strings"
	"time"
)

type Containerd struct{}

func (c Containerd) GetHost() string {
	panic("implement me")
}

func (c Containerd) CreateNode(ctx context.Context, node *k3d.Node) error {

	panic("implement me")
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
	namespacedCtx := namespaces.WithNamespace(ctx, namespaces.Default)
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

func (c Containerd) GetNode(ctx context.Context, node *k3d.Node) (*k3d.Node, error) {
	panic("implement me")
}

func (c Containerd) GetNodeStatus(ctx context.Context, node *k3d.Node) (bool, string, error) {
	panic("implement me")
}

func (c Containerd) GetNodesInNetwork(ctx context.Context, s string) ([]*k3d.Node, error) {
	panic("implement me")
}

func (c Containerd) CreateNetworkIfNotPresent(ctx context.Context, network *k3d.ClusterNetwork) (*k3d.ClusterNetwork, bool, error) {
	panic("implement me")
}

func (c Containerd) GetKubeconfig(ctx context.Context, node *k3d.Node) (io.ReadCloser, error) {
	panic("implement me")
}

func (c Containerd) DeleteNetwork(ctx context.Context, s string) error {
	panic("implement me")
}

func (c Containerd) StartNode(ctx context.Context, node *k3d.Node) error {

	panic("implement me")
}

func (c Containerd) StopNode(ctx context.Context, node *k3d.Node) error {
	panic("implement me")
}

func (c Containerd) CreateVolume(ctx context.Context, s string, m map[string]string) error {
	panic("implement me")
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
