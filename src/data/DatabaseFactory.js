import MongoDatabase from 'data/MongoDatabase';

class DatabaseFactory {
  constructor(type, settings) {
    this.type = type;
    this.settings = settings;
  }

  create() {
    switch (this.type) {
      case 'MONGO':
        return new MongoDatabase(this.settings);
      default:
        throw new Error('Unknown cloud');
    }
  }
}

export default DatabaseFactory;
