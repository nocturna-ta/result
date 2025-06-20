export const config = {
    base_url: __ENV.BASE_URL || 'http://localhost:8903',

    threshold: {
        http_req_duration: ['p(95)<500'],
        http_req_failed: ['rate<0.01'],
        http_reqs: ['rate>100'],
    },

    stages:{
        smoke: [
            { duration: '1m', target: 10 },
        ],
        load: [
            {
                duration: '2m',
                target: 10
            },
            {
                duration: '5m',
                target: 10
            },
            {
                duration: '2m',
                target: 0
            }
        ],
        endurance : [
            {
                duration: '2m',
                target: 20
            },
            {
                duration: '30m',
                target: 20
            },
            {
                duration: '2m',
                target: 0
            },

        ]
    }
}