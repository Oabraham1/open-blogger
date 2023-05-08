import { getModelForClass, prop } from '@typegoose/typegoose';

export class Post {
  @prop({
    type: String,
    required: true,
  })
  postId: string;

  @prop({
    type: String,
    required: true,
  })
  title: string;

  @prop({
    type: String,
    required: true,
  })
  content: string;

  @prop({
    type: Array<String>,
    required: false,
  })
  images: string[];

  @prop({
    type: String,
    required: true,
  })
  author: string;

  @prop({
    type: String,
    default: '0',
  })
  likes: string;

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

export const PostModel = getModelForClass(Post, {
  schemaOptions: { timestamps: true },
});
