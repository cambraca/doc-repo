import Model, { attr, belongsTo } from '@ember-data/model';

export default class UserModel extends Model {
  @attr name;
  @attr email;

  @belongsTo('account', {async: false, inverse: null}) account;
}
