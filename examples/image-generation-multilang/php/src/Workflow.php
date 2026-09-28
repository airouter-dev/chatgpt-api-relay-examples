<?php

declare(strict_types=1);

namespace AiRouter\ImageDemo;

final class Workflow
{
    public function __construct(private readonly ImageClient $client)
    {
    }

    /**
     * Generate an image, persist it, then pass that artifact to /images/edits.
     * @return array{generated: string, edited: string}
     */
    public function run(
        string $prompt,
        string $editPrompt,
        string $generatedPath,
        string $editedPath,
        string $model,
        string $size,
    ): array {
        $generated = $this->client->generate($prompt, $model, $size);
        ArtifactStore::save($generatedPath, $this->client->imageBytes($generated));

        $edited = $this->client->edit($generatedPath, $editPrompt, $model, $size);
        ArtifactStore::save($editedPath, $this->client->imageBytes($edited));
        return ['generated' => $generatedPath, 'edited' => $editedPath];
    }
}
