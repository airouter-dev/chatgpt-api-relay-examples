import test from 'node:test';
import assert from 'node:assert/strict';
import { configFromEnv, normalizeBaseUrl } from '../src/config.mjs';

test('normalizes an OpenAI-compatible base URL exactly once', () => {
  assert.equal(normalizeBaseUrl('https://example.test/'), 'https://example.test/v1');
  assert.equal(normalizeBaseUrl('https://example.test/v1/'), 'https://example.test/v1');
});

test('accepts OPENAI_API_KEY as a convenient fallback', () => {
  const config = configFromEnv({ OPENAI_API_KEY: 'test-key' });
  assert.equal(config.apiKey, 'test-key');
  assert.equal(config.model, 'gpt-image-2');
});
