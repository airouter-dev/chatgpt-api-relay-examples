import fs from 'node:fs';
import { wrapFetchWithPaymentFromConfig } from '@x402/fetch';
import { ExactEvmScheme } from '@x402/evm';
import { privateKeyToAccount } from 'viem/accounts';

const NETWORK = 'eip155:8453';
const ASSET = '0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913'.toLowerCase();
const PAY_TO = '0xEe027fE58800Bcef81469EDF0C2E7f2ab54965De'.toLowerCase();

function challengeGuard(version, accepts) {
  const matching = accepts.filter((a) => a.scheme === 'exact' && a.network === NETWORK && String(a.asset).toLowerCase() === ASSET && String(a.payTo).toLowerCase() === PAY_TO);
  if (matching.length !== 1) throw new Error('402 challenge did not contain the expected exact Base USDC recipient');
  const amount = BigInt(matching[0].amount);
  if (amount < 1n || amount > 18000n) throw new Error(`402 amount ${amount} is outside the configured 1-4 image range`);
  return matching[0];
}

function buildRequest(config) {
  const url = `${config.baseUrl}/x402/v1/images/${config.endpoint}`;
  if (config.endpoint === 'generations') {
    return { url, init: { method: 'POST', headers: { 'content-type': 'application/json' }, body: JSON.stringify({ model: config.model, prompt: config.prompt, n: config.n }) } };
  }
  const form = new FormData();
  form.set('model', config.model); form.set('prompt', config.prompt); form.set('n', String(config.n));
  form.set('image', new Blob([fs.readFileSync(config.imagePath)], { type: 'image/png' }), config.imagePath.split(/[\\/]/).pop());
  return { url, init: { method: 'POST', body: form } };
}

export async function run(config) {
  const { url, init } = buildRequest(config);
  if (config.mode === 'probe') {
    const response = await fetch(url, init);
    if (response.status !== 402) throw new Error(`probe returned HTTP ${response.status}: ${(await response.text()).slice(0, 300)}`);
    const encoded = response.headers.get('payment-required');
    if (!encoded) throw new Error('402 response has no Payment-Required header');
    const payment = JSON.parse(Buffer.from(encoded, 'base64url').toString('utf8'));
    const challenge = challengeGuard(payment.x402Version, payment.accepts || []);
    console.log(JSON.stringify({ mode: 'probe', status: response.status, endpoint: config.endpoint, model: config.model, amount: challenge.amount, network: challenge.network, asset: challenge.asset, payTo: challenge.payTo }));
    return payment;
  }
  const account = privateKeyToAccount(config.privateKey);
  const paidFetch = wrapFetchWithPaymentFromConfig(fetch, { schemes: [{ network: NETWORK, client: new ExactEvmScheme(account) }], paymentRequirementsSelector: (_version, accepts) => challengeGuard(_version, accepts) });
  const response = await paidFetch(url, init);
  if (!response.ok) throw new Error(`paid request returned HTTP ${response.status}: ${(await response.text()).slice(0, 500)}`);
  const data = await response.json();
  console.log(JSON.stringify({ mode: 'pay', status: response.status, endpoint: config.endpoint, model: config.model, data }));
  return data;
}
