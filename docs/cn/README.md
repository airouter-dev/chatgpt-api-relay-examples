# AI ROUTER ChatGPT API 中转

AI ROUTER 是面向开发者的 ChatGPT API 中转服务。你可以使用兼容 OpenAI 的调用方式接入，创建 API Key，并在控制台查看用量、余额、额度和订阅状态。

中文入口：https://ai-router.dev/cn

## 接口地址

```text
https://api.ai-router.dev/v1
```

## 适合场景

- 已有 OpenAI SDK 项目，只想替换 `base_url`。
- AI 应用原型开发和调试。
- Coding Agent / Codex 类工作流。
- 短期项目按天或按周购买 ChatGPT API 额度。
- 团队需要查看 API Key 的用量和额度状态。

## 接入流程

1. 注册 AI ROUTER 账号。
2. 购买 ChatGPT API 日套餐、周套餐，或进行余额充值。
3. 创建 API Key。
4. 将 `base_url` 设置为 `https://api.ai-router.dev/v1`。
5. 按 OpenAI 兼容方式发起请求。

AI ROUTER 不是 OpenAI 官方服务。我们提供独立的大模型中转服务，调用方式兼容 OpenAI，具体套餐、模型和价格以站内展示为准。

