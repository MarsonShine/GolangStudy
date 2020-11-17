package main

// 移除文件路径和文件格式尾缀
// e.g., a => a, a.go => a, a/b/c.go => c, a/b.c.go => b.c
func basename(s string) string {
	for i := len(s); i > 0; i-- {
		if s[i] == '/' {
			s = s[i+1:]
			break
		}
	}

	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '.' {
			s = s[:i]
			break
		}
	}
	return s
}
