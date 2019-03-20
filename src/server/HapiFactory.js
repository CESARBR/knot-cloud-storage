import Hapi from 'server/Hapi';
import SaveDataInteractor from 'interactors/SaveData';
import ListDataInteractor from 'interactors/ListData';
import VerifySignatureInteractor from 'interactors/VerifySignature';
import DataController from 'controllers/DataController';
import DataStore from 'data/DataStore';
import UuidAliasResolverFactory from 'network/UuidAliasResolverFactory';
import CloudFactory from 'network/CloudFactory';

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
    const verifySignatureInteractor = new VerifySignatureInteractor(this.settings.server.publicKey);
    const dataController = new DataController(saveDataInteractor, listDataInteractor,
      verifySignatureInteractor);

    return new Hapi(this.settings, dataController);
  }
}

export default HapiFactory;
