import swaggerJsdoc from 'swagger-jsdoc';

const swaggerDefinition = {
  openapi: '3.0.0',
  info: {
    title: 'Resource Manager API',
    version: '1.0.0',
    description: 'API para towers, vehicles e structures',
  },
  servers: [{ url: 'http://localhost:4000' }],
  components: {
    schemas: {
      Tower: {
        type: 'object',
        properties: {
          id: { type: 'string', format: 'uuid' },
          name: { type: 'string' },
          latitude: { type: 'number' },
          longitude: { type: 'number' },
          is_leader: { type: 'boolean' },
        },
        required: ['id', 'name', 'latitude', 'longitude', 'is_leader'],
      },
      Vehicle: {
        type: 'object',
        properties: {
          id: { type: 'string', format: 'uuid' },
          name: { type: 'string' },
          type: { type: 'string' },
          latitude: { type: 'number' },
          longitude: { type: 'number' },
        },
        required: ['id', 'name', 'type', 'latitude', 'longitude'],
      },
      Structure: {
        type: 'object',
        properties: {
          id: { type: 'string', format: 'uuid' },
          name: { type: 'string' },
          type: { type: 'string' },
          latitude: { type: 'number' },
          longitude: { type: 'number' },
        },
        required: ['id', 'name', 'type', 'latitude', 'longitude'],
      },
      Error: {
        type: 'object',
        properties: { message: { type: 'string' } },
      },
    },
  },
  paths: {
    '/health': {
      get: {
        summary: 'Health check',
        responses: { 200: { description: 'OK' } },
      },
    },
    '/api/towers': {
      get: {
        summary: 'Listar torres',
        responses: {
          200: {
            description: 'Lista de torres',
            content: {
              'application/json': {
                schema: { type: 'array', items: { $ref: '#/components/schemas/Tower' } },
              },
            },
          },
        },
      },
      post: {
        summary: 'Criar torre',
        requestBody: {
          required: true,
          content: {
            'application/json': {
              schema: {
                type: 'object',
                properties: {
                  name: { type: 'string' },
                  latitude: { type: 'number' },
                  longitude: { type: 'number' },
                  is_leader: { type: 'boolean' },
                },
                required: ['name', 'latitude', 'longitude'],
              },
            },
          },
        },
        responses: {
          201: {
            description: 'Torre criada',
            content: { 'application/json': { schema: { $ref: '#/components/schemas/Tower' } } },
          },
          400: {
            description: 'Dados invalidos',
            content: { 'application/json': { schema: { $ref: '#/components/schemas/Error' } } },
          },
        },
      },
    },
    '/api/towers/{id}': {
      get: {
        summary: 'Obter torre por ID',
        parameters: [{ name: 'id', in: 'path', required: true, schema: { type: 'string', format: 'uuid' } }],
        responses: {
          200: {
            description: 'Torre encontrada',
            content: { 'application/json': { schema: { $ref: '#/components/schemas/Tower' } } },
          },
          404: {
            description: 'Nao encontrada',
            content: { 'application/json': { schema: { $ref: '#/components/schemas/Error' } } },
          },
        },
      },
      delete: {
        summary: 'Remover torre',
        parameters: [{ name: 'id', in: 'path', required: true, schema: { type: 'string', format: 'uuid' } }],
        responses: {
          204: { description: 'Removida' },
          404: {
            description: 'Nao encontrada',
            content: { 'application/json': { schema: { $ref: '#/components/schemas/Error' } } },
          },
        },
      },
    },
    '/api/vehicles': {
      get: {
        summary: 'Listar veiculos',
        responses: {
          200: {
            description: 'Lista de veiculos',
            content: {
              'application/json': {
                schema: { type: 'array', items: { $ref: '#/components/schemas/Vehicle' } },
              },
            },
          },
        },
      },
      post: {
        summary: 'Criar veiculo',
        requestBody: {
          required: true,
          content: {
            'application/json': {
              schema: {
                type: 'object',
                properties: {
                  name: { type: 'string' },
                  type: { type: 'string' },
                  latitude: { type: 'number' },
                  longitude: { type: 'number' },
                },
                required: ['name', 'type', 'latitude', 'longitude'],
              },
            },
          },
        },
        responses: {
          201: {
            description: 'Veiculo criado',
            content: { 'application/json': { schema: { $ref: '#/components/schemas/Vehicle' } } },
          },
          400: {
            description: 'Dados invalidos',
            content: { 'application/json': { schema: { $ref: '#/components/schemas/Error' } } },
          },
        },
      },
    },
    '/api/vehicles/{id}': {
      get: {
        summary: 'Obter veiculo por ID',
        parameters: [{ name: 'id', in: 'path', required: true, schema: { type: 'string', format: 'uuid' } }],
        responses: {
          200: {
            description: 'Veiculo encontrado',
            content: { 'application/json': { schema: { $ref: '#/components/schemas/Vehicle' } } },
          },
          404: {
            description: 'Nao encontrado',
            content: { 'application/json': { schema: { $ref: '#/components/schemas/Error' } } },
          },
        },
      },
    },
    '/api/structures': {
      get: {
        summary: 'Listar estruturas',
        responses: {
          200: {
            description: 'Lista de estruturas',
            content: {
              'application/json': {
                schema: { type: 'array', items: { $ref: '#/components/schemas/Structure' } },
              },
            },
          },
        },
      },
      post: {
        summary: 'Criar estrutura',
        requestBody: {
          required: true,
          content: {
            'application/json': {
              schema: {
                type: 'object',
                properties: {
                  name: { type: 'string' },
                  type: { type: 'string' },
                  latitude: { type: 'number' },
                  longitude: { type: 'number' },
                },
                required: ['name', 'type', 'latitude', 'longitude'],
              },
            },
          },
        },
        responses: {
          201: {
            description: 'Estrutura criada',
            content: { 'application/json': { schema: { $ref: '#/components/schemas/Structure' } } },
          },
          400: {
            description: 'Dados invalidos',
            content: { 'application/json': { schema: { $ref: '#/components/schemas/Error' } } },
          },
        },
      },
    },
    '/api/structures/{id}': {
      get: {
        summary: 'Obter estrutura por ID',
        parameters: [{ name: 'id', in: 'path', required: true, schema: { type: 'string', format: 'uuid' } }],
        responses: {
          200: {
            description: 'Estrutura encontrada',
            content: { 'application/json': { schema: { $ref: '#/components/schemas/Structure' } } },
          },
          404: {
            description: 'Nao encontrada',
            content: { 'application/json': { schema: { $ref: '#/components/schemas/Error' } } },
          },
        },
      },
    },
  },
};

const options = { definition: swaggerDefinition, apis: [] };
const swaggerSpec = swaggerJsdoc(options);

export default swaggerSpec;
