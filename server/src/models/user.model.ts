import { prop } from '@typegoose/typegoose';

export class User {
  @prop({
    type: String,
    required: true,
  })
  userId: string;

  @prop({
    type: String,
    required: true,
  })
  username: string;

  @prop({
    type: String,
    required: true,
  })
  email: string;

  @prop({
    type: String,
    required: true,
  })
  avatar: string;

  @prop({
    type: String,
    required: true,
  })
  password: string;

  @prop({
    type: String,
    default: new Date().toUTCString(),
  })
  createdAt: string;
}
