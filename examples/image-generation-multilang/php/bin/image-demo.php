#!/usr/bin/env php
<?php

declare(strict_types=1);

use AiRouter\ImageDemo\ArtifactStore;
use AiRouter\ImageDemo\Config;
use AiRouter\ImageDemo\ImageClient;
use AiRouter\ImageDemo\Workflow;

require dirname(__DIR__) . '/src/autoload.php';

const DEFAULT_MODEL = 'gpt-image-2';
const DEFAULT_SIZE = '1024x1024';

function usage(): void
{
    echo <<<'HELP'
AI-ROUTER image generation demo

Commands:
  generate --prompt TEXT [--model MODEL] [--size SIZE] [--output FILE]
  edit --image FILE --prompt TEXT [--model MODEL] [--size SIZE] [--output FILE]
  workflow --prompt TEXT --edit-prompt TEXT [--model MODEL] [--size SIZE]
           [--generated FILE] [--edited FILE]

Options accept both --key=value and --key value. Set AI_ROUTER_API_KEY before running.
HELP;
    echo PHP_EOL;
}

/** @return array{0: string|null, 1: array<string, string|bool>} */
function parseArguments(array $arguments): array
{
    $command = $arguments[0] ?? null;
    $start = 1;
    if (is_string($command) && str_starts_with($command, '--')) {
        $command = 'help';
        $start = 0;
    }
    $options = [];
    for ($index = $start; $index < count($arguments); $index++) {
        $token = $arguments[$index];
        if (!str_starts_with($token, '--')) {
            throw new InvalidArgumentException("Unexpected argument: {$token}");
        }
        $token = substr($token, 2);
        if (str_contains($token, '=')) {
            [$key, $value] = explode('=', $token, 2);
        } else {
            $key = $token;
            $next = $arguments[$index + 1] ?? null;
            if ($next !== null && !str_starts_with($next, '--')) {
                $value = $next;
                $index++;
            } else {
                $value = true;
            }
        }
        $options[$key] = $value;
    }
    return [$command, $options];
}

/** @param array<string, string|bool> $options */
function option(array $options, string $name, ?string $default = null, bool $required = false): string
{
    $value = $options[$name] ?? $default;
    if (!is_string($value) || trim($value) === '') {
        if ($required) {
            throw new InvalidArgumentException("Missing required option --{$name}");
        }
        return (string) $default;
    }
    return $value;
}

try {
    [$command, $options] = parseArguments(array_slice($argv, 1));
    if ($command === null || $command === 'help' || isset($options['help'])) {
        usage();
        exit(0);
    }

    $config = Config::fromEnvironment();
    if (isset($options['base-url']) && is_string($options['base-url'])) {
        $config = $config->withBaseUrl($options['base-url']);
    }
    $client = new ImageClient($config);
    $model = option($options, 'model', getenv('AI_ROUTER_MODEL') ?: DEFAULT_MODEL);
    $size = option($options, 'size', getenv('AI_ROUTER_SIZE') ?: DEFAULT_SIZE);

    if ($command === 'generate') {
        $image = $client->generate(option($options, 'prompt', required: true), $model, $size);
        $output = option($options, 'output', 'output/php-generated.png');
        ArtifactStore::save($output, $client->imageBytes($image));
        echo "Generated image saved to {$output}" . PHP_EOL;
    } elseif ($command === 'edit') {
        $image = $client->edit(
            option($options, 'image', required: true),
            option($options, 'prompt', required: true),
            $model,
            $size,
        );
        $output = option($options, 'output', 'output/php-edited.png');
        ArtifactStore::save($output, $client->imageBytes($image));
        echo "Edited image saved to {$output}" . PHP_EOL;
    } elseif ($command === 'workflow') {
        $result = (new Workflow($client))->run(
            option($options, 'prompt', required: true),
            option($options, 'edit-prompt', required: true),
            option($options, 'generated', 'output/php-generated.png'),
            option($options, 'edited', 'output/php-edited.png'),
            $model,
            $size,
        );
        echo "Generated image saved to {$result['generated']}" . PHP_EOL;
        echo "Edited image saved to {$result['edited']}" . PHP_EOL;
    } else {
        throw new InvalidArgumentException("Unknown command: {$command}");
    }
} catch (Throwable $error) {
    fwrite(STDERR, 'Error: ' . $error->getMessage() . PHP_EOL);
    exit(1);
}
