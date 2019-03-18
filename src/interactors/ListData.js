class ListData {
  constructor(dataStore) {
    this.dataStore = dataStore;
  }

  async execute(query) {
    return this.dataStore.list(query);
  }
}

export default ListData;
