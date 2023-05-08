import { FastifyReply, FastifyRequest } from 'fastify';
import { CreatePostSchema } from '../schema/post.schema';
import { logger } from '../../log/logger';
import { createPost } from '../services/post.service';

export async function createPostHandler(
  request: FastifyRequest<{ Body: CreatePostSchema }>,
  response: FastifyReply
) {
  try {
    const post = await createPost(request.body);

    return response.status(201).send(post);
  } catch (error) {
    logger.error(error, 'createPostHandler: Error creating post');
    return response.status(400).send({ message: 'Error creating post' });
  }
}
