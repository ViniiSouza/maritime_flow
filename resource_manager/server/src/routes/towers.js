import { Router } from 'express';
import pool from '../db.js';

const router = Router();

router.get('/', async (_req, res, next) => {
  try {
    const { rows } = await pool.query(
      'SELECT id, name, latitude, longitude, is_leader FROM towers ORDER BY id ASC'
    );
    res.json(rows);
  } catch (error) {
    next(error);
  }
});

router.get('/:id', async (req, res, next) => {
  try {
    const { id } = req.params;
    const { rows } = await pool.query(
      'SELECT id, name, latitude, longitude, is_leader FROM towers WHERE id = $1',
      [id]
    );
    if (!rows.length) {
      return res.status(404).json({ message: 'Tower not found' });
    }
    res.json(rows[0]);
  } catch (error) {
    next(error);
  }
});

router.post('/', async (req, res, next) => {
  try {
    const { name, latitude, longitude, is_leader = false } = req.body;
    if (!name || latitude === undefined || longitude === undefined) {
      return res.status(400).json({ message: 'name, latitude and longitude are required' });
    }
    const { rows } = await pool.query(
      `INSERT INTO towers (name, latitude, longitude, is_leader)
       VALUES ($1, $2, $3, $4)
       RETURNING id, name, latitude, longitude, is_leader`,
      [name, latitude, longitude, is_leader]
    );
    res.status(201).json(rows[0]);
  } catch (error) {
    next(error);
  }
});

router.delete('/:id', async (req, res, next) => {
  try {
    const { id } = req.params;
    const { rowCount } = await pool.query('DELETE FROM towers WHERE id = $1', [id]);
    if (!rowCount) {
      return res.status(404).json({ message: 'Tower not found' });
    }
    res.status(204).send();
  } catch (error) {
    next(error);
  }
});

export default router;
