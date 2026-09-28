<?php

declare(strict_types=1);

namespace AiRouter\ImageDemo;

use InvalidArgumentException;

final class ImageData
{
    public function __construct(
        public ?string $base64,
        public ?string $url,
        public ?string $revisedPrompt = null,
    ) {
        if ($this->base64 === null && $this->url === null) {
            throw new InvalidArgumentException('Image response contains neither b64_json nor url.');
        }
    }

    /** @param array<string, mixed> $data */
    public static function fromArray(array $data): self
    {
        $base64 = isset($data['b64_json']) && is_string($data['b64_json'])
            ? $data['b64_json']
            : null;
        $url = isset($data['url']) && is_string($data['url']) ? $data['url'] : null;
        $revisedPrompt = isset($data['revised_prompt']) && is_string($data['revised_prompt'])
            ? $data['revised_prompt']
            : null;
        return new self($base64, $url, $revisedPrompt);
    }
}
