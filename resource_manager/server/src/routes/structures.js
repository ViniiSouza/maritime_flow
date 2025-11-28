import { Router } from 'express';
import { randomUUID } from 'crypto';
import pool from '../db.js';
import { commitManifest, deleteResource } from '../gitops.js';

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
  let structureId = randomUUID();
  let { name, type, latitude, longitude } = req.body;
  try {
    client = await pool.connect();
    if (!name || !type || latitude === undefined || longitude === undefined) {
      return res
        .status(400)
        .json({ message: 'name, type, latitude and longitude are required' });
    }
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
      [dockSlotId, structureId, 1, 'dock', helipadSlotId, 1, 'helipad']
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

  const directoryType = (type.toLowerCase() == 'platform' ? 'platforms' : 'centrals');

  commitManifest({
    directory: `station-core/${directoryType}/${structureId}`,
    files: {
      "station-core/deployment.yaml.tpl": `station-core/${directoryType}/${structureId}/deployment.yaml`,
      "station-core/svc.yaml.tpl": `station-core/${directoryType}/${structureId}/svc.yaml`,
    },
    replacements: {
      SID: structureId,
      STYPE: type.toLowerCase(),
    }
  }).catch(console.error);
});

router.delete('/:id', async (req, res, next) => {
  let client;
  let { id } = req.params;
  let type;
  try {
    client = await pool.connect();
    await client.query('BEGIN');

    const { rows } = await client.query(
      'SELECT type FROM structures WHERE id = $1',
      [id]
    );
    if (!rows.length) {
      await client.query('ROLLBACK');
      return res.status(404).json({ message: 'Structure not found' });
    }

    type = rows[0].type;

    await client.query('DELETE FROM slots WHERE structure_id = $1', [id]);
    await client.query('DELETE FROM structures WHERE id = $1', [id]);

    await client.query('COMMIT');
    res.status(204).send();
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

  const directoryType = (type.toLowerCase() == 'platform' ? 'platforms' : 'centrals');
  deleteResource(`station-core/${directoryType}/${id}`).catch(console.error);
});

export default router;
