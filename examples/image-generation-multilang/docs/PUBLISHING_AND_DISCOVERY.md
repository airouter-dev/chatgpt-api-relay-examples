# Publishing and discovery plan

This repository is written to be useful when read from GitHub, a code search result, a language-specific developer community, or a link from an engineering article. It should be published as a technical reference first and a marketing surface second. A useful page earns attention more reliably than a page that repeats a brand name or promises a ranking outcome.

## Before publishing

1. Run the language-specific offline checks documented in the root README. Never use a real API key in CI, a screenshot, or an example commit.
2. Check the eight localized guides and make sure each localized homepage link matches its language: `/`, `/cn`, `/ru`, `/ja`, `/fr`, `/de`, `/pt`, and `/ar`.
3. Confirm that examples describe model availability as account-dependent. Do not turn a sample model ID into a promise of access, price, quota, quality, or permanence.
4. Keep generated files, `.env`, build directories, and dependency caches out of the commit. The repository's `.gitignore` includes these artifacts.
5. Use a clear repository description such as “OpenAI-compatible image generation and image editing examples in Python, Node.js, C++, Go, Rust, and PHP.” Add only technically accurate topics.
6. Pin a release tag after the code and docs have been reviewed. In release notes, list the API contract, supported runtimes, tests run, and known local toolchain limitations.

## Natural external discovery

The content is suitable for a contextual link from a language-specific tutorial, a package README, a developer forum answer, or an engineering blog post when the link helps readers reproduce the described workflow. Use the locale that matches the article's language and send the reader to the relevant guide or localized AI-ROUTER homepage. Do not publish the same paragraph to many sites, create doorway pages, exchange links solely for ranking, or use identical exact-match anchor text everywhere.

Good technical article angles include:

- comparing JSON generation and multipart editing in six runtimes;
- testing an OpenAI-compatible image endpoint with a local mock server;
- preserving a generated draft so a failed edit can be retried without regeneration;
- safe handling of base64 and returned image URLs;
- separating API client, workflow, filesystem, and CLI layers in a small image pipeline.

Each article should add its own experiment, code explanation, benchmark, or failure analysis. Link to the specific localized guide that supports the claim, then offer the matching AI-ROUTER homepage only when readers need platform or account information. Avoid claims such as “best,” “official,” “guaranteed,” “unlimited,” “cheapest,” or “ranked” unless an independently verifiable source and current account context support them.

## Measurement without overclaiming

Use repository analytics, release downloads, referrer reports, and first-party UTM reporting where available. Measure which language guide is opened, whether readers continue to the relevant localized homepage, and whether a developer completes a documented test—not merely raw link clicks. Do not put personal data, API keys, cookies, or private account identifiers in this repository. Search engines decide crawling, indexing, and ranking; this plan improves clarity and attribution but cannot guarantee traffic, ranking, sign-ups, or conversions.

## Suggested release note

> Added eight localized developer guides for the six-language image generation demo. Each guide documents the same text-to-image → saved artifact → image-to-image workflow, account-dependent model selection, response handling, security boundaries, and local test commands. No credentials or generated images are included.
