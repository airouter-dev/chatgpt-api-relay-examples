/**
 * Small transport adapter for OpenAI-compatible Images endpoints.
 * Node 18's fetch/FormData/Blob are enough; no SDK or runtime dependency is
 * needed. The workflow layer never needs to know how multipart is encoded.
 */

import { readFile } from 'node:fs/promises';
import { basename, extname } from 'node:path';
import { ImageApiError } from './errors.mjs';

const MIME_BY_EXTENSION = {
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.webp': 'image/webp',
};

function mimeTypeFor(path) {
  return MIME_BY_EXTENSION[extname(path).toLowerCase()] || 'application/octet-stream';
}

async function responseBody(response) {
  const text = await response.text();
  try {
    return JSON.parse(text);
  } catch {
    return text.slice(0, 1000);
  }
}

function errorMessage(payload) {
  if (payload && typeof payload === 'object' && payload.error?.message) return String(payload.error.message);
  return typeof payload === 'string' ? payload : 'provider returned an error';
}

async function fetchWithTimeout(url, options, timeoutMs) {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), timeoutMs);
  try {
    return await fetch(url, { ...options, signal: controller.signal });
  } catch (error) {
    if (error?.name === 'AbortError') throw new ImageApiError(`request timed out after ${timeoutMs} ms`);
    throw new ImageApiError(`could not reach image API: ${error.message}`);
  } finally {
    clearTimeout(timer);
  }
}

async function decodeImageResponse(payload, timeoutMs) {
  if (!payload || !Array.isArray(payload.data) || !payload.data[0]) {
    throw new ImageApiError('Images API returned no image data');
  }
  const item = payload.data[0];
  const revisedPrompt = typeof item.revised_prompt === 'string' ? item.revised_prompt : undefined;
  if (typeof item.b64_json === 'string' && item.b64_json) {
    return {
      bytes: Buffer.from(item.b64_json, 'base64'),
      mediaType: 'image/png',
      revisedPrompt,
    };
  }
  if (typeof item.url !== 'string' || !item.url) throw new ImageApiError('Images API item contained neither b64_json nor url');
  if (item.url.startsWith('data:')) {
    const match = item.url.match(/^data:([^;,]+)?;base64,(.+)$/s);
    if (!match) throw new ImageApiError('Images API returned an invalid data URL');
    return { bytes: Buffer.from(match[2], 'base64'), mediaType: match[1] || 'image/png', sourceUrl: item.url, revisedPrompt };
  }
  const response = await fetchWithTimeout(item.url, { headers: { 'User-Agent': 'ai-router-image-demo/0.1' } }, timeoutMs);
  if (!response.ok) throw new ImageApiError('could not download image URL', { status: response.status, body: await response.text() });
  const contentType = response.headers.get('content-type')?.split(';', 1)[0] || 'image/png';
  return { bytes: Buffer.from(await response.arrayBuffer()), mediaType: contentType, sourceUrl: item.url, revisedPrompt };
}

export class ImagesClient {
  constructor(config) {
    this.config = config;
  }

  async #request(path, options) {
    const response = await fetchWithTimeout(`${this.config.baseUrl}${path}`, {
      ...options,
      headers: {
        Authorization: `Bearer ${this.config.apiKey}`,
        Accept: 'application/json',
        ...(options.headers || {}),
      },
    }, this.config.timeoutMs);
    const payload = await responseBody(response);
    if (!response.ok) throw new ImageApiError(`Images API request failed: ${errorMessage(payload)}`, { status: response.status, body: JSON.stringify(payload) });
    if (typeof payload !== 'object' || payload === null) throw new ImageApiError('Images API returned invalid JSON');
    return payload;
  }

  async generate(prompt, { model, size, quality } = {}) {
    if (!String(prompt).trim()) throw new TypeError('prompt cannot be empty');
    const payload = await this.#request('/images/generations', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        model: model || this.config.model,
        prompt,
        size: size || this.config.size,
        quality: quality || this.config.quality,
        response_format: 'b64_json',
      }),
    });
    return decodeImageResponse(payload, this.config.timeoutMs);
  }

  async edit(imagePath, prompt, { model, size, quality } = {}) {
    if (!String(prompt).trim()) throw new TypeError('edit prompt cannot be empty');
    const bytes = await readFile(imagePath);
    const form = new FormData();
    form.set('model', model || this.config.model);
    form.set('prompt', prompt);
    form.set('size', size || this.config.size);
    form.set('quality', quality || this.config.quality);
    form.set('response_format', 'b64_json');
    form.set('image', new Blob([bytes], { type: mimeTypeFor(imagePath) }), basename(imagePath));
    const payload = await this.#request('/images/edits', { method: 'POST', body: form });
    return decodeImageResponse(payload, this.config.timeoutMs);
  }
}

export function extensionForMediaType(mediaType) {
  const subtype = String(mediaType || 'image/png').toLowerCase().split('/', 2)[1];
  return subtype === 'jpeg' || subtype === 'jpg' ? '.jpg' : subtype === 'webp' ? '.webp' : '.png';
}
