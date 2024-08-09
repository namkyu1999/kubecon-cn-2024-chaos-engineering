import http from 'k6/http';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
    discardResponseBodies: true,
    scenarios: {
        contacts: {
            executor: 'constant-arrival-rate',
            duration: '40s',
            rate: 30,
            timeUnit: '1s',
            preAllocatedVUs: 2,
            maxVUs: 50,
        },
    },
};

export default function () {
    const url = 'http://order-app.v2.svc.cluster.local:8080/order';
    const payload = JSON.stringify({
        order_id: uuidv4(),
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    http.post(url, payload, params);
}