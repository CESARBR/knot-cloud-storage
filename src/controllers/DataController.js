import Joi from 'joi';

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
      const { error } = Joi.validate(message, messageSchema);
      if (error) {
        throw error;
      }

      await this.saveDataInteractor.execute(message);
      return h.response().code(201);
    } catch (err) {
      return h.response().code(400);
    }
  }
}

export default DataController;
