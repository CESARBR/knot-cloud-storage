class DataController {
  constructor(saveDataInteractor, listDataInteractor) {
    this.saveDataInteractor = saveDataInteractor;
    this.listDataInteractor = listDataInteractor;
  }

  async save(request, h) {
    try {
      await this.saveDataInteractor.execute(request.payload);
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
    const dataQuery = request.query;
    dataQuery.from = request.params.deviceId;
    dataQuery.sensorId = request.params.sensorId;

    try {
      const data = await this.listDataInteractor.execute(dataQuery);
      return h.response(data).code(200);
    } catch (error) {
      return h.response(error.message).code(error.code);
    }
  }
}

export default DataController;
