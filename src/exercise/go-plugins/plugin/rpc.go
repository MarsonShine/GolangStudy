package plugin

import (
	"go-plugins/content"
	"log"
	"net/rpc"
)

type HooksArgs struct{}

type HooksReply struct {
	Hooks []string
}

type ContentsArgs struct {
	Value string
	Post  content.Post
}

type ContentsReply struct {
	Value string
}

type RoleArgs struct {
	Role  string
	Value string
	Post  content.Post
}

type RoleReply struct {
	Value string
}

/*
PluginServerRPC被插件用来映射来自客户端的RPC调用到Htmlizer接口的方法。
*/
type PluginServerRPC struct {
	Impl Htmlizer
}

func (s *PluginServerRPC) Hooks(args HooksArgs, reply *HooksReply) error {
	reply.Hooks = s.Impl.Hooks()
	return nil
}

func (s *PluginServerRPC) ProcessContents(args ContentsArgs, reply *ContentsReply) error {
	reply.Value = s.Impl.ProcessContents(args.Value, args.Post)
	return nil
}

func (s *PluginServerRPC) ProcessRole(args RoleArgs, reply *RoleReply) error {
	reply.Value = s.Impl.ProcessRole(args.Role, args.Value, args.Post)
	return nil
}

/*
PluginClientRPC被客户端（主程序）用来将插件的Htmlize接口翻译成RPC调用。
*/
type PluginClientRPC struct {
	client *rpc.Client
}

func (c *PluginClientRPC) Hooks() []string {
	var reply HooksReply
	if err := c.client.Call("Plugin.Hooks", HooksArgs{}, &reply); err != nil {
		log.Fatal(err)
	}
	return reply.Hooks
}

func (c *PluginClientRPC) ProcessContents(val string, post content.Post) string {
	var reply ContentsReply
	if err := c.client.Call(
		"Plugin.ProcessContents",
		ContentsArgs{Value: val, Post: post},
		&reply); err != nil {
		log.Fatal(err)
	}
	return reply.Value
}

func (c *PluginClientRPC) ProcessRole(role string, val string, post content.Post) string {
	var reply RoleReply
	if err := c.client.Call(
		"Plugin.ProcessRole",
		RoleArgs{Role: role, Value: val, Post: post},
		&reply); err != nil {
		log.Fatal(err)
	}
	return reply.Value
}
