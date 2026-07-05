# International SEO Link Map

AI ROUTER uses English at the root path and localized subdirectories for non-English pages.

Primary SEO targets:

- English: https://ai-router.dev
- Chinese: https://ai-router.dev/cn
- Russian: https://ai-router.dev/ru
- Persian: https://ai-router.dev/fa

Implementation principle:

- Keep English at `/`.
- Use `/cn`, `/ru`, `/fa`, and other locale prefixes for localized landing pages.
- Link each language to its own localized path.
- Use canonical and hreflang tags on the website, not query parameters such as `?lang=`.
- Avoid thin translated pages. Each priority language should have useful localized content, not just translated navigation.

See the full supported language map in [languages.md](languages.md).

