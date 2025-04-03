export async function initialize(applicationInstance) {
  const url = applicationInstance.router.rootURL + 'config.json';
  const response = await fetch(url);
  const { api_url } = await response.json();
  applicationInstance.api_url = api_url;
}

export default {
  initialize,
};
