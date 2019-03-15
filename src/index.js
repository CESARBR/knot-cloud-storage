import SettingsFactory from 'data/SettingsFactory';
import DatabaseFactory from 'data/DatabaseFactory';

async function main() {
  const settings = new SettingsFactory().create();
  const database = new DatabaseFactory(settings.databaseType, settings.database).create();
  await database.start();
}

main();
