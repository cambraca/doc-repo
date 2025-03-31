import config from "./config.json" with {type: "json"};

function log(value) {
    document.getElementById("output").value += value + '\n';
}

async function run() {
    log('running...');
    log('build time: ' + config.build_time);
    log('calling api: ' + config.api_url);
    const response = await fetch(config.api_url + '/status');
    log(await response.text());
    log('done!');
}

run();
