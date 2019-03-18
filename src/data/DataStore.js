import mongoose from 'mongoose';

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
    start: query.start,
    finish: query.finish,
    limit: parseInt(query.limit, 10) || 10,
    order: query.order || 1,
    orderBy: query.orderBy,
  };
}

function extractQueryBase(query) {
  const base = query;
  delete base.start;
  delete base.finish;
  delete base.limit;
  delete base.order;
  delete base.orderBy;
  return base;
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
    const queryBase = extractQueryBase(query);

    return this.database.find('Data', DataSchema, queryBase, queryOptions);
  }
}

export default DataStore;
