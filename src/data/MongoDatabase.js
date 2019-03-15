import mongoose from 'mongoose';

class MongoDatabase {
  constructor(settings) {
    this.url = `mongodb://${settings.hostname}:${settings.port}/${settings.name}`;
  }

  async start() {
    this.mongoose = await mongoose.connect(this.url);
    return this.mongoose;
  }
}

export default MongoDatabase;
