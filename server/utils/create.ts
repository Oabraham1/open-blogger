import fastify from 'fastify';
import { postRoutes } from '../src/routes/post.route';

export async function createServer() {
  const app = fastify();

  app.register(postRoutes, { prefix: '/api/post' });

  return app;
}
