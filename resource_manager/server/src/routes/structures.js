import { Router } from 'express';
import { randomUUID } from 'crypto';
import pool from '../db.js';

const router = Router();

router.get('/', async (_req, res, next) => {
  try {
    const { rows } = await pool.query(
      `SELECT s.id,
              s.name,
              s.type,
              s.latitude,
              s.longitude,
              COALESCE(
                json_agg(
                  json_build_object(
                    'id', sl.id,
                    'structure_id', sl.structure_id,
                    'number', sl.number,
                    'type', sl.type
                  )
                ) FILTER (WHERE sl.id IS NOT NULL),
                '[]'
              ) AS slots
       FROM structures s
       LEFT JOIN slots sl ON sl.structure_id = s.id
       GROUP BY s.id
       ORDER BY s.id ASC`
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
      `SELECT s.id,
              s.name,
              s.type,
              s.latitude,
              s.longitude,
              COALESCE(
                json_agg(
                  json_build_object(
                    'id', sl.id,
                    'structure_id', sl.structure_id,
                    'number', sl.number,
                    'type', sl.type
                  )
                ) FILTER (WHERE sl.id IS NOT NULL),
                '[]'
              ) AS slots
       FROM structures s
       LEFT JOIN slots sl ON sl.structure_id = s.id
       WHERE s.id = $1
       GROUP BY s.id`,
      [id]
    );
    if (!rows.length) {
      return res.status(404).json({ message: 'Structure not found' });
    }
    res.json(rows[0]);
  } catch (error) {
    next(error);
  }
});

router.post('/', async (req, res, next) => {
  let client;
  try {
    client = await pool.connect();
    const { name, type, latitude, longitude } = req.body;
    if (!name || !type || latitude === undefined || longitude === undefined) {
      return res
        .status(400)
        .json({ message: 'name, type, latitude and longitude are required' });
    }
    const structureId = randomUUID();
    const dockSlotId = randomUUID();
    const helipadSlotId = randomUUID();

    await client.query('BEGIN');

    const { rows: structureRows } = await client.query(
      `INSERT INTO structures (id, name, type, latitude, longitude)
       VALUES ($1, $2, $3, $4, $5)
       RETURNING id, name, type, latitude, longitude`,
      [structureId, name, type, latitude, longitude]
    );

    const { rows: slotRows } = await client.query(
      `INSERT INTO slots (id, structure_id, number, type)
       VALUES ($1, $2, $3, $4),
              ($5, $2, $6, $7)
       RETURNING id, structure_id, number, type`,
      [dockSlotId, structureId, 1, 'dock', helipadSlotId, 2, 'helipad']
    );

    await client.query('COMMIT');

    res.status(201).json({ ...structureRows[0], slots: slotRows });
  } catch (error) {
    if (client) {
      await client.query('ROLLBACK');
    }
    next(error);
  } finally {
    if (client) {
      client.release();
    }
  }
});

export default router;
