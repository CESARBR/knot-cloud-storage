import Hapi from 'server/Hapi';
import SaveDataInteractor from 'interactors/SaveData';
import ListDataInteractor from 'interactors/ListData';
import DataController from 'controllers/DataController';
import DataStore from 'data/DataStore';
import UuidAliasResolverFactory from 'network/UuidAliasResolverFactory';

class HapiFactory {
  constructor(settings, database) {
    this.settings = settings;
    this.database = database;
  }

  create() {
    const dataStore = new DataStore(this.database);
    const uuidAliasResolver = new UuidAliasResolverFactory(this.settings).create();
    const saveDataInteractor = new SaveDataInteractor(dataStore, uuidAliasResolver);
    const listDataInteractor = new ListDataInteractor(dataStore);
    const dataController = new DataController(saveDataInteractor, listDataInteractor);

    return new Hapi(this.settings, dataController);
  }
}

export default HapiFactory;
