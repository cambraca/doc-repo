// import RESTSerializer from '@ember-data/serializer/rest';
import RESTSerializer, { EmbeddedRecordsMixin } from '@ember-data/serializer/rest';

export default class ApplicationSerializer extends RESTSerializer.extend(EmbeddedRecordsMixin) {
  // attrs = {
  //   createdBy: { embedded: 'always' }
  // };
}
