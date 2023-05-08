import { Static, Type } from '@sinclair/typebox';

const post = Type.Object({
  _id: Type.String(),
  postId: Type.String(),
  title: Type.String(),
  content: Type.String(),
  images: Type.Array(Type.String()),
  author: Type.String(),
  likes: Type.String(),
  createdAt: Type.String(),
  updatedAt: Type.String(),
});

export const createPostSchema = {
  tags: ['Post'],
  description: 'Create a blog post',
  body: Type.Object({
    title: Type.String(),
    content: Type.String(),
    author: Type.String(),
    images: Type.Array(Type.String()),
  }),
  response: {
    201: post,
  },
};

export type CreatePostSchema = Static<typeof createPostSchema.body>;
