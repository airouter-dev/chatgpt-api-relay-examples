import { loadConfig } from './config.mjs';
import { run } from './client.mjs';

const parsedArgs = process.argv.slice(2).reduce((out, item, index, values) => {
  if (!item.startsWith('--')) return out;
  out.push([item.slice(2), values[index + 1] && !values[index + 1].startsWith('--') ? values[index + 1] : true]);
  return out;
}, []);
const args = Object.fromEntries(parsedArgs);
await run(loadConfig({ mode: args.mode || undefined, endpoint: args.endpoint || undefined, model: args.model || undefined, imagePath: args.image || undefined }));
