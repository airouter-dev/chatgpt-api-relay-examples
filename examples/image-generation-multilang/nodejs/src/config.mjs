/** Runtime settings shared by the generation and edit requests. */

export function normalizeBaseUrl(value) {
  const candidate = String(value ?? '').trim().replace(/\/+$/, '');
  if (!candidate) throw new TypeError('IMAGE_API_BASE_URL cannot be empty');
  let parsed;
  try {
    parsed = new URL(candidate);
  } catch {
    throw new TypeError('IMAGE_API_BASE_URL must be an http(s) URL');
  }
  if (!['http:', 'https:'].includes(parsed.protocol) || !parsed.host) {
    throw new TypeError('IMAGE_API_BASE_URL must be an http(s) URL');
  }
  return parsed.pathname.replace(/\/+$/, '').endsWith('/v1') ? candidate : `${candidate}/v1`;
}

export function configFromEnv(env = process.env) {
  const apiKey = env.IMAGE_API_KEY || env.OPENAI_API_KEY || env.AI_ROUTER_API_KEY;
  if (!apiKey) throw new TypeError('Set AI_ROUTER_API_KEY (or IMAGE_API_KEY/OPENAI_API_KEY) before running the demo');
  const timeoutMs = Number(env.IMAGE_TIMEOUT_MS ?? 120000);
  if (!Number.isFinite(timeoutMs) || timeoutMs <= 0) throw new TypeError('IMAGE_TIMEOUT_MS must be greater than zero');
  return {
    apiKey,
    baseUrl: normalizeBaseUrl(env.IMAGE_API_BASE_URL || env.AI_ROUTER_BASE_URL || 'https://api.ai-router.dev/v1'),
    model: env.IMAGE_MODEL || env.AI_ROUTER_MODEL || 'gpt-image-2',
    size: env.IMAGE_SIZE || env.AI_ROUTER_SIZE || '1024x1024',
    quality: env.IMAGE_QUALITY || 'auto',
    timeoutMs,
  };
}

export function withOverrides(config, overrides = {}) {
  const next = { ...config };
  for (const [key, value] of Object.entries(overrides)) {
    if (value !== undefined && value !== null) next[key] = value;
  }
  if (next.baseUrl) next.baseUrl = normalizeBaseUrl(next.baseUrl);
  if (!Number.isFinite(Number(next.timeoutMs)) || Number(next.timeoutMs) <= 0) {
    throw new TypeError('timeout must be greater than zero');
  }
  next.timeoutMs = Number(next.timeoutMs);
  return next;
}
