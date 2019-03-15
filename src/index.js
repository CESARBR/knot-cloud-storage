import SettingsFactory from 'data/SettingsFactory';
import MongoDatabaseFactory from 'data/MongoDatabaseFactory';

import HapiFactory from 'server/HapiFactory';

async function main() {
  const settings = new SettingsFactory().create();
  const database = new MongoDatabaseFactory(settings).create();

  const server = new HapiFactory(settings, database).create();

  await database.start();
  await server.start();
}

main();
