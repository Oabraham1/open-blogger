import { FastifyInstance, FastifyPluginOptions } from 'fastify';
import { createPostHandler } from '../controllers/post.controller';
import { createPostSchema } from '../schema/post.schema';
import { createCommentSchema } from '../schema/comment.schema';
import { createCommentHandler } from '../controllers/comment.controllers';

export function postRoutes(
  app: FastifyInstance,
  options: FastifyPluginOptions,
  done: () => void
) {
  app.post('/', { schema: createPostSchema }, createPostHandler);

  app.post('/comment', { schema: createCommentSchema }, createCommentHandler);

  done();
}
