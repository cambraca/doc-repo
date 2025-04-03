import RESTAdapter from '@ember-data/adapter/rest';

export default class ApplicationAdapter extends RESTAdapter {
  namespace = 'api/v1';

  constructor(config = {}) {
    super(config);
    if (config.api_url)
      this.host = config.api_url;
    else
      throw new Error('Missing API URL');
  }
}
