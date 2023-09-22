package boilerplate

import (
	"embed"
)

// DefaultBoilerplateFs 新的默认模版需要在这里添加
//
//go:embed default.toml
//go:embed git/java.gitignore
//go:embed git/javascript.gitignore
//go:embed git/python.gitignore
//go:embed git/golang.gitignore
//go:embed plantuml/class.puml
//go:embed plantuml/mindmap.puml
//go:embed prettier/prettierignore.ignore
//go:embed prettier/prettierrc.json
var DefaultBoilerplateFs embed.FS

//go:embed custom.toml
var CustomConfigExample []byte
