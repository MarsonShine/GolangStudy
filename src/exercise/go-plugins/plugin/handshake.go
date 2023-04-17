package plugin

import goplugin "github.com/hashicorp/go-plugin"

var Handshake = goplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "HTMLIZE_PLUGIN",
	MagicCookieValue: "hello",
}
