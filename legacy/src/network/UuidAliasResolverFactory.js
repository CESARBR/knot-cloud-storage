import Redis from 'ioredis';
import RedisNS from '@octoblu/redis-ns';
import UuidAliasResolver from 'meshblu-uuid-alias-resolver';

class UuidAliasResolverFactory {
  constructor(settings) {
    this.settings = settings;
  }

  create() {
    const cacheClient = new Redis(this.settings.meshblu.cacheRedisUri, { dropBufferSupport: true });
    const uuidAliasClient = new RedisNS('uuid-alias', cacheClient);
    return new UuidAliasResolver({
      cache: uuidAliasClient,
      aliasServerUri: this.settings.meshblu.aliasLookupServerUri,
    });
  }
}

export default UuidAliasResolverFactory;
