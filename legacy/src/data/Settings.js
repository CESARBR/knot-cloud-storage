class Settings {
  constructor(database, databaseConfig, server, meshblu, logger) {
    this.database = database;
    this.databaseConfig = databaseConfig;
    this.server = server;
    this.meshblu = meshblu;
    this.logger = logger;
  }
}

export default Settings;
