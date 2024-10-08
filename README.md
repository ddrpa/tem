# tem: 样板文件注入器

# 支持列表

- `Dockerfile` 模版，有 rootless 和 rootful 版本
- `Jenkinsfile`
- `.gitignore` 文件（Go / Java / JavaScript / Python）
- PlantUML Class 模版和 MindMap 模版
- Prettier 配置和 ignore 文件

输入小写字母的名称，自动补全将会列出可选项。

# 简介

很多时候我们需要快速创建一些固定的样板文件：

- 补充特定语言的 gitignore 文件
- 添加项目级的配置文件（Maven / Prettier / Webpack 等）
- 从模版开始编写代码（PlantUML / Jenkinsfile / Dockerfile / SystemD Unit 等等）
- ……

大部分情况下我们只是需要把项目中的某个默认文件替换为自己的版本，或者添加一两个文件。
使用 yeoman 这样的工具创建工程模版似乎有点小题大做。
2020 年我从找一个以前的项目拷贝文件进化到使用一个 Zsh 脚本从 `$HOME/public/`
下拷贝样板文件到当前目录。这个脚本还支持了自动补全，不过每次新增文件就要修改，并不是特别方便。

```shell
# 其中一段截取出来长这样
...
UML_MINDMAP="$CT_TEMPLATES_PATH/plantuml/mindmap"
PRETTIER_IGNORE="$CT_TEMPLATES_PATH/prettier/prettierignore"
...
case $1 {
    # gitignore
    (gitignore-java)
    cp $GITIGNORE_JAVA .gitignore
    ;;
    (gitignore-javascript);&
    (gitignore-js)
    cp $GITIGNORE_JAVASCRIPT .gitignore
    ;;
...
```

Tem 是我的最新方案，`tem add gitignore-java -y` 这样的命令将在当前目录下创建一个 `.gitignore` 文件，
相较 Spring Initializr 提供的 `.gitignore` 添加了不少符合我本人以及团队实践习惯的排除项。

## 安装后的首次使用与升级

从 release 页面下载二进制文件，解压后放到 `$PATH` 下即可使用。

首次安装或升级后需要执行 `tem init`。
tem 会检查 `$HOME/.tem/` 目录，并替换该目录下 `default/` 的内容：

```shell
$ tree /Users/yufan/.tem

├── custom
│   └── custom.toml
└── default
    ├── default.toml
    ├── git
    │   ├── java.gitignore
    ...
```

如果用户还没有编写自定义配置，`tem init` 也会创建一个 `custom/custom.toml` 文件作为示例。

`default/` 目录预置了一些样板文件，`default/default.toml` 描述了复制这些样板文件需要输入的命令。
用户自定义配置也遵循相同的语法。

### 配置自动补全

tem 通过 `github.com/spf13/cobra` 实现了强大的自动补全功能，支持 bash、Zsh、fish、PowerShell 等 shell。
以 Zsh 为例，执行如下命令，在新的会话中即可使用 <kbd>Tab</kbd> 补全。

```shell
$ tem completion zsh > "$ZSH_CACHE_DIR/completions/_tem"
```

### 使用方法

`add` 命令用于将模板文件写入以当前工作目录为起始点的预设相对路径。

在输入 `tem add $key -y` 的 `key` 参数时，可以使用 <kbd>Tab</kbd> 进行自动补全。如果不添加 `-y` 参数，tem 会告诉你将要执行什么操作，但不会把文件写入磁盘。

## 使用自定义样板

`tem init` 生成的示例 `custom/custom.toml` 文件为例（移除了注释符号）：

```toml
[[template]]
key = "maven-wrapper"
assets = [
    ["custom/maven/mvnw", "mvnw"],
    ["custom/maven/mvnw.cmd", "mvnw.cmd"],
    ["custom/maven/maven-wrapper.jar", ".mvn/wrapper/maven-wrapper.jar"],
    ["custom/maven/maven-wrapper.properties", ".mvn/wrapper/maven-wrapper.properties"],
]

[[template]]
key = "golang"
assets = [
    ["https://www.toptal.com/developers/gitignore/api/goland+all,macos,linux,windows,visualstudiocode", ".gitignore"]
]
```

可以使用 `tem add maven-wrapper -y`（因为有自动补全，`key` 以方便记忆和识别为佳）创建 maven-wrapper 需要的样板文件。

`assets` 中的每一项都是一个数组，第一个元素是源文件，第二个元素是目标文件。
源文件的路径是 `custom/` 目录下的相对路径，目标文件的路径是当前目录下的相对路径。
tem 会自动创建目标文件的父目录。

`assets` 中的源文件可以是本地文件，也可以是网络文件，网络文件目前支持 HTTP 和 HTTPS。

### 添加子配置

tem 现在支持子配置，在读取 `$TEM_HOME/custom/custom.toml` 时，会检查 `[config]` 下的 `imports` 字段，载入其他配置文件。

模版配置的加载使用深度优先的方案，即读取到子配置文件就先加载子配置文件中的内容，因此需要将可能重名且优先度较高的的模版 key 放在加载队列的后方。
