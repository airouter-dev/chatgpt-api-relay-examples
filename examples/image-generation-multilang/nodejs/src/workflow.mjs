/** Generate, persist, edit, and persist: the reusable application workflow. */

import { mkdir, writeFile } from 'node:fs/promises';
import { join } from 'node:path';
import { extensionForMediaType } from './client.mjs';

async function saveAsset(asset, outputDir, stem) {
  await mkdir(outputDir, { recursive: true });
  const target = join(outputDir, `${stem}${extensionForMediaType(asset.mediaType)}`);
  await writeFile(target, asset.bytes);
  return target;
}

export class ImageWorkflow {
  constructor(client) {
    this.client = client;
  }

  async run(prompt, editPrompt, outputDir) {
    const generated = await this.client.generate(prompt);
    const generatedPath = await saveAsset(generated, outputDir, 'generated');
    const edited = await this.client.edit(generatedPath, editPrompt);
    const editedPath = await saveAsset(edited, outputDir, 'edited');
    return { generated, generatedPath, edited, editedPath };
  }
}
