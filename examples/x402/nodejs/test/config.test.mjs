import test from 'node:test';
import assert from 'node:assert/strict';
import { loadConfig } from '../src/config.mjs';

test('probe configuration has safe defaults', () => {
  const cfg = loadConfig({ mode: 'probe' });
  assert.equal(cfg.endpoint, 'generations');
  assert.equal(cfg.model, 'gpt-image-2');
});

test('pay mode rejects an absent private key', () => {
  assert.throws(() => loadConfig({ mode: 'pay', privateKey: '' }), /PRIVATE_KEY/);
});
