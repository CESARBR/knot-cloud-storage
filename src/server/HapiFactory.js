import Hapi from 'server/Hapi';
import SaveDataInteractor from 'interactors/SaveData';
import DataController from 'controllers/DataController';
import DataStore from 'data/DataStore';

class HapiFactory {
  constructor(settings, database) {
    this.settings = settings;
    this.database = database;
  }

  create() {
    const dataStore = new DataStore(this.database);
    const saveDataInteractor = new SaveDataInteractor(dataStore);
    const dataController = new DataController(saveDataInteractor);

    return new Hapi(this.settings, dataController);
  }
}

export default HapiFactory;
