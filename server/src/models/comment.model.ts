import { getModelForClass, prop } from '@typegoose/typegoose';

export class Comment {
  @prop({
    type: String,
    required: true,
  })
  commentId: string;

  @prop({
    type: String,
    required: true,
  })
  postId: string;

  @prop({
    type: String,
    required: true,
  })
  content: string;

  @prop({
    type: String,
    required: true,
  })
  author: string;

  @prop({
    type: String,
    required: true,
  })
  authorId: string;

  @prop({
    type: String,
    default: new Date().toUTCString(),
  })
  createdAt: string;

  @prop({
    type: String,
    default: new Date().toUTCString(),
  })
  updatedAt: string;
}

export const CommentModel = getModelForClass(Comment, {
  schemaOptions: { timestamps: true },
});
