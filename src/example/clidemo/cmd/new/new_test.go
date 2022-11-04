package new

import "testing"

func TestRenamePackageName(t *testing.T) {
	rewrite(`E:\repositories\GolangStudy\src\example\clidemo\marsonshine\api\helloworld\v1\error_reason.pb.go`, "template", "marsonshine")
}
