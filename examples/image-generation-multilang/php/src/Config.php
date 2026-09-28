<?php

declare(strict_types=1);

namespace AiRouter\ImageDemo;

use RuntimeException;

final class Config
{
    public function __construct(
        public string $baseUrl,
        public string $apiKey,
        public int $timeoutSeconds = 180,
    ) {
    }

    public static function fromEnvironment(): self
    {
        $apiKey = trim((string) (getenv('AI_ROUTER_API_KEY') ?: getenv('OPENAI_API_KEY')));
        if ($apiKey === '') {
            throw new RuntimeException(
                'AI_ROUTER_API_KEY (or OPENAI_API_KEY) must be set before making a request.'
            );
        }

        $baseUrl = trim((string) (getenv('AI_ROUTER_BASE_URL') ?: 'https://api.ai-router.dev/v1'));
        $baseUrl = rtrim($baseUrl, '/');
        if ($baseUrl === '') {
            throw new RuntimeException('AI_ROUTER_BASE_URL must not be empty.');
        }

        $timeout = (int) (getenv('AI_ROUTER_TIMEOUT_SECONDS') ?: 180);
        return new self($baseUrl, $apiKey, max(1, $timeout));
    }

    public function withBaseUrl(string $baseUrl): self
    {
        return new self(rtrim($baseUrl, '/'), $this->apiKey, $this->timeoutSeconds);
    }
}
