import Hapi from 'server/Hapi';

class HapiFactory {
  constructor(settings) {
    this.settings = settings;
  }

  create() {
    return new Hapi(this.settings);
  }
}

export default HapiFactory;
