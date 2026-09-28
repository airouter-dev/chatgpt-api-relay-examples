/** Provider failures with safe, actionable diagnostics. */

export class ImageDemoError extends Error {}

export class ImageApiError extends ImageDemoError {
  constructor(message, { status = undefined, body = '' } = {}) {
    super(`${message}${status ? ` (HTTP ${status})` : ''}`);
    this.name = 'ImageApiError';
    this.status = status;
    this.body = body;
  }
}
