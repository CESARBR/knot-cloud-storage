class ListData {
  constructor(dataStore, cloud) {
    this.dataStore = dataStore;
    this.cloud = cloud;
  }

  async execute(credentials, query) {
    if (!query.from) {
      return this.getDevicesData(credentials, query);
    }

    await this.cloud.getDevice(credentials, query.from);
    return this.dataStore.list(query);
  }

  async getDevicesData(credentials, query) {
    const dataQuery = query;
    const devices = await this.cloud.getDevices(credentials, { type: 'knot:thing' });
    const data = await Promise.all(devices.map(async (device) => {
      dataQuery.from = device.knot.id;
      return this.dataStore.list(dataQuery);
    }));

    return [].concat(...data);
  }
}

export default ListData;
