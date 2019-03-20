import httpSignature from 'http-signature';

function throwError(message, code) {
  const error = new Error(message);
  error.code = code;
  throw error;
}

class VerifySignature {
  constructor(publicKey) {
    this.publicKey = publicKey;
  }

  execute(request) {
    let parsedReq;
    try {
      parsedReq = httpSignature.parseRequest(request);
    } catch (error) {
      throwError('No authorization on request', 401);
    }
    if (!httpSignature.verifySignature(parsedReq, Buffer.from(this.publicKey, 'base64').toString('ascii'))) {
      throwError('Signature failed', 401);
    }
  }
}

export default VerifySignature;
