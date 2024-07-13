require('dotenv').config();
const fastify = require('fastify')({ logger: true });
const mongoose = require('mongoose');
const fastifySwagger = require('@fastify/swagger');
const swaggerUi = require('@fastify/swagger-ui');
const jwt = require('fastify-jwt');

fastify.register(jwt, { secret: process.env.JWT_SECRET });
fastify.register(fastifySwagger, {
  swagger: {
    info: {
      title: 'Auth Service API',
      description: 'API documentation for the Auth Service',
      version: '1.0.0',
    },
    host: 'localhost:8081',
    schemes: ['http'],
    consumes: ['application/json'],
    produces: ['application/json'],
  },
});

fastify.register(swaggerUi, {
  routePrefix: '/documentation',
  uiConfig: {
    docExpansion: 'full',
    deepLinking: false
  },
  uiHooks: {
    onRequest: function (request, reply, next) { next(); },
    preHandler: function (request, reply, next) { next(); }
  },
  staticCSP: true,
  transformStaticCSP: (header) => header,
  transformSpecification: (swaggerObject, request, reply) => { return swaggerObject },
  transformSpecificationClone: true
});



mongoose.connect(process.env.MONGO_URI, { useNewUrlParser: true, useUnifiedTopology: true });

const authRoutes = require('./routes/Auth.route');
fastify.register(authRoutes);


const start = async () => {
  try {
    console.log("try to start at ", process.env.APP_PORT || 8081);
    await fastify.listen({ port: process.env.APP_PORT || 8081, host: '0.0.0.0' });
    fastify.log.info(`Server listening on ${fastify.server.address().port}`);
  } catch (err) {
    console.log("app crash");
    fastify.log.error(err);
    process.exit(1);
  }
};
start();