import { module, test } from 'qunit';
import { setupTest } from 'docrepo/tests/helpers';

module('Unit | Controller | temptest', function (hooks) {
  setupTest(hooks);

  // TODO: Replace this with your real tests.
  test('it exists', function (assert) {
    let controller = this.owner.lookup('controller:temptest');
    assert.ok(controller);
  });
});
