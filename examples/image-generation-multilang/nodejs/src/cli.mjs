#!/usr/bin/env node

import { resolve } from 'node:path';
import { ImagesClient } from './client.mjs';
import { configFromEnv, withOverrides } from './config.mjs';
import { ImageWorkflow } from './workflow.mjs';

function usage() {
  console.log(`Usage: node src/cli.mjs --prompt "..." --edit-prompt "..." [options]

Options:
  --output-dir DIR   Where generated.png and edited.png are saved (default: outputs)
  --model ID         Model id (default: IMAGE_MODEL or gpt-image-2)
  --size VALUE       Provider-supported image size
  --quality VALUE    Provider-supported quality tier
  --base-url URL     OpenAI-compatible API base URL
  --api-key KEY      API key (prefer IMAGE_API_KEY in the environment)
  --timeout MS       Request timeout in milliseconds
  --help             Show this help
`);
}

function parseArgs(argv) {
  const values = {};
  const takesValue = new Set(['prompt', 'edit-prompt', 'output-dir', 'model', 'size', 'quality', 'base-url', 'api-key', 'timeout']);
  for (let index = 0; index < argv.length; index += 1) {
    const token = argv[index];
    if (token === '--help') return { help: true };
    if (!token.startsWith('--') || !takesValue.has(token.slice(2))) throw new TypeError(`unknown option: ${token}`);
    const key = token.slice(2).replace(/-([a-z])/g, (_, letter) => letter.toUpperCase());
    const value = argv[++index];
    if (!value || value.startsWith('--')) throw new TypeError(`${token} requires a value`);
    values[key] = value;
  }
  if (!values.prompt || !values.editPrompt) throw new TypeError('--prompt and --edit-prompt are required');
  return values;
}

export async function main(argv = process.argv.slice(2)) {
  const args = parseArgs(argv);
  if (args.help) {
    usage();
    return 0;
  }
  const config = withOverrides(configFromEnv(), {
    apiKey: args.apiKey,
    baseUrl: args.baseUrl,
    model: args.model,
    size: args.size,
    quality: args.quality,
    timeoutMs: args.timeout,
  });
  const outputDir = resolve(args.outputDir || 'outputs');
  const result = await new ImageWorkflow(new ImagesClient(config)).run(args.prompt, args.editPrompt, outputDir);
  console.log(`Generated image: ${result.generatedPath}`);
  console.log(`Edited image:    ${result.editedPath}`);
  if (result.generated.revisedPrompt) console.log(`Provider revised prompt: ${result.generated.revisedPrompt}`);
  return 0;
}

if (import.meta.url === `file://${process.argv[1]}`) {
  main().then((code) => process.exitCode = code).catch((error) => {
    console.error(`image demo failed: ${error.message}`);
    process.exitCode = 1;
  });
}
