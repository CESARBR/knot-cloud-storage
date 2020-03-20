import Cloud from 'network/Cloud';

class CloudFactory {
  constructor(cloudRequester, uuidAliasResolver) {
    this.cloudRequester = cloudRequester;
    this.uuidAliasResolver = uuidAliasResolver;
  }

  create() {
    return new Cloud(this.cloudRequester, this.uuidAliasResolver);
  }
}

export default CloudFactory;
