import http from 'k6/http';
import { check, sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
    stages: [
        { duration: '30s', target: 20 },  // Ramp up to 20
        { duration: '1m', target: 100 }, // Push to 100
        { duration: '2m', target: 100 }, // Sustained heavy load
        { duration: '30s', target: 0 },   // Ramp down
    ],
    thresholds: {
        http_req_failed: ['rate<0.05'], // Allow 5% failure for high-concurrency lifecycle
        http_req_duration: ['p(95)<1000'], // P95 under 1s for full operations
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
    const uniqueId = uuidv4().substring(0, 8);
    const email = `user-${uniqueId}@loadtest.local`;
    const password = 'Password123!';

    const headers = { 'Content-Type': 'application/json' };

    // 1. REGISTER
    const regPayload = JSON.stringify({ email, password, name: `User ${uniqueId}` });
    const regRes = http.post(`${BASE_URL}/auth/register`, regPayload, { headers });
    check(regRes, { 'reg success': (r) => r.status === 201 || r.status === 200 });

    if (regRes.status !== 201 && regRes.status !== 200) {
        console.error(`Registration failed for ${email}: ${regRes.body}`);
        return;
    }

    // 2. LOGIN
    const loginRes = http.post(`${BASE_URL}/auth/login`, regPayload, { headers });
    const loginCheck = check(loginRes, { 'login success': (r) => r.status === 200 });
    if (!loginCheck) return;

    const apiKey = loginRes.json('data.api_key');
    if (!apiKey) {
        console.error(`API Key not found in login response: ${loginRes.body}`);
    }
    const authHeaders = { ...headers, 'X-API-Key': apiKey };

    // 3. CREATE VPC
    const vpcPayload = JSON.stringify({ name: `vpc-${uniqueId}`, cidr_block: '10.0.0.0/16' });
    const vpcRes = http.post(`${BASE_URL}/vpcs`, vpcPayload, { headers: authHeaders });
    check(vpcRes, { 'vpc created': (r) => r.status === 201 });
    if (vpcRes.status !== 201) return;
    const vpcId = vpcRes.json('data.id');

    // 4. CREATE SUBNET
    const subnetPayload = JSON.stringify({
        name: `subnet-${uniqueId}`,
        vpc_id: vpcId,
        cidr_block: '10.0.1.0/24'
    });
    const subnetRes = http.post(`${BASE_URL}/vpcs/${vpcId}/subnets`, subnetPayload, { headers: authHeaders });
    check(subnetRes, { 'subnet created': (r) => r.status === 201 });
    if (subnetRes.status !== 201) return;
    const subnetId = subnetRes.json('data.id');

    // 5. LAUNCH INSTANCE
    const instPayload = JSON.stringify({
        name: `inst-${uniqueId}`,
        image: 'alpine:latest',
        vpc_id: vpcId,
        subnet_id: subnetId,
        ports: '80:80'
    });
    const instRes = http.post(`${BASE_URL}/instances`, instPayload, { headers: authHeaders });
    check(instRes, { 'instance launch accepted': (r) => r.status === 202 });
    if (instRes.status !== 202) return;
    const instId = instRes.json('data.id');

    // 6. WAIT FOR RUNNING (Async Provisioning)
    let isRunning = false;
    for (let i = 0; i < 10; i++) {
        const getRes = http.get(`${BASE_URL}/instances/${instId}`, { headers: authHeaders });
        if (getRes.status === 200 && getRes.json('data.status') === 'running') {
            isRunning = true;
            break;
        }
        sleep(1);
    }
    check(isRunning, { 'instance is running': (val) => val === true });

    // 7. GET STATS (Simulate monitoring)
    if (isRunning) {
        const statsRes = http.get(`${BASE_URL}/instances/${instId}/stats`, { headers: authHeaders });
        check(statsRes, { 'stats retrieved': (r) => r.status === 200 });
    }

    // 7. CLEANUP (Delete in reverse)
    const delInstRes = http.del(`${BASE_URL}/instances/${instId}`, null, { headers: authHeaders });
    check(delInstRes, { 'inst deleted': (r) => r.status === 204 || r.status === 200 });

    const delVpcRes = http.del(`${BASE_URL}/vpcs/${vpcId}`, null, { headers: authHeaders });
    check(delVpcRes, { 'vpc deleted': (r) => r.status === 204 || r.status === 200 });

    sleep(1);
}
