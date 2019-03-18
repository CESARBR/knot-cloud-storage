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
    const data = await this.listDataInteractor.execute(request.query);
    return h.response(data).code(200);
  }

  async listByDevice(request, h) {
    const dataQuery = request.query;
    dataQuery.from = request.params.id;
    const data = await this.listDataInteractor.execute(dataQuery);
    return h.response(data).code(200);
  }
}

export default DataController;
