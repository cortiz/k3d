package containerd

import (
	"bytes"
	"context"
	"fmt"
	l "github.com/rancher/k3d/v5/pkg/logger"
	"os"
	"path/filepath"
	"text/template"

	//	"github.com/containernetworking/cni/libcni"
	//	"github.com/containernetworking/cni/pkg/types"
	//	l "github.com/rancher/k3d/v5/pkg/logger"
	k3d "github.com/rancher/k3d/v5/pkg/types"
	"net"
	//gocni "github.com/containerd/go-cni"
)

const (
	cniConfig       = ".network"
	interfacePrefix = "eth"
	defaultName     = "eth0"
	SUBNET          = "10.10.0.0/24"
	GATEWAY         = "10.10.0.1"
	//	https://github.com/AkihiroSuda/cni-isolation
	// This is the same plugin that nerdctl uses so if it is good for
	// containerd should be good for this.
	// Downside, it does require the isolation plugin to be installed, alongside cni reference plugins
	CNI_TEMPLATE = `{
   "cniVersion": "0.4.0",
   "name": "{{.NetworkName}}",
   "labels": {{.Labels}},
   "plugins": [
      {
         "type": "bridge",
         "bridge": "{{.NetworkName}}-cni",
         "isGateway": true,
         "ipMasq": true,
         "hairpinMode": true,
         "ipam": {
            "type": "host-local",
            "routes": [
               {
                  "dst": "0.0.0.0/0"
               }
            ],
            "ranges": [
               [
                  {
                     "subnet": "{{.Subnet}}",
                     "gateway": "{{.Gateway}}"
                  }
               ]
            ]
         }
      },
      {
         "type": "firewall"
      },
      {
         "type": "isolation"
      }
   ]
}`
)

type NetworkConfig struct {
	NetworkName string
	Labels      map[string]string
	Subnet      string
	Gateway     string
}

func (c Containerd) GetNetwork(ctx context.Context, network *k3d.ClusterNetwork) (*k3d.ClusterNetwork, error) {
	panic("implement me")
}

func (c Containerd) GetHostIP(ctx context.Context, network string) (net.IP, error) {
	////cni, err := gocni.New()
	//if err != nil {
	//	l.Log().Infof("Error getting Filesystem '%s'", err)
	//}

	panic("implement me")

}

func (c Containerd) DeleteNetwork(ctx context.Context, s string) error {
	panic("implement me")
}

func (c Containerd) CreateNetworkIfNotPresent(ctx context.Context, network *k3d.ClusterNetwork) (*k3d.ClusterNetwork, bool, error) {

	parse, err := template.New("").Parse(CNI_TEMPLATE)
	if err != nil {
		return nil, false, err
	}
	var buf bytes.Buffer
	err = parse.Execute(&buf, &NetworkConfig{NetworkName: network.Name,
		Subnet: SUBNET, Gateway: GATEWAY})
	if err != nil {
		l.Log().Errorf("Unable to Create network config file")
		return nil, false, err
	}
	cwd, _ := os.Getwd()
	configFilePath := checkNetworkFolder(cwd, network.Name)
	err = os.WriteFile(configFilePath, buf.Bytes(), 0644)
	if err != nil {
		return nil, false, err
	}
	return &k3d.ClusterNetwork{}, true, nil
}

func checkNetworkFolder(cwd string, name string) string {
	netFolder := filepath.Join(cwd, cniConfig)
	_, err := os.Stat(netFolder)
	if os.IsNotExist(err) {
		err := os.MkdirAll(netFolder, os.ModePerm)
		if err != nil {
			l.Log().Errorf("Unable to create Network config Folder %s", err)
			return ""
		}
	}
	return filepath.Join(netFolder, fmt.Sprintf("%s.conf", name))
}

func createCNINetworkConfig() {

}
