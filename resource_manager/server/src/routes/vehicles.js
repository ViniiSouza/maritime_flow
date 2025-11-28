import { Router } from 'express';
import { randomUUID } from 'crypto';
import pool from '../db.js';
import { commitManifest, deleteResource } from '../gitops.js';

const router = Router();

router.get('/', async (_req, res, next) => {
  try {
    const { rows } = await pool.query(
      'SELECT id, name, type, latitude, longitude FROM vehicles ORDER BY id ASC'
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
      'SELECT id, name, type, latitude, longitude FROM vehicles WHERE id = $1',
      [id]
    );
    if (!rows.length) {
      return res.status(404).json({ message: 'Vehicle not found' });
    }
    res.json(rows[0]);
  } catch (error) {
    next(error);
  }
});

router.post('/', async (req, res, next) => {
  let id = randomUUID();
  let { name, type, latitude, longitude } = req.body;
  try {
    if (!name || !type || latitude === undefined || longitude === undefined) {
      return res
        .status(400)
        .json({ message: 'name, type, latitude and longitude are required' });
    }
    const { rows } = await pool.query(
      `INSERT INTO vehicles (id, name, type, latitude, longitude)
       VALUES ($1, $2, $3, $4, $5)
       RETURNING id, name, type, latitude, longitude`,
      [id, name, type, latitude, longitude]
    );
    res.status(201).json(rows[0]);
  } catch (error) {
    next(error);
  }

  const directoryType = (type.toLowerCase() == 'helicopter' ? 'helicopters' : 'ships');

  commitManifest({
    directory: `mobility-core/${directoryType}/${id}`,
    files: {
      "mobility-core/deployment.yaml.tpl": `mobility-core/${directoryType}/${id}/deployment.yaml`,
    },
    replacements: {
      VID: id,
      VTYPE: type.toLowerCase(),
      LATITUDE: String(latitude),
      LONGITUDE: String(longitude),
      VELOCITY: '800'
    }
  }).catch(console.error);
});

router.delete('/:id', async (req, res, next) => {
  let { id } = req.params;
  let type;
  try {
    const { rows } = await pool.query('SELECT type FROM vehicles WHERE id = $1', [id]);
    if (!rows.length) {
      return res.status(404).json({ message: 'Vehicle not found' });
    }

    type = rows[0].type;
    const { rowCount } = await pool.query('DELETE FROM vehicles WHERE id = $1', [id]);
    res.status(204).send();
  } catch (error) {
    next(error);
  }

  const directoryType = (type.toLowerCase() == 'helicopter' ? 'helicopters' : 'ships');
  deleteResource(`mobility-core/${directoryType}/${id}`).catch(console.error);
});

export default router;
