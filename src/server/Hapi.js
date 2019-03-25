import hapi from 'hapi';
import good from 'good';
import goodWinston from 'hapi-good-winston';

class HapiServer {
  constructor(settings, dataController, logger) {
    this.settings = settings;
    this.dataController = dataController;
    this.logger = logger;
  }

  async start() {
    const server = hapi.server({
      port: this.settings.server.port,
      router: {
        stripTrailingSlash: true,
      },
    });
    const goodOptions = {
      ops: false,
      reporters: {
        winston: [goodWinston(this.logger)],
      },
    };

    try {
      await server.register([
        {
          plugin: good,
          options: goodOptions,
        },
      ]);
      await server.route(this.createServerRoutes());
      await server.start();
      this.logger.info(`Listening on ${this.settings.server.port}`);
    } catch (err) {
      this.logger.error(err);
    }
  }

  createServerRoutes() {
    const routes = [
      {
        method: 'GET',
        path: '/healthcheck',
        handler: this.healthCheckHandler.bind(this),
      },
      {
        method: 'GET',
        path: '/data',
        handler: this.dataController.list.bind(this.dataController),
      },
      {
        method: 'GET',
        path: '/data/{id}',
        handler: this.dataController.listByDevice.bind(this.dataController),
      },
      {
        method: 'GET',
        path: '/data/{deviceId}/sensor/{sensorId}',
        handler: this.dataController.listBySensor.bind(this.dataController),
      },
      {
        method: 'POST',
        path: '/data',
        handler: this.dataController.save.bind(this.dataController),
      },
    ];

    return routes;
  }

  async healthCheckHandler(request, h) {
    return h.response({ online: true }).code(200);
  }
}

export default HapiServer;
