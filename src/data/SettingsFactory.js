import config from 'config';
import Joi from 'joi';
import _ from 'lodash';

import Settings from 'data/Settings';

const supportedDatabases = ['MONGO'];
const databaseSchema = Joi.string().valid(supportedDatabases).required();

const databaseConfigSchema = Joi.object().keys({
  hostname: Joi.string().required(),
  port: Joi.number().port().required(),
  name: Joi.string().required(),
});

const serverSchema = Joi.object().keys({
  port: Joi.number().port().required(),
});

class SettingsFactory {
  create() {
    const database = this.loadDatabaseSettings();
    const databaseConfig = this.loadDatabaseConfigSettings();
    const server = this.loadServerSettings();
    return new Settings(database, databaseConfig, server);
  }

  loadDatabaseSettings() {
    const database = config.get('database');
    this.validate('database', database, databaseSchema);
    return database;
  }

  loadDatabaseConfigSettings() {
    const databaseConfig = config.get('databaseConfig');
    this.validate('databaseConfig', databaseConfig, databaseConfigSchema);
    return databaseConfig;
  }

  loadServerSettings() {
    const server = config.get('server');
    this.validate('server', server, serverSchema);
    return server;
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
