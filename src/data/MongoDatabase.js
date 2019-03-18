import mongoose from 'mongoose';
import moment from 'moment';
import _ from 'lodash';

function getISODate(date) {
  return `${moment(date).format('YYYY-MM-DDTHH:mm:ss')}Z`;
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
    const queryBase = query;
    const queryOptions = options;

    if (queryBase.from && queryBase.sensorId) {
      const sensorId = parseInt(queryBase.sensorId, 10);
      const combinedQuery = { $and: [{ from: queryBase.from }, { 'payload.sensorId': sensorId }] };
      return Model.find(combinedQuery, { _id: 0, __v: 0 }, queryOptions).exec();
    }

    if (queryOptions.orderBy) {
      queryOptions.sort = {};
      queryOptions.sort[queryOptions.orderBy] = queryOptions.order;
    }
    if (queryOptions.start) {
      _.merge(queryBase, { timestamp: { $gte: getISODate(queryOptions.start) } });
    }
    if (queryOptions.finish) {
      _.merge(queryBase, { timestamp: { $lte: getISODate(queryOptions.finish) } });
    }

    return Model.find(queryBase, { _id: 0, __v: 0 }, queryOptions).exec();
  }
}

export default MongoDatabase;
