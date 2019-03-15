import mongoose from 'mongoose';

class MongoDatabase {
  constructor(hostname, port, name) {
    this.url = `mongodb://${hostname}:${port}/${name}`;
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
}

export default MongoDatabase;
