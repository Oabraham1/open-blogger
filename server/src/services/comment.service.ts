import { nanoid } from 'nanoid';
import { CommentModel } from '../models/comment.model';
import { CreateCommentSchema } from '../schema/comment.schema';

export async function createComment(input: CreateCommentSchema) {
  const commentId = `comment_${nanoid()}`;

  return CommentModel.create({
    commentId,
    ...input,
  });
}
