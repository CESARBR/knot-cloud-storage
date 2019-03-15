import _ from 'lodash';

function throwError(message, code) {
  const error = new Error(message);
  error.code = code;
  throw error;
}

class SaveData {
  constructor(dataStore) {
    this.dataStore = dataStore;
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

    await this.dataStore.save({
      from: routeData.from,
      payload: message.data.payload,
    });
  }
}

export default SaveData;
