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
  constructor(saveDataInteractor) {
    this.saveDataInteractor = saveDataInteractor;
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
}

export default DataController;
