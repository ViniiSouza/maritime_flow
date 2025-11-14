import { config } from 'dotenv';
import pkg from 'pg';

config();

const { Pool } = pkg;

const connectionOptions = {};

if (process.env.DATABASE_URL) {
  connectionOptions.connectionString = process.env.DATABASE_URL;
}

if (process.env.DB_SSL === 'true') {
  connectionOptions.ssl = { rejectUnauthorized: false };
}

const pool = new Pool(connectionOptions);

pool.on('error', (err) => {
  console.error('Unexpected PG error', err);
  process.exit(1);
});

export default pool;
