import ApplicationSerializer from './application';

export default class DocumentSerializer extends ApplicationSerializer {
  attrs = {
    createdBy: { embedded: 'always' }
  };
}
