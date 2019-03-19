import SettingsFactory from 'data/SettingsFactory';
import DatabaseFactory from 'data/DatabaseFactory';
import CloudRequesterFactory from 'network/CloudRequesterFactory';

import HapiFactory from 'server/HapiFactory';

async function main() {
  const settings = new SettingsFactory().create();
  const database = new DatabaseFactory(settings.database, settings.databaseConfig).create();
  const cloudRequester = new CloudRequesterFactory(settings).create();
  const server = new HapiFactory(settings, database, cloudRequester).create();

  await cloudRequester.start();
  await database.start();
  await server.start();
}

main();
