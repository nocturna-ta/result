import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';


export function generateHeaders(){
    return {
        'Content-Type': 'application/json',
        'X-User-Id': uuidv4(),
        'X-Address-Id': '0x' + Math.random().toString(16).substring(2, 40),
        'X-Role': 'voter',
    };
}

export function getRandomRegion(){
    const regions = [
        'Jakarta', 'Bandung', 'Surabaya', 'Medan', 'Semarang',
        'Makassar', 'Palembang', 'Tangerang', 'Depok', 'Bekasi'
    ];
    return regions[Math.floor(Math.random() * regions.length)];
}

export function getRandomVoteStatus(){
    const statuses = ['pending', 'confirmed', 'rejected', 'error', 'queued', 'retrying'];
    return statuses[Math.floor(Math.random() * statuses.length)];
}

export function generateRandomVoteId(){
    return uuidv4();
}

