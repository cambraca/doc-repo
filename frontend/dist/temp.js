import config from "./config.json" with {type: "json"};

function log(value) {
    document.getElementById("output").value += value + '\n';
}

async function run() {
    log('running...');
    // const response = await fetch('https://' + config.api_url + '/status');
    // log(await response.text());
}

log('api url is: ' + config.api_url);
run();
