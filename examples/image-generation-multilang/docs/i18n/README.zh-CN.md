# 多语言 AI 生图 API 示例：从文生图到图生图（Python、Node.js、C++、Go、Rust、PHP）

这份中文指南介绍一个可阅读、可运行、便于扩展的多语言图片生成示例项目。六种实现都演示相同的核心工作流：先把文字提示发送到 OpenAI-compatible Images API 的 `/images/generations` 端点；收到图片后将其保存为本地产物；再把该文件作为输入提交给 `/images/edits`，依据第二条指令生成修改后的图片。你可以逐步查看请求、复用客户端模块，也可以使用同一组提示词比较不同语言的集成方式。

本项目面向开发者学习和集成验证，不是 OpenAI 官方 SDK，也不是一份对所有环境开箱即用的生产服务。示例默认连接独立服务 [AI-ROUTER](https://ai-router.dev/cn)，但 API 根地址可以配置到其他实现了兼容接口的网关。模型、尺寸、配额、价格和响应能力由你的账户及上游配置决定；代码中出现某个模型名称，不代表该模型对每个账户均开放。

## 你能用这个项目做什么

- **验证完整的图像创作循环：**以文本生成初稿，保留初稿，再用实际图片文件进行图像编辑。编辑并非第二次纯文本生成，而是带有输入图片的 multipart 请求。
- **学习通用 Images API 集成：**查看 Bearer 认证、JSON 请求、multipart 文件上传、响应解析、本地保存和错误报告如何组织。
- **对照不同生态的工程结构：**比较标准库客户端、异步客户端、原生 HTTP 库和命令行解析方案，而不是比较不同的业务流程。
- **作为应用原型的起点：**将客户端嵌入自己的 API 服务、队列消费者、批处理任务、桌面工具或内容制作流程；对外开放前仍需补齐鉴权、限额和内容治理。

每次完整工作流至少会请求两次图像 API，因此可能分别产生生成和编辑费用。正式运行前，请先确认账户、模型权限、价格和上游保留策略。

## 六种语言如何分层

六个 demo 共享端点与业务顺序，但按各语言惯例组织配置、HTTP 客户端、工作流和 CLI。让 CLI 只处理参数、客户端只处理协议、工作流只连接业务步骤，可以减少未来更换存储或增加语言时的交叉修改。

| 语言 | 代码组织与职责 | 本项目中的实际入口 | 适合观察的特点 |
| --- | --- | --- | --- |
| **Python** | `image_demo/config.py` 读取环境与命令行覆盖；`client.py` 使用标准库 `urllib` 处理 JSON、multipart、响应解码和 URL 图片下载；`workflow.py` 串联生成、保存、编辑、再保存；`cli.py` 负责参数和终端输出。 | 从 `python/` 目录运行 `python main.py --prompt ... --edit-prompt ...`。一次命令完成两步工作流。 | 零运行时依赖、易读，适合脚本、Notebook、批处理和快速原型。 |
| **Node.js** | `src/config.mjs` 统一设置；`src/client.mjs` 封装内置 `fetch` 与图片解析；`src/workflow.mjs` 管理两步流程；`src/cli.mjs` 处理 CLI；`test/` 使用内置 `node:test`。 | 从 `nodejs/` 目录运行 `npm run demo -- --prompt ... --edit-prompt ...`。 | 无第三方 npm 运行时依赖，适合 JavaScript 后端、队列和无服务器运行时。 |
| **C++** | `include/image_client.hpp` 定义可复用接口；`src/image_client.cpp` 用 libcurl 处理网络、multipart、解码与超时；`src/main.cpp` 解析 CLI 并执行工作流；CMake 单独构建库、可执行文件与离线测试。 | 在 `cpp/` 配置 CMake 后运行 `imagegen --prompt ... --edit-prompt ...`。 | 适合原生应用和受控运行环境；展示 C++17、libcurl 与 CMake 的边界划分。 |
| **Go** | `internal/config` 读取环境默认值；`internal/client` 使用 `net/http`、`encoding/json` 和 multipart；`cmd/imagegen` 解析 flag、保存产物并协调两次请求；`client_test.go` 用 `httptest` 验证协议。 | 从 `go/` 目录运行 `go run ./cmd/imagegen --prompt ... --edit-prompt ...`。 | 标准库实现简洁，适合小型服务、命令行工具和后台 worker。 |
| **Rust** | `config.rs` 加载密钥与 URL；`client.rs` 使用 `reqwest` / `serde`；`types.rs` 描述响应；`workflow.rs` 编排操作；`artifacts.rs` 写入文件；`main.rs` 通过 `clap` 提供 `generate`、`edit`、`workflow` 子命令。 | 从 `rust/` 目录运行 `cargo run --release -- workflow ...`，也能单独调用 generate 与 edit。 | 展示显式类型、异步 HTTP 和把传输、文件系统、业务组合分离的做法。 |
| **PHP** | `Config` 管理环境设置；`ImageClient` 处理 cURL 和 Images API；`ImageData` 归一化响应；`ArtifactStore` 写文件；`Workflow` 串联两步；`bin/image-demo.php` 负责 CLI。 | 从 `php/` 目录运行 `php bin/image-demo.php workflow ...`，也提供独立的 `generate`、`edit`。 | 无必须安装的 Composer 运行时依赖；需要 PHP cURL 与 JSON 扩展，适合 PHP 网站后台与 CLI 作业。 |

所有实现都先保存生成结果，再发起 edit 请求。因此如果编辑失败，生成初稿仍在磁盘上，便于检查、重试或手动复用。架构的更多说明见[架构指南](../ARCHITECTURE.md)；新增语言时可按[贡献指南](../../CONTRIBUTING.md)保持协议和工作流一致。

## API 契约：文字生成图片，然后编辑图片

默认 API 根地址为 `https://api.ai-router.dev/v1`。这里的 `{base_url}` 指包含 `/v1` 的根地址，客户端会在其后拼接端点。真实请求使用账户 API key 作为 Bearer token。

### 1. 文生图：`POST /images/generations`

请求是 `application/json`。下面是便于理解的 OpenAI-compatible 请求示例；具体可选字段由兼容服务和账户能力决定：

```http
POST https://api.ai-router.dev/v1/images/generations
Authorization: Bearer $AI_ROUTER_API_KEY
Content-Type: application/json
```

```json
{
  "model": "gpt-image-2",
  "prompt": "雨后屋顶上的小型玻璃温室，清晨柔和逆光，编辑摄影风格",
  "size": "1024x1024",
  "n": 1,
  "response_format": "b64_json"
}
```

示例客户端从响应 `data` 数组的第一项解析图片。支持的实现可处理 `b64_json`，并在网关返回图片 `url` 时下载该图片。优先使用 base64 的好处是工作流不需要额外依赖一个匿名下载地址；具体响应格式以当前模型和账户为准。

### 2. 图生图：`POST /images/edits`

编辑使用 `multipart/form-data`。除 `model`、`prompt`、可选的 `size` 与 `response_format` 外，还必须以 `image` 字段附上刚才保存的图片：

| 字段 | 用途 |
| --- | --- |
| `model` | 当前账户可用、且支持编辑的图像模型 ID。 |
| `prompt` | 描述希望怎样修改图片；也可以明确写出希望保留的主体或构图。 |
| `size` | 可选输出尺寸，例如 `1024x1024`；是否可用取决于模型和账户。 |
| `response_format` | 示例默认使用 `b64_json`；服务也可能按账户配置返回 URL。 |
| `image` | 必需的输入图片文件。此项目把文生图步骤写入磁盘的图片作为该字段上传。 |

成功响应沿用 `data` 图片结果结构。完整的跨语言字段说明见[API 契约文档](../API_CONTRACT.md)。

### 模型、尺寸与“OpenAI-compatible”的含义

`AI_ROUTER_MODEL` 或 `--model` 是传给网关的模型 ID，不是本地模型实现。可尝试 `gpt-image-2`、`gpt-image-2.5`、`gpt-image-2.5-flare`、`gpt-image-2.5-burst` 等名称，前提是它们确实出现在你的账户模型列表中，并且支持相应操作。示例没有绕过账户权限，也不承诺特定名称、尺寸、质量等级、价格或可用性长期不变。选择尺寸时，优先采用所选模型文档明确支持的值；`1024x1024` 是示例默认值，不是所有模型的保证。

“OpenAI-compatible”表示客户端采用相似的端点和数据交换形状，便于迁移已有集成思路；它不表示 AI-ROUTER 是 OpenAI，亦不表示所有可选参数、模型或行为完全一致。请以所使用服务的当前文档为准。

## 环境变量和配置优先级

仓库根目录的 `.env.example` 仅供参考。每个 CLI 并不会自动加载 `.env`；请在启动进程的 shell 中设置变量，或使用你信任的 dotenv 工具。真实 key 不要提交到 Git。

通用起点：

```bash
export AI_ROUTER_API_KEY="替换为你自己的密钥"
export AI_ROUTER_BASE_URL="https://api.ai-router.dev/v1"
export AI_ROUTER_MODEL="gpt-image-2"
export AI_ROUTER_SIZE="1024x1024"
```

不同 demo 对别名的支持有差异，以下是源码实际读取的配置：

| 示例 | API Key | Base URL | 模型 / 尺寸 / 超时 |
| --- | --- | --- | --- |
| Python | `IMAGE_API_KEY`、`OPENAI_API_KEY`、`AI_ROUTER_API_KEY`，按顺序回退 | `IMAGE_API_BASE_URL`、`AI_ROUTER_BASE_URL` | 模型 `IMAGE_MODEL` / `AI_ROUTER_MODEL`；尺寸 `IMAGE_SIZE` / `AI_ROUTER_SIZE`；质量 `IMAGE_QUALITY`；超时 `IMAGE_TIMEOUT_SECONDS`（秒） |
| Node.js | `IMAGE_API_KEY`、`OPENAI_API_KEY`、`AI_ROUTER_API_KEY` | `IMAGE_API_BASE_URL`、`AI_ROUTER_BASE_URL` | `IMAGE_MODEL` / `AI_ROUTER_MODEL`；`IMAGE_SIZE` / `AI_ROUTER_SIZE`；`IMAGE_QUALITY`；`IMAGE_TIMEOUT_MS`（毫秒） |
| C++ | `OPENAI_API_KEY`、`AI_ROUTER_API_KEY` | `OPENAI_BASE_URL`、`AI_ROUTER_BASE_URL` | `IMAGE_MODEL` / `AI_ROUTER_MODEL`；`IMAGE_SIZE` / `AI_ROUTER_SIZE`；`IMAGE_RESPONSE_FORMAT`；`IMAGE_API_TIMEOUT` / `AI_ROUTER_IMAGE_TIMEOUT`（秒） |
| Go | `OPENAI_API_KEY`、`AI_ROUTER_API_KEY` | `OPENAI_BASE_URL`、`AI_ROUTER_BASE_URL` | `IMAGE_MODEL` / `AI_ROUTER_MODEL`；`IMAGE_SIZE` / `AI_ROUTER_SIZE`；`IMAGE_RESPONSE_FORMAT`；`IMAGE_API_TIMEOUT` / `AI_ROUTER_IMAGE_TIMEOUT`（秒） |
| Rust | `AI_ROUTER_API_KEY`、`OPENAI_API_KEY` | `AI_ROUTER_BASE_URL` | `AI_ROUTER_TIMEOUT_SECONDS`；模型和尺寸由子命令选项提供，可从 `AI_ROUTER_MODEL`、`AI_ROUTER_SIZE` 读取 |
| PHP | `AI_ROUTER_API_KEY`、`OPENAI_API_KEY` | `AI_ROUTER_BASE_URL` | 超时 `AI_ROUTER_TIMEOUT_SECONDS`；模型和尺寸通过 CLI 选项或 CLI 默认环境读取 |

命令行选项通常覆盖环境默认值。注意 PHP 不接受 `--api-key`，请通过环境传密钥；Rust 的 CLI 主 URL 选项是 `--base-url`，模型和尺寸是每个子命令自己的选项。

## 安装和运行：按语言选择

下面的路径均相对于仓库根目录。一次工作流会发送两次 API 请求；如暂时只想检查编译、测试或参数解析，请先运行各自的本地测试，不必配置真实 key。

### Python（3.10+）

实现只使用 Python 标准库，没有运行时依赖。可用虚拟环境隔离开发：

```bash
cd marketing/image-generation-multilang/python
python -m venv .venv
```

激活环境（按你使用的 shell 选择一组）：

```bash
# macOS / Linux
source .venv/bin/activate

# Windows PowerShell
.\.venv\Scripts\Activate.ps1
```

设置密钥并运行完整两步流程：

```bash
export IMAGE_API_KEY="你的密钥"
python main.py \
  --prompt "一座雨后屋顶上的玻璃温室，电影感编辑摄影" \
  --edit-prompt "将光线改成温暖的日落色调，同时保留温室结构" \
  --output-dir ./outputs
```

PowerShell 中可用 `$env:IMAGE_API_KEY = "你的密钥"`，其余参数按 PowerShell 换行规则传入。可选地运行 `python -m pip install -e .` 安装 `image-demo` 命令；也可验证 `python -m unittest discover -s tests -v`。成功后默认目录含 `generated` 和 `edited` 两个图片文件；扩展名依响应文件类型而定。

### Node.js（18.17+）

此目录不依赖第三方 npm 包。安装 Node.js 后从项目目录运行：

```bash
cd marketing/image-generation-multilang/nodejs
npm test
```

`npm test` 是不调用真实 API 的本地测试。执行两阶段工作流：

```bash
export AI_ROUTER_API_KEY="你的密钥"
npm run demo -- \
  --prompt "一座雨后屋顶上的玻璃温室，电影感编辑摄影" \
  --edit-prompt "将光线改成温暖的日落色调，同时保留温室结构" \
  --model gpt-image-2
```

PowerShell 使用 `$env:AI_ROUTER_API_KEY = "你的密钥"`。默认输出到 `nodejs/outputs/generated.png` 与 `nodejs/outputs/edited.png`；两者都先后写入本地，编辑请求以上一个文件为输入。`--size`、`--quality`、`--base-url`、`--timeout`（毫秒）、`--output-dir` 均可按 CLI 用法调整。

### C++17（CMake 3.20+ 和 libcurl）

安装 C++17 编译器、CMake 和 libcurl 开发文件。Windows vcpkg 用户可安装依赖：

```powershell
vcpkg install curl:x64-windows
```

配置时按本机工具链添加 vcpkg toolchain 参数；通用单配置构建命令如下：

```bash
cd marketing/image-generation-multilang/cpp
cmake -S . -B build -DBUILD_TESTING=ON
cmake --build build --config Release
ctest --test-dir build --output-on-failure
```

设置 `AI_ROUTER_API_KEY` 后执行（Linux/macOS 单配置示例）：

```bash
export AI_ROUTER_API_KEY="你的密钥"
./build/imagegen \
  --base-url https://api.ai-router.dev/v1 \
  --model gpt-image-2 \
  --prompt "一座雨后屋顶上的玻璃温室" \
  --edit-prompt "加入温暖的金色夕照，保持建筑主体不变" \
  --output generated.png --edited-output edited.png
```

Visual Studio 多配置生成器通常将程序放在 `build/Release/imagegen.exe`。离线 CTest 不需要 API key；示例命令会按指定路径保存两张图。CLI 也支持 `--size`、`--response-format`、`--timeout` 和其他输出路径。

### Go（1.22+）

该实现使用 Go 标准库，无需 `go get` 添加额外依赖：

```bash
cd marketing/image-generation-multilang/go
go test ./...
go vet ./...
```

运行端到端生成与编辑：

```bash
export AI_ROUTER_API_KEY="你的密钥"
go run ./cmd/imagegen \
  --base-url https://api.ai-router.dev/v1 \
  --model gpt-image-2 \
  --prompt "一座雨后屋顶上的玻璃温室" \
  --edit-prompt "加入温暖的金色夕照，保持建筑主体不变" \
  --output generated.png --edited-output edited.png
```

默认产物为当前工作目录的 `generated.png` 和 `edited.png`。`go test` 会用本地 `httptest` 服务检查 HTTP 请求和响应处理，不会产生真实生成请求或费用。

### Rust（稳定版工具链）

Cargo 根据 `Cargo.toml` 获取所需依赖，首次构建可能需要网络访问 crates.io：

```bash
cd marketing/image-generation-multilang/rust
export AI_ROUTER_API_KEY="你的密钥"
```

通过一个子命令先文生图、保存、再 edit：

```bash
cargo run --release -- workflow \
  --prompt "雨后屋顶上的玻璃温室" \
  --edit-prompt "加入温暖的悬挂灯，保留原有构图" \
  --model gpt-image-2 --size 1024x1024 \
  --generated output/generated.png \
  --edited output/edited.png
```

输出父目录会自动创建。`cargo run -- generate --prompt ... --output ...` 和 `cargo run -- edit --image ... --prompt ... --output ...` 也可分开调用。`cargo test` 可用于检查编译与测试目标；Rust 客户端请求使用 `b64_json`，并能解析 URL 图片响应。

### PHP（8.1+、cURL 和 JSON 扩展）

项目包含轻量 autoloader，因此 demo 运行本身不要求 Composer 安装包。先确认扩展可用，再运行本地 CLI：

```bash
cd marketing/image-generation-multilang/php
php -m
export AI_ROUTER_API_KEY="你的密钥"
```

完整工作流命令：

```bash
php bin/image-demo.php workflow \
  --prompt "雨后屋顶上的玻璃温室" \
  --edit-prompt "加入温暖的悬挂灯，保留原有构图" \
  --model gpt-image-2 --size 1024x1024 \
  --generated output/php-generated.png \
  --edited output/php-edited.png
```

也可以通过 `generate` 和 `edit` 子命令分开操作。默认输出位于 `output/php-generated.png` 和 `output/php-edited.png`；`AI_ROUTER_BASE_URL` 可修改 API 根地址。密钥只能从 `AI_ROUTER_API_KEY` 或 `OPENAI_API_KEY` 环境变量读取。若要把包接入 Composer 管理的应用，可再执行 `composer dump-autoload` 并使用 `composer.json` 中的 PSR-4 映射。

## 结果文件和执行成本

工作流会保存两种产物：第一张是从 prompt 生成的初稿；第二张是把初稿上传至 edit 端点后得到的编辑结果。文件名和目录由语言示例的默认值或命令行参数决定。根 `.gitignore` 忽略了常见输出目录，但请在运行前确认实际产物不会进入提交，以免把私有素材或客户图像发布到公开仓库。

示例没有设置自动重试：图像操作可能计费，盲目重试容易造成重复请求。生成成功但编辑失败时，请保留并检查初稿；修正问题后再有意识地单独重试 edit 步骤。

## 常见问题排查

| 现象 | 建议检查 |
| --- | --- |
| `401` / `403` 或提示缺少 API key | 确认启动程序的同一 shell 中设置了密钥；检查 key 是否有效、账户权限是否足够。注意各语言读取变量名称不同。不要将 key 加入源代码或问题日志。 |
| 模型不存在或无权限 | 检查账户当前模型列表与访问权限。模型 ID 必须准确，并确认所选模型支持 generations 和 edits 两个操作。 |
| 尺寸或质量参数无效 | 先使用当前模型文档允许的值；可从示例默认 `1024x1024` 开始，但并非每个上游都保证该值。 |
| 连接超时或网络错误 | 检查 API 根地址、网络代理、TLS 证书及超时单位；Node 的超时按毫秒，其余大多按秒。不要在不确认计费状态前无限重试。 |
| 响应中没有图片或 base64 解码失败 | 检查响应格式、网关兼容性和上游返回的错误详情。某些客户端支持 URL 回退，但 URL 下载也可能受对象存储链接时效影响。 |
| 生成成功而 edit 失败 | 确认生成文件存在、非空且可读，edit prompt 非空，并确认模型支持图片编辑。保留初稿后单独排错，不必重新生成。 |
| C++ 找不到 CURL 或链接失败 | 确认 libcurl 开发文件与 CMake 可发现的配置已安装；vcpkg 环境应向 CMake 传入 toolchain 文件。 |
| Rust 首次编译失败 | 确认安装 Rust stable 与 Cargo，且依赖下载可用；先运行 `cargo check` 查看具体依赖/编译错误。 |
| PHP 缺少 cURL 函数 | 在当前 PHP CLI 使用的 `php.ini` 启用 `curl` 扩展，命令行所用 PHP 可能与 Web PHP 配置不同。 |

## 安全、隐私与负责任使用

- 将 API key 放在环境变量或密钥管理服务中，不要提交 `.env`、命令历史或终端截图。生产服务不要把密钥放到浏览器或移动端。
- 使用 HTTPS，并确认 base URL 指向可信服务。这里只展示请求格式，服务端实际数据处理、日志、保留期限和训练政策须查阅账户服务条款与隐私文档。
- 提示词、源图、生成图可能包含个人或商业敏感信息。仅上传你有权处理的内容，并为本地图片设置适当的访问权限、留存期限和删除流程。
- 把模型返回 URL 当作外部不可信输入。用于公网服务时增加域名 allow-list、重定向限制、文件大小上限、超时、MIME 与真实格式校验，避免 SSRF 和资源耗尽。
- 公开服务需另行增加终端用户鉴权、配额与速率限制、内容审核、审计日志脱敏、请求取消、可观测性与异常恢复策略。本仓库示例没有把这些包装成现成的生产保障。
- 按上游条款、版权许可和本地法律使用模型结果；面向用户发布前进行人工质量与合规审查。

更多建议见[安全指南](../SECURITY.md)。

## 从 Demo 扩展到自己的服务

1. **把本地输出替换为对象存储。**保留 `generate → persist → edit → persist` 的明确边界，给临时产物配置访问控制、加密和生命周期规则。
2. **设计请求和成本控制。**校验 prompt 长度、上传类型和图片尺寸；按用户、模型和操作设置配额；把取消信号与 HTTP 超时传递到底层。
3. **只对可重试错误重试。**在了解服务端幂等性和计费语义后，为短暂网络错误配置有限次数的退避重试，避免重复收费。
4. **增加可观测性而不过度采集。**记录请求 ID、操作阶段、状态码和耗时；按数据最小化原则处理 prompt、图片 URL 与个人信息；绝不记录 API key。
5. **用测试固定服务契约。**扩展请求字段时同步更新 API 契约、各语言实现和离线测试；使用 mock server 验证 multipart 上传，不必每次改动都消耗真实 API 配额。
6. **逐步加入新语言或存储后端。**先实现配置、独立客户端、工作流和轻量入口，再为新边界加测试。避免把网络实现复制到 CLI，也不要假装各 provider 的扩展字段完全一致。

项目的边界和扩展路径详见[架构说明](../ARCHITECTURE.md)与[共享 API 契约](../API_CONTRACT.md)。你也可以用[示例提示词](../../examples/prompts.md)在六种实现中做较公平的流程比较。

## 相关资源

- [AI-ROUTER 中文主页](https://ai-router.dev/cn)：了解平台、账户和 API 接入入口。模型和功能以账户实际可见信息为准。
- [AI-ROUTER API 根地址](https://api.ai-router.dev/v1)：示例默认使用的 API 地址；该链接本身不替代认证或接口文档。
- [共享接口契约](../API_CONTRACT.md) · [项目架构](../ARCHITECTURE.md) · [安全建议](../SECURITY.md)
- [可复用的提示词示例](../../examples/prompts.md) · [贡献指南](../../CONTRIBUTING.md)

如果这个跨语言参考项目对你有帮助，欢迎在仓库中提出具体的兼容性问题、测试结果或新语言实现。本站链接用于了解 AI-ROUTER 服务；外链和示例不会保证特定搜索排名、流量或模型效果。
