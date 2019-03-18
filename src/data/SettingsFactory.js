import config from 'config';
import Joi from 'joi';
import _ from 'lodash';

import Settings from 'data/Settings';

const databaseSchema = Joi.object().keys({
  hostname: Joi.string(),
  port: Joi.number().port(),
  name: Joi.string(),
});

const serverSchema = Joi.object().keys({
  port: Joi.number().port().required(),
});

const meshbluSchema = Joi.object().keys({
  aliasLookupServerUri: Joi.string().required(),
  cacheRedisUri: Joi.string().required(),
});

class SettingsFactory {
  create() {
    const database = this.loadDatabaseSettings();
    const server = this.loadServerSettings();
    const meshblu = this.loadMeshbluSettings();
    return new Settings(database, server, meshblu);
  }

  loadDatabaseSettings() {
    const database = config.get('database');
    this.validate('database', database, databaseSchema);
    return database;
  }

  loadServerSettings() {
    const server = config.get('server');
    this.validate('server', server, serverSchema);
    return server;
  }

  loadMeshbluSettings() {
    const meshblu = config.get('meshblu');
    this.validate('meshblu', meshblu, meshbluSchema);
    return meshblu;
  }

  validate(propertyName, propertyValue, schema) {
    const { error } = Joi.validate(propertyValue, schema, { abortEarly: false });
    if (error) {
      throw this.mapJoiError(propertyName, error);
    }
  }

  mapJoiError(propertyName, error) {
    const reasons = _.map(error.details, 'message');
    const formattedReasons = reasons.length > 1
      ? `\n${_.chain(reasons).map(reason => `- ${reason}`).join('\n').value()}`
      : reasons[0];
    return new Error(`Invalid "${propertyName}" property: ${formattedReasons}`);
  }
}

export default SettingsFactory;
