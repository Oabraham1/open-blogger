import pino from 'pino';
export const logger = pino({
  redact: ['hostname'],
  timestamp() {
    return `,"time":"${new Date().toDateString()}"`;
  },
});
