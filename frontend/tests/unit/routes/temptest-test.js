import { module, test } from 'qunit';
import { setupTest } from 'docrepo/tests/helpers';

module('Unit | Route | temptest', function (hooks) {
  setupTest(hooks);

  test('it exists', function (assert) {
    let route = this.owner.lookup('route:temptest');
    assert.ok(route);
  });
});
