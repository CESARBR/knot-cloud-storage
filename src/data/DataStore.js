import mongoose from 'mongoose';
import _ from 'lodash';

const DataSchema = mongoose.Schema({
  from: {
    type: String,
    required: true,
  },
  timestamp: {
    type: Date,
    default: Date.now,
    required: true,
  },
  payload: {
    type: Object,
    required: true,
  },
});

function extractQueryOptions(query) {
  return {
    startDate: query.startDate,
    finishDate: query.finishDate,
    skip: parseInt(query.skip, 10) || 0,
    take: parseInt(query.take, 10) || 10,
    order: query.order || 1,
    orderBy: query.orderBy,
  };
}

class DataStore {
  constructor(database) {
    this.database = database;
  }

  async save(data) {
    await this.database.save('Data', DataSchema, data);
  }

  async list(query) {
    const queryOptions = extractQueryOptions(query);
    const queryBase = _.omit(query, ['order', 'orderBy', 'skip', 'take', 'startDate, finishDate']);

    queryOptions.take = queryOptions.take > 100 ? 100 : queryOptions.take;

    if (queryBase.from && queryBase.sensorId) {
      const sensorId = parseInt(queryBase.sensorId, 10);
      const combinedQuery = { $and: [{ from: queryBase.from }, { 'payload.sensorId': sensorId }] };
      return this.database.find('Data', DataSchema, combinedQuery, queryOptions);
    }

    return this.database.find('Data', DataSchema, queryBase, queryOptions);
  }
}

export default DataStore;
