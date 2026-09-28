<?php

declare(strict_types=1);

namespace AiRouter\ImageDemo;

use RuntimeException;

final class ApiException extends RuntimeException
{
    public function __construct(
        string $message,
        public readonly int $statusCode = 0,
        ?\Throwable $previous = null,
    ) {
        parent::__construct($message, $statusCode, $previous);
    }
}
