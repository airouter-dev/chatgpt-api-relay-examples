# x402 付费生图调用

AI ROUTER 的生图接口会先返回机器可读的 HTTP 402 付款挑战，不会在第一
次请求时扣款。x402 SDK 会校验挑战、签名精确的 EVM 授权，再携带签名重试
原请求。客户端不需要保存商户或 PayAI 密钥。

## 已验证的生产契约

当前使用 Base 主网 USDC，支持 `gpt-image-2`、`gpt-image-2.5`、
`gpt-image-2.5-sunburst`、`gpt-image-2.5-flare`。`n` 必须显式传入，范围
为 1–4；金额是 `4500 × n` 个 USDC 最小单位，即每张 `$0.0045`。

生成接口使用 JSON；编辑接口使用 multipart 表单并上传 `image` 文件。
Go、Python、Node.js 示例均从 `.env` 或 TOML 读取配置，默认只探测不付款。

## 私钥安全

请使用余额很小的专用付款钱包。`AI_ROUTER_X402_EVM_PRIVATE_KEY` 只能放在
本地被 Git 忽略的 `.env` 或密钥管理器中，不能写入源代码、命令历史、Issue
或日志。先运行 probe 模式检查金额，再显式选择 pay 模式。

访问 [中文主页](https://ai-router.dev/cn) 了解账户和 API Key；也可查看
[英文主页](https://ai-router.dev) 与 [俄文主页](https://ai-router.dev/ru)。
