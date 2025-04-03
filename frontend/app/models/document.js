import Model, { attr, belongsTo } from '@ember-data/model';

export default class DocumentModel extends Model {
  @attr name;
  @attr mimeType;
  @attr complexInfo;

  @belongsTo('user', {async: false, inverse: null}) createdBy;
}
