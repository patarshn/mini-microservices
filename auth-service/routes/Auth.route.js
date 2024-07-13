// routes/authRoutes.js
const AuthController = require('../controllers/Auth.controller');

async function routes(fastify, options) {
  fastify.get('/ping', {
    schema: {
      description: 'Ping the server',
      response: {
        200: {
          type: 'object',
          properties: {
            success: { type: 'string' },
          },
        },
      },
    },
    handler: AuthController.ping
  });

  fastify.post('/register', {
    schema: {
      description: 'Register a new user',
      body: {
        type: 'object',
        required: ['username', 'password'],
        properties: {
          username: { type: 'string' },
          password: { type: 'string' },
        },
      },
      response: {
        200: {
          type: 'object',
          properties: {
            error: { type: 'boolean' },
            message: { type: 'string'}
          },
        },
        '4xx': {
          type: 'object',
          properties: {
            error: { type: 'boolean' },
            message: { type: 'string'}
          },
        },
        '5xx': {
          type: 'object',
          properties: {
            error: { type: 'boolean' },
            message: { type: 'string'}
          },
        },
      },
    },
    handler: AuthController.register
  });

  fastify.post('/login', {
    schema: {
      description: 'Login a user',
      body: {
        type: 'object',
        required: ['username', 'password'],
        properties: {
          username: { type: 'string' },
          password: { type: 'string' },
        },
      },
      response: {
        200: {
          type: 'object',
          properties: {
            error: { type: 'boolean' },
            data: {
              type: 'object',
              properties: {
                jwt: { type: 'string' },
              },
            },
            message: { type: 'string' },
          },
        },
        '4xx': {
          type: 'object',
          properties: {
            error: { type: 'boolean' },
            message: { type: 'string'}
          }
        },
        '5xx': {
          type: 'object',
          properties: {
            error: { type: 'boolean' },
            message: { type: 'string'}
          }
        },
      },
    },
    handler: AuthController.login
  });
}

module.exports = routes;
