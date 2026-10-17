import 'dotenv/config';

export function loadConfig(overrides = {}) {
  const config = {
    baseUrl: process.env.AI_ROUTER_X402_BASE_URL || 'https://api.ai-router.dev',
    mode: process.env.AI_ROUTER_X402_MODE || 'probe',
    endpoint: process.env.AI_ROUTER_X402_ENDPOINT || 'generations',
    model: process.env.AI_ROUTER_X402_MODEL || 'gpt-image-2',
    prompt: process.env.AI_ROUTER_X402_PROMPT || 'A small red fox under a moonlit sky',
    n: Number(process.env.AI_ROUTER_X402_N || 1),
    imagePath: process.env.AI_ROUTER_X402_IMAGE_PATH || '',
    privateKey: process.env.AI_ROUTER_X402_EVM_PRIVATE_KEY || '',
  };
  for (const [key, value] of Object.entries(overrides)) {
    if (value !== undefined && value !== null) config[key] = value;
  }
  config.baseUrl = String(config.baseUrl).replace(/\/+$/, '');
  if (!['probe', 'pay'].includes(config.mode)) throw new TypeError('mode must be probe or pay');
  if (!['generations', 'edits'].includes(config.endpoint)) throw new TypeError('endpoint must be generations or edits');
  if (!config.model || !config.prompt || !Number.isInteger(config.n) || config.n < 1 || config.n > 4) throw new TypeError('model, prompt, and n (1-4) are required');
  if (config.mode === 'pay' && !/^0x[0-9a-fA-F]{64}$/.test(config.privateKey)) throw new TypeError('pay mode requires AI_ROUTER_X402_EVM_PRIVATE_KEY');
  if (config.endpoint === 'edits' && config.mode === 'pay' && !config.imagePath) throw new TypeError('edits pay mode requires AI_ROUTER_X402_IMAGE_PATH');
  return config;
}
