import SettingsFactory from 'data/SettingsFactory';
import MongoDatabaseFactory from 'data/MongoDatabaseFactory';
import CloudRequesterFactory from 'network/CloudRequesterFactory';

import HapiFactory from 'server/HapiFactory';

async function main() {
  const settings = new SettingsFactory().create();
  const database = new MongoDatabaseFactory(settings).create();
  const cloudRequester = new CloudRequesterFactory(settings).create();

  const server = new HapiFactory(settings, database, cloudRequester).create();

  await cloudRequester.start();
  await database.start();
  await server.start();
}

main();
