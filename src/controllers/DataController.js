import Joi from 'joi';
import _ from 'lodash';

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

function mapJoiError(propertyName, error) {
  const reasons = _.map(error.details, 'message');
  const formattedReasons = reasons.length > 1
    ? `\n${_.chain(reasons).map(reason => `- ${reason}`).join('\n').value()}`
    : reasons[0];
  return new Error(`Invalid "${propertyName}" property: ${formattedReasons}`);
}

function validate(propertyName, propertyValue, schema) {
  const { error } = Joi.validate(propertyValue, schema, { abortEarly: false });
  if (error) {
    throw mapJoiError(propertyName, error);
  }
}


class DataController {
  constructor(saveDataInteractor, listDataInteractor) {
    this.saveDataInteractor = saveDataInteractor;
    this.listDataInteractor = listDataInteractor;
  }

  async save(request, h) {
    try {
      const message = {
        data: request.payload,
        metadata: {
          route: JSON.parse(request.headers['x-meshblu-route']),
          date: new Date(request.headers.date),
        },
      };

      validate('message', message, messageSchema);
      await this.saveDataInteractor.execute(message);
      return h.response().code(201);
    } catch (err) {
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
      return h.response(data).code(200);
    } catch (error) {
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
      return h.response(data).code(200);
    } catch (error) {
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
      return h.response(data).code(200);
    } catch (error) {
      return h.response(error.message).code(error.code);
    }
  }
}

export default DataController;
