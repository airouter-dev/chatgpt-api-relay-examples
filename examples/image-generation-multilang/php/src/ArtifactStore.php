<?php

declare(strict_types=1);

namespace AiRouter\ImageDemo;

use RuntimeException;

final class ArtifactStore
{
    public static function save(string $path, string $bytes): void
    {
        $parent = dirname($path);
        if ($parent !== '.' && !is_dir($parent) && !mkdir($parent, 0775, true) && !is_dir($parent)) {
            throw new RuntimeException("Could not create output directory: {$parent}");
        }
        if (file_put_contents($path, $bytes) === false) {
            throw new RuntimeException("Could not write image output: {$path}");
        }
    }
}
