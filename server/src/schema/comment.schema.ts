import { Static, Type } from '@sinclair/typebox';

const comment = Type.Object({
  _id: Type.String(),
  commentId: Type.String(),
  postId: Type.String(),
  content: Type.String(),
  author: Type.String(),
  authorId: Type.String(),
  createdAt: Type.String(),
  updatedAt: Type.String(),
});

export const createCommentSchema = {
  tags: ['Comment'],
  description: 'Comment on a blog post',
  body: Type.Object({
    postId: Type.String(),
    content: Type.String(),
    author: Type.String(),
    authorId: Type.String(),
  }),
  response: {
    201: comment,
  },
};

export type CreateCommentSchema = Static<typeof createCommentSchema.body>;
