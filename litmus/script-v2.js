import http from 'k6/http';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { sleep } from 'k6';

export const options = {
    vus: 100,
    duration: '60s',
};

export default function () {
    const url = 'http://order-app-v2.v2.svc.cluster.local:8080/order';
    const payload = JSON.stringify({
        order_id: uuidv4(),
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    http.post(url, payload, params);
    sleep(1);
}