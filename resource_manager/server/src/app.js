import express from 'express';
import cors from 'cors';
import towersRouter from './routes/towers.js';
import vehiclesRouter from './routes/vehicles.js';
import structuresRouter from './routes/structures.js';

const app = express();

app.use(cors());
app.use(express.json());

app.get('/health', (_req, res) => {
  res.json({ status: 'ok' });
});

app.use('/api/towers', towersRouter);
app.use('/api/vehicles', vehiclesRouter);
app.use('/api/structures', structuresRouter);

app.use((err, _req, res, _next) => {
  console.error(err);
  res.status(err.status || 500).json({
    message: err.message || 'Internal server error',
  });
});

export default app;
