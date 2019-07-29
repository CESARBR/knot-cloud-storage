class ListData {
  constructor(dataStore, cloud) {
    this.dataStore = dataStore;
    this.cloud = cloud;
  }

  async execute(credentials, query) {
    if (!query.from) {
      return this.getDevicesData(credentials, query);
    }

    await this.verifyAuthorization(credentials, query.from);
    return this.dataStore.list(query);
  }

  async verifyAuthorization(credentials, deviceId) {
    const authDevice = await this.cloud.getDevice(credentials, credentials.uuid);
    return authDevice.type === 'knot:app'
      ? this.cloud.getDevice(credentials, deviceId, authDevice.knot.router)
      : this.cloud.getDevice(credentials, deviceId);
  }

  async getDevicesData(credentials, query) {
    const queryBase = query;
    const devices = await this.cloud.getDevices(credentials, { type: 'knot:thing' });
    queryBase.deviceIds = devices.map(device => device.knot.id);
    return this.dataStore.list(queryBase);
  }
}

export default ListData;
