class Cloud {
  constructor(requester, uuidAliasResolver) {
    this.requester = requester;
    this.uuidAliasResolver = uuidAliasResolver;
  }

  async getDevice(credentials, id, as) {
    const uuid = await this.resolveUuidAlias(id);
    const request = {
      metadata: {
        jobType: 'GetDevice',
        auth: credentials,
        toUuid: uuid,
        fromUuid: as,
      },
    };

    const response = await this.sendRequest(request);
    this.checkResponseHasError(response, 200);
    return JSON.parse(response.rawData);
  }

  async getDevices(credentials, query) {
    const device = await this.getDevice(credentials, credentials.uuid);
    const request = {
      metadata: {
        jobType: 'SearchDevices',
        auth: credentials,
        fromUuid: device.type === 'knot:app' ? device.knot.router : credentials.uuid,
      },
      data: query,
    };

    const response = await this.sendRequest(request);
    this.checkResponseHasError(response, 200);
    return JSON.parse(response.rawData);
  }

  async resolveUuidAlias(id) {
    return new Promise((resolve, reject) => {
      this.uuidAliasResolver.resolve(id, (error, uuid) => {
        if (error) {
          reject(error);
          return;
        }

        resolve(uuid);
      });
    });
  }

  async sendRequest(request) {
    let response;

    try {
      response = await this.send(request);
    } catch (error) {
      this.throwError('Bad Gateway', 502);
    }

    return response;
  }

  throwError(message, code) {
    const error = new Error(message);
    error.code = code;
    throw error;
  }

  checkResponseHasError(response, successCode) {
    if (!response) {
      this.throwError('Gateway Timeout', 504);
    }

    if (response.metadata.code !== successCode) {
      this.throwError(response.metadata.status, response.metadata.code);
    }
  }

  async send(request) {
    return new Promise((resolve, reject) => {
      this.requester.do(request, (error, response) => {
        if (error) {
          reject(error);
        } else {
          resolve(response);
        }
      });
    });
  }
}

export default Cloud;
