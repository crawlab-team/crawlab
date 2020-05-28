package local_node

import (
	"crawlab/model"
	"github.com/spf13/viper"
)

func GetLocalNode() *LocalNode {
	return localNode
}
func CurrentNode() *model.Node {
	return GetLocalNode().Current()
}

func InitLocalNodeInfo() (err error) {
	registerType := viper.GetString("server.register.type")
	ip := viper.GetString("server.register.ip")
	customNodeName := viper.GetString("server.register.customNodeName")

	localNode, err = NewLocalNode(ip, customNodeName, registerType)
	if err != nil {
		return err
	}
	if model.IsMaster() {
		err = model.UpdateMasterNodeInfo(localNode.Identify, localNode.Ip, localNode.Mac, localNode.Hostname)

		if err != nil {
			return err
		}
	}
	return localNode.Ready()
}
