import Joi from 'joi';
import httpSignature from 'http-signature';

const messageSchema = Joi.object().keys({
  data: Joi.object().keys({
    devices: Joi.array(),
    topic: Joi.string().required(),
    payload: Joi.object().required(),
  }),
  metadata: Joi.object().keys({
    route: Joi.array().required(),
    date: Joi.date().required(),
  }),
});

function mapRequestToMessage(request) {
  return {
    data: request.payload,
    metadata: {
      route: JSON.parse(request.headers['x-meshblu-route']),
      date: new Date(request.headers.date),
    },
  };
}

function validateMessage(message) {
  const { error } = Joi.validate(message, messageSchema);
  if (error) {
    throw error;
  }
}

function verifySignature(request, publicKey) {
  const parsedReq = httpSignature.parseRequest(request);
  if (!httpSignature.verifySignature(parsedReq, Buffer.from(publicKey, 'base64').toString('ascii'))) {
    throw new Error('Signature failed');
  }
}

class DataController {
  constructor(settings, saveDataInteractor, listDataInteractor, logger) {
    this.publicKey = settings.server.publicKey;
    this.saveDataInteractor = saveDataInteractor;
    this.listDataInteractor = listDataInteractor;
    this.logger = logger;
  }

  async save(request, h) {
    try {
      verifySignature(request, this.publicKey);
      const message = mapRequestToMessage(request);
      validateMessage(message);
      await this.saveDataInteractor.execute(message);
      this.logger.info('Data saved');
      return h.response().code(201);
    } catch (err) {
      this.logger.error(`Failed saving data: ${err.message}`);
      return h.response().code(400);
    }
  }

  async list(request, h) {
    const credentials = {
      uuid: request.headers.auth_id,
      token: request.headers.auth_token,
    };

    try {
      const data = await this.listDataInteractor.execute(credentials, request.query);
      this.logger.info('Data obtained');
      return h.response(data).code(200);
    } catch (error) {
      this.logger.error(`Failed to list data (${error.code || 500}): ${error.message}`);
      return h.response(error.message).code(error.code);
    }
  }

  async listByDevice(request, h) {
    const credentials = {
      uuid: request.headers.auth_id,
      token: request.headers.auth_token,
    };
    const dataQuery = request.query;
    dataQuery.from = request.params.id;

    try {
      const data = await this.listDataInteractor.execute(credentials, dataQuery);
      this.logger.info('Data obtained');
      return h.response(data).code(200);
    } catch (error) {
      this.logger.error(`Failed to list data (${error.code || 500}): ${error.message}`);
      return h.response(error.message).code(error.code);
    }
  }

  async listBySensor(request, h) {
    const credentials = {
      uuid: request.headers.auth_id,
      token: request.headers.auth_token,
    };
    const dataQuery = request.query;
    dataQuery.from = request.params.deviceId;
    dataQuery.sensorId = request.params.sensorId;

    try {
      const data = await this.listDataInteractor.execute(credentials, dataQuery);
      this.logger.info('Data obtained');
      return h.response(data).code(200);
    } catch (error) {
      this.logger.error(`Failed to list data (${error.code || 500}): ${error.message}`);
      return h.response(error.message).code(error.code);
    }
  }
}

export default DataController;
