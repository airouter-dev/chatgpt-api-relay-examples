# AI ROUTER: رله API چت جی پی تی

AI ROUTER یک سرویس رله API چت جی پی تی با الگوی سازگار با OpenAI برای توسعه دهندگان است. این سرویس برای اتصال برنامه ها، اسکریپت ها و جریان های کاری توسعه با API key، مشاهده مصرف، سهمیه و پلن های روزانه یا هفتگی طراحی شده است.

صفحه فارسی: https://ai-router.dev/fa

## نشانی Endpoint

```text
https://api.ai-router.dev/v1
```

## کاربردهای رایج

- اتصال کدهای موجود OpenAI SDK با تغییر `base_url`.
- آزمایش فراخوانی های ChatGPT API در برنامه ها و اسکریپت ها.
- استفاده از پلن های روزانه یا هفتگی برای پروژه های کوتاه.
- مشاهده مصرف، موجودی، سهمیه و وضعیت اشتراک API key.
- استفاده در ابزارهای داخلی و جریان های کاری coding agent.

## صفحه های موضوعی

- [رله ChatGPT API به فارسی](chatgpt-api-relay.md)
- [OpenAI-compatible endpoint](../openai-compatible-endpoint.md)
- [Usage tracking برای API key](../api-key-usage-tracking.md)
- [پلن های روزانه و هفتگی](../daily-weekly-chatgpt-api-plans.md)
- [Coding-agent workflows](../coding-agent-chatgpt-api.md)
- [FAQ](../faq.md)

## مسیر شروع

1. یک حساب AI ROUTER بسازید.
2. یک پلن روزانه یا هفتگی ChatGPT API انتخاب کنید، یا موجودی اضافه کنید.
3. یک API key ایجاد کنید.
4. مقدار `base_url` را روی `https://api.ai-router.dev/v1` بگذارید.
5. درخواست های سازگار با OpenAI ارسال کنید.

AI ROUTER سرویس رسمی OpenAI نیست. این سرویس یک relay مستقل با الگوی درخواست سازگار با OpenAI است. مدل ها، پلن ها و قیمت ها در داشبورد محصول نمایش داده می شوند.
