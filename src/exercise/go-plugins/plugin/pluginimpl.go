package plugin

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

/*
HtmlizePlugin实现了plugin.Plugin接口，以提供RPC服务器或客户端返回到插件机械。服务器端应该用Htmlizer接口的具体实现来证明Impl字段。
*/
type HtmlizePlugin struct {
	Impl Htmlizer
}

func (p *HtmlizePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &PluginServerRPC{
		Impl: p.Impl,
	}, nil
}

func (p *HtmlizePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &PluginClientRPC{client: c}, nil
}
