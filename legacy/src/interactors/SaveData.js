import _ from 'lodash';

function throwError(message, code) {
  const error = new Error(message);
  error.code = code;
  throw error;
}

class SaveData {
  constructor(dataStore, uuidAliasResolver) {
    this.dataStore = dataStore;
    this.uuidAliasResolver = uuidAliasResolver;
  }

  async execute(message) {
    if (message.data.topic !== 'data') {
      throwError('Received message isn\'t a data', 400);
    }

    const types = ['broadcast.sent', 'message.received'];
    const routeData = _.chain(message.metadata.route)
      .filter(data => _.includes(types, data.type))
      .head()
      .value();

    const id = await this.reverseUuidLookup(routeData.from);

    await this.dataStore.save({
      from: id,
      payload: message.data.payload,
    });
  }

  async reverseUuidLookup(uuid) {
    return new Promise((resolve, reject) => {
      this.uuidAliasResolver.reverseLookup(uuid, (error, aliases) => {
        if (error) {
          reject(error);
          return;
        }

        resolve(aliases[0]);
      });
    });
  }
}

export default SaveData;
