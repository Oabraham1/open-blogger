import { connect, disconnect } from './config/mongo';
import { logger } from './log/logger';
import { createServer } from './utils/create';
import { config } from './config/config';

const signals = ['SIGINT', 'SIGTERM', 'SIGHUP'] as const;

async function stopServer({
  signal,
  server,
}: {
  signal: (typeof signals)[number];
  server: Awaited<ReturnType<typeof createServer>>;
}) {
  logger.info(`Server is stopping due to ${signal}`);
  await server.close();
  await disconnect();
  logger.info('Server is stopped');
  process.exit(0);
}

async function startServer() {
  const server = await createServer();
  server.listen({
    port: config.PORT,
    host: config.HOST,
  });

  await connect();
  logger.info('Server is running');

  signals.forEach((signal) => {
    process.on(signal, () => stopServer({ signal, server }));
  });
}

startServer();
