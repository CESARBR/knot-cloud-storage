class DataController {
  constructor(saveDataInteractor) {
    this.saveDataInteractor = saveDataInteractor;
  }

  async save(request, h) {
    try {
      await this.saveDataInteractor.execute(request.payload);
      return h.response().code(201);
    } catch (err) {
      return h.response().code(400);
    }
  }
}

export default DataController;
