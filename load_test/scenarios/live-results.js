import http from 'k6/http';
import { check, sleep } from 'k6';
import { config } from '../utils/config.js';
import { generateHeaders, getRandomRegion } from '../utils/data-generator.js';
import { checkResponse, randomSleep, logResponse } from '../utils/helpers.js';

// Define load test options using the 'load' stage from config.js
export const options = {
    stages: config.stages.load,
    thresholds: config.threshold,
};

// Main function executed by each virtual user
export default function () {
    const baseUrl = config.base_url;
    const headers = generateHeaders();

    // Sample election pair IDs (replace with actual IDs if available)
    const electionPairIds = ['c3834ab2-7735-44b3-a4bc-7509c6c37d17', 'd17af195-440c-4132-8169-98472fa48bb4', 'e031d6aa-4a8d-40df-bcd8-764d9a4ba5f7'];
    const randomElectionId = electionPairIds[Math.floor(Math.random() * electionPairIds.length)];

    // Test /v1/live/elections/{election_pair_id}
    const electionUrl = `${baseUrl}/v1/live/elections/${randomElectionId}`;
    let electionResponse = http.get(electionUrl, { headers });
    checkResponse(electionResponse, 200, 'Live Election Results');
    logResponse(electionResponse, 'Live Election Results');

    randomSleep();

    // Test /v1/live/cities/{city_name}
    const cityName = getRandomRegion();
    const cityUrl = `${baseUrl}/v1/live/cities/${cityName}`;
    let cityResponse = http.get(cityUrl, { headers });
    checkResponse(cityResponse, 200, 'Live City Results');
    logResponse(cityResponse, 'Live City Results');

    randomSleep();
}