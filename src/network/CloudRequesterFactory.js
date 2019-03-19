import { JobManagerRequester } from 'meshblu-core-job-manager';

class CloudRequesterFactory {
  constructor(settings) {
    this.settings = settings;
  }

  create() {
    return new JobManagerRequester({
      namespace: this.settings.meshblu.namespace,
      redisUri: this.settings.meshblu.redisUri,
      maxConnections: 1,
      jobTimeoutSeconds: this.settings.meshblu.jobTimeoutSeconds,
      jobLogSampleRate: this.settings.meshblu.jobLogSampleRate,
      requestQueueName: this.settings.meshblu.requestQueueName,
      responseQueueName: this.settings.meshblu.responseQueueName,
      queueTimeoutSeconds: this.settings.meshblu.jobTimeoutSeconds,
    });
  }
}

export default CloudRequesterFactory;
