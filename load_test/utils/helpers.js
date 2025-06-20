import {check, sleep} from 'k6';
import {Rate} from 'k6/metrics';

export const errorRate = new Rate('errors');

export function checkResponse(response, expectedStatus = 200, checkName = 'response check'){
    const result = check(response, {
        [`${checkName} - status is ${expectedStatus}`]: (r) => r.status === expectedStatus,
        [`${checkName} - response time < 500ms`]: (r) => r.timings.duration < 500,
        [`${checkName} - has valid JSON`]: (r) => {
            try {
                JSON.parse(r.body);
                return true;
            } catch (e) {
                return false;
            }
        }
    });

    errorRate.add(!result);
    return result;
}

export function randomSleep(min = 1, max = 3) {
    sleep(Math.random() * (max - min) + min);
}

export function logResponse(response, context = '') {
    console.log(`${context} - Status: ${response.status}, Duration: ${response.timings.duration}ms`);
    if (response.status >= 400) {
        console.log(`${context} - Error Body: ${response.body}`);
    }
}