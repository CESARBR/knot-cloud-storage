import mongoose from 'mongoose';

class MongoDatabase {
  constructor(hostname, port, name) {
    this.url = `mongodb://${hostname}:${port}/${name}`;
  }

  async start() {
    this.mongoose = await mongoose.connect(this.url);
    return this.mongoose;
  }
}

export default MongoDatabase;
