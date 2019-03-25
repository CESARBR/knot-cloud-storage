import Hapi from 'server/Hapi';
import SaveDataInteractor from 'interactors/SaveData';
import ListDataInteractor from 'interactors/ListData';
import DataController from 'controllers/DataController';
import DataStore from 'data/DataStore';
import UuidAliasResolverFactory from 'network/UuidAliasResolverFactory';
import CloudFactory from 'network/CloudFactory';
import LoggerFactory from 'LoggerFactory';

class HapiFactory {
  constructor(settings, database, cloudRequester) {
    this.settings = settings;
    this.database = database;
    this.cloudRequester = cloudRequester;
  }

  create() {
    const dataStore = new DataStore(this.database);
    const uuidAliasResolver = new UuidAliasResolverFactory(this.settings).create();
    const cloud = new CloudFactory(this.cloudRequester, uuidAliasResolver).create();
    const saveDataInteractor = new SaveDataInteractor(dataStore, uuidAliasResolver);
    const listDataInteractor = new ListDataInteractor(dataStore, cloud);
    const logger = new LoggerFactory(this.settings).create('');
    const dataController = new DataController(
      this.settings,
      saveDataInteractor,
      listDataInteractor,
      logger,
    );

    return new Hapi(this.settings, dataController, logger);
  }
}

export default HapiFactory;
