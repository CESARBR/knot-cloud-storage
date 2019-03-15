import SettingsFactory from 'data/SettingsFactory';
import DatabaseFactory from 'data/DatabaseFactory';

import HapiFactory from 'server/HapiFactory';

async function main() {
  const settings = new SettingsFactory().create();
  const database = new DatabaseFactory(settings.database, settings.databaseConfig).create();
  const server = new HapiFactory(settings).create();

  await database.start();
  await server.start();
}

main();
