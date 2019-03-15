import SettingsFactory from 'data/SettingsFactory';
import MongoDatabaseFactory from 'data/MongoDatabaseFactory';

async function main() {
  const settings = new SettingsFactory().create();
  const database = new MongoDatabaseFactory(settings).create();
  await database.start();
}

main();
