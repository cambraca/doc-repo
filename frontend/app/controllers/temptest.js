import Controller from '@ember/controller';
import { action } from '@ember/object';
import { service } from '@ember/service';

export default class TemptestController extends Controller {
	@service store;

	docs = [];

	@action
	async loadDocs() {
		this.set('docs', await this.store.findAll('document'));
	}
}
