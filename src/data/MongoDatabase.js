import mongoose from 'mongoose';
import moment from 'moment';
import _ from 'lodash';

function getISODate(date) {
  return `${moment(date).format('YYYY-MM-DDTHH:mm:ss')}Z`;
}

function buildQueryOptions(query) {
  const options = query;

  if (query.take) {
    options.limit = query.take;
  }
  if (query.orderBy) {
    options.sort = {};
    options.sort[query.orderBy] = query.order;
  }

  return options;
}

function buildQueryBase(base, options) {
  const queryBase = base;

  if (options.startDate) {
    _.merge(queryBase, { timestamp: { $gte: getISODate(options.startDate) } });
  }
  if (options.finishDate) {
    _.merge(queryBase, { timestamp: { $lte: getISODate(options.finishDate) } });
  }

  return queryBase;
}

class MongoDatabase {
  constructor(settings) {
    this.url = `mongodb://${settings.hostname}:${settings.port}/${settings.name}`;
  }

  async start() {
    this.mongoose = await mongoose.connect(this.url);
    return this.mongoose;
  }

  async save(name, schema, data) {
    const Model = this.mongoose.model(name, schema);
    const dataModel = new Model(data);
    return dataModel.save();
  }

  async find(name, schema, query, options) {
    const Model = this.mongoose.model(name, schema);
    const queryBase = buildQueryBase(query, options);
    const queryOptions = buildQueryOptions(options);

    return Model.find(queryBase, { _id: 0, __v: 0 }, queryOptions).exec();
  }
}

export default MongoDatabase;
