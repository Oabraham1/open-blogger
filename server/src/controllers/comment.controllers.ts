import { FastifyReply, FastifyRequest } from 'fastify';
import { CreateCommentSchema } from '../schema/comment.schema';
import { logger } from '../../log/logger';
import { createComment } from '../services/comment.service';

export async function createCommentHandler(
  request: FastifyRequest<{ Body: CreateCommentSchema }>,
  response: FastifyReply
) {
  try {
    const comment = await createComment(request.body);

    return response.status(201).send(comment);
  } catch (error) {
    logger.error(error, 'createCommentHandler: Error adding comment');
    return response.status(400).send({ message: 'Error adding comment' });
  }
}
