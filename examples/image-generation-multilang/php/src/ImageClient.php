<?php

declare(strict_types=1);

namespace AiRouter\ImageDemo;

use CURLFile;
use JsonException;
use RuntimeException;

final class ImageClient
{
    public function __construct(private readonly Config $config)
    {
    }

    public function generate(string $prompt, string $model = 'gpt-image-2', string $size = '1024x1024'): ImageData
    {
        $this->assertPrompt($prompt);
        $payload = $this->request('POST', '/images/generations', [
            'model' => $model,
            'prompt' => $prompt,
            'size' => $size,
            'n' => 1,
            'response_format' => 'b64_json',
        ]);
        return $this->firstImage($payload, 'generation');
    }

    public function edit(
        string $imagePath,
        string $prompt,
        string $model = 'gpt-image-2',
        string $size = '1024x1024',
    ): ImageData {
        $this->assertPrompt($prompt);
        if (!is_file($imagePath) || !is_readable($imagePath)) {
            throw new RuntimeException("Input image is not readable: {$imagePath}");
        }
        if ((int) filesize($imagePath) === 0) {
            throw new RuntimeException("Input image is empty: {$imagePath}");
        }

        $mime = function_exists('mime_content_type')
            ? (mime_content_type($imagePath) ?: 'image/png')
            : 'image/png';
        $payload = $this->request('POST', '/images/edits', [
            'model' => $model,
            'prompt' => $prompt,
            'size' => $size,
            'response_format' => 'b64_json',
            'image' => new CURLFile($imagePath, $mime, basename($imagePath)),
        ], multipart: true);
        return $this->firstImage($payload, 'edit');
    }

    public function imageBytes(ImageData $image): string
    {
        if ($image->base64 !== null) {
            $bytes = base64_decode($image->base64, true);
            if ($bytes === false) {
                throw new RuntimeException('Provider returned invalid base64 image data.');
            }
            if ($bytes === '') {
                throw new RuntimeException('Provider returned an empty image.');
            }
            return $bytes;
        }

        // URL responses are fetched separately so API response parsing remains pure.
        return $this->download($image->url ?? '');
    }

    /** @return array<string, mixed> */
    private function request(string $method, string $path, array $payload, bool $multipart = false): array
    {
        $handle = curl_init($this->config->baseUrl . '/' . ltrim($path, '/'));
        if ($handle === false) {
            throw new RuntimeException('Could not initialise cURL.');
        }

        $headers = [
            'Accept: application/json',
            'Authorization: Bearer ' . $this->config->apiKey,
        ];
        $options = [
            CURLOPT_CUSTOMREQUEST => $method,
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_FOLLOWLOCATION => true,
            CURLOPT_CONNECTTIMEOUT => 20,
            CURLOPT_TIMEOUT => $this->config->timeoutSeconds,
            CURLOPT_HTTPHEADER => $headers,
        ];
        if ($multipart) {
            $options[CURLOPT_POSTFIELDS] = $payload;
        } else {
            $headers[] = 'Content-Type: application/json';
            $options[CURLOPT_HTTPHEADER] = $headers;
            try {
                $options[CURLOPT_POSTFIELDS] = json_encode($payload, JSON_THROW_ON_ERROR);
            } catch (JsonException $error) {
                curl_close($handle);
                throw new RuntimeException('Could not encode request JSON.', 0, $error);
            }
        }

        curl_setopt_array($handle, $options);
        $body = curl_exec($handle);
        $curlError = curl_error($handle);
        $status = (int) curl_getinfo($handle, CURLINFO_HTTP_CODE);
        curl_close($handle);
        if ($body === false) {
            throw new RuntimeException('HTTP request failed: ' . ($curlError ?: 'unknown cURL error'));
        }
        if ($status < 200 || $status >= 300) {
            throw new ApiException($this->apiErrorMessage((string) $body), $status);
        }

        try {
            $decoded = json_decode((string) $body, true, 512, JSON_THROW_ON_ERROR);
        } catch (JsonException $error) {
            throw new RuntimeException('Provider returned invalid JSON.', 0, $error);
        }
        if (!is_array($decoded)) {
            throw new RuntimeException('Provider returned an unexpected JSON shape.');
        }
        return $decoded;
    }

    private function download(string $url): string
    {
        if ($url === '') {
            throw new RuntimeException('Provider returned an empty image URL.');
        }
        $handle = curl_init($url);
        if ($handle === false) {
            throw new RuntimeException('Could not initialise cURL for image download.');
        }
        curl_setopt_array($handle, [
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_FOLLOWLOCATION => true,
            CURLOPT_CONNECTTIMEOUT => 20,
            CURLOPT_TIMEOUT => $this->config->timeoutSeconds,
        ]);
        $body = curl_exec($handle);
        $error = curl_error($handle);
        $status = (int) curl_getinfo($handle, CURLINFO_HTTP_CODE);
        curl_close($handle);
        if ($body === false || $status < 200 || $status >= 300) {
            throw new RuntimeException('Image download failed: ' . ($error ?: "HTTP {$status}"));
        }
        return (string) $body;
    }

    /** @param array<string, mixed> $payload */
    private function firstImage(array $payload, string $operation): ImageData
    {
        $data = $payload['data'] ?? null;
        if (!is_array($data) || !isset($data[0]) || !is_array($data[0])) {
            throw new RuntimeException("Provider returned no image for {$operation}.");
        }
        try {
            return ImageData::fromArray($data[0]);
        } catch (\InvalidArgumentException $error) {
            throw new RuntimeException("Provider returned an unusable {$operation} image.", 0, $error);
        }
    }

    private function assertPrompt(string $prompt): void
    {
        if (trim($prompt) === '') {
            throw new RuntimeException('Prompt must not be empty.');
        }
    }

    private function apiErrorMessage(string $body): string
    {
        try {
            $decoded = json_decode($body, true, 512, JSON_THROW_ON_ERROR);
            $message = $decoded['error']['message'] ?? null;
            if (is_string($message) && $message !== '') {
                return $message;
            }
        } catch (JsonException) {
            // Fall through to the provider's raw response.
        }
        return trim($body) !== '' ? trim($body) : 'Provider returned an HTTP error.';
    }
}
