import { nanoid } from 'nanoid';
import { PostModel } from '../models/post.model';
import { CreatePostSchema } from '../schema/post.schema';

export async function createPost(input: CreatePostSchema) {
  const postId = `post_${nanoid()}`;

  return PostModel.create({
    postId,
    ...input,
  });
}
