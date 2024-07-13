// controllers/authController.js
const bcrypt = require('bcryptjs');
const UserModel = require('../models/User.model');

exports.ping = async (request, reply) => {
  reply.send({ success: "true" });
};

exports.register = async (request, reply) => {
  const { username, password } = request.body;
  const hashedPassword = await bcrypt.hash(password, 10);

  try {
    const checkUser = await UserModel.findOne({ username });
    if (checkUser) {
      reply.code(400).send({ error: true, message: 'Username already exist' });
      return;
    }

    const user = new UserModel({ username, password: hashedPassword });
    await user.save();
    reply.code(200).send({ error: false, message: "Register success" });
  } catch (err) {
    reply.code(500).send({ error: true, message: err.message });
  }
};

exports.login = async (request, reply) => {
  const { username, password } = request.body;

  try {
    const user = await UserModel.findOne({ username });
    if (!user) {
      reply.code(400).send({ error: true, message: 'User not found' });
      return;
    }

    const isValid = await bcrypt.compare(password, user.password);
    if (!isValid) {
      reply.code(400).send({ error: true, message: 'Invalid password' });
      return;
    }

    const token = await reply.jwtSign({username})
    reply.code(200).send({ error: false, message: "Login success", data: { jwt: token } });
  } catch (err) {
    reply.code(500).send({ error: true, message: err.message });
  }
};
