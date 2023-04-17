package plugin

import "go-plugins/content"

/*
Htmlizer 是插件必须实现的接口。为了避免调用它不支持的角色，它必须告诉插件管理器。
它必须通过实现Hooks()方法来告诉插件管理器它想在哪些角色上被调用。
*/
type Htmlizer interface {
	// 钩子返回这个插件想要注册的钩子的列表。
	// 钩子可以有以下形式之一：
	// "contents"：该插件的ProcessContents方法将被调用到帖子的完整内容上。
	// "role:NN"：当输入中遇到:NN:role 时，该插件的ProcessRole方法将被调用，并带有role=NN和role的值。
	Hooks() []string
	// ProcessRole是对Hooks()返回的列表中的插件所要求的角色进行调用。它接收角色名称、输入的角色值和帖子，并应返回转换后的角色值。
	ProcessRole(role string, val string, post content.Post) string
	// 如果在Hooks()中要求的话，ProcessContents会被调用到整个帖子内容上。它接收内容和帖子，并应返回转换后的内容。
	ProcessContents(val string, post content.Post) string
}
