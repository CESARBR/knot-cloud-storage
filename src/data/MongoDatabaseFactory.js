import MongoDatabase from 'data/MongoDatabase';

class MongoDatabaseFactory {
  constructor(settings) {
    this.settings = settings;
  }

  create() {
    return new MongoDatabase(
      this.settings.database.hostname,
      this.settings.database.port,
      this.settings.database.name,
    );
  }
}

export default MongoDatabaseFactory;
