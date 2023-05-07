import mongoose from 'mongoose';
import { logger } from '../log/logger';
import { config } from './config';

export async function connect() {
  try {
    await mongoose.connect(config.DATABASE_URL);
    logger.info('Connected to database');
  } catch (error) {
    logger.error('Error connecting to database', error);
    process.exit(1);
  }
}

export function disconnect() {
  return mongoose.connection.close();
}
