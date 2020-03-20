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
  publicKey: Joi.string().base64().required(),
});

const meshbluSchema = Joi.object().keys({
  namespace: Joi.string().required(),
  messagesNamespace: Joi.string().required(),
  redisUri: Joi.string().uri({ scheme: 'redis' }).required(),
  cacheRedisUri: Joi.string().uri({ scheme: 'redis' }).required(),
  aliasLookupServerUri: Joi.string().uri().required(),
  jobTimeoutSeconds: Joi.number().positive().required(),
  jobLogSampleRate: Joi.number().integer().min(0).required(),
  requestQueueName: Joi.string().required(),
  responseQueueName: Joi.string().required(),
});

const levels = ['error', 'warn', 'info', 'verbose', 'debug', 'silly'];
const loggerSchema = Joi.object().keys({
  level: Joi.string().valid(levels).required(),
});

class SettingsFactory {
  create() {
    const database = this.loadDatabaseSettings();
    const databaseConfig = this.loadDatabaseConfigSettings();
    const server = this.loadServerSettings();
    const meshblu = this.loadMeshbluSettings();
    const logger = this.loadLoggerSettings();
    return new Settings(database, databaseConfig, server, meshblu, logger);
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

  loadMeshbluSettings() {
    const meshblu = config.get('meshblu');
    this.validate('meshblu', meshblu, meshbluSchema);
    return meshblu;
  }

  loadLoggerSettings() {
    const logger = config.get('logger');
    this.validate('logger', logger, loggerSchema);
    return logger;
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
