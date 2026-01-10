import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '1m', target: 50 },  // Ramp up to 50 users
        { duration: '2m', target: 50 },  // Stay at 50 users
        { duration: '1m', target: 100 }, // Push to 100 users
        { duration: '2m', target: 100 }, // Stay at 100 users
        { duration: '1m', target: 0 },   // Ramp down
    ],
    thresholds: {
        http_req_failed: ['rate<0.02'],
        http_req_duration: ['p(99)<2000'], // Allow up to 2s for P99 at high load
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const API_KEY = __ENV.API_KEY || 'test-api-key';

export default function () {
    const params = {
        headers: {
            'X-API-Key': API_KEY,
            'Content-Type': 'application/json',
        },
    };

    // Simulate common dashboard user pattern
    const responses = http.batch([
        ['GET', `${BASE_URL}/api/dashboard/summary`, null, params],
        ['GET', `${BASE_URL}/instances`, null, params],
        ['GET', `${BASE_URL}/vpcs`, null, params],
    ]);

    check(responses[0], { 'summary success': (r) => r.status === 200 });
    check(responses[1], { 'instances success': (r) => r.status === 200 });

    // Every 5th request, create a temporary resource
    if (__ITER % 5 === 0) {
        const payload = JSON.stringify({
            name: `scale-test-${__VU}-${__ITER}`,
            image: 'alpine'
        });
        const res = http.post(`${BASE_URL}/instances`, payload, params);
        check(res, { 'inst created during scale': (r) => r.status === 201 });

        if (res.status === 201) {
            const instId = res.json('id');
            sleep(0.5);
            http.del(`${BASE_URL}/instances/${instId}`, null, params);
        }
    }

    sleep(1);
}
