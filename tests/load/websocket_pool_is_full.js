import http from 'k6/http';
import ws from 'k6/ws';
import {check} from 'k6';
import {Counter} from 'k6/metrics';

export const options = {
    scenarios: {
        websocket_full_pool: {
            executor: 'shared-iterations',
            exec: 'websocket_full_pool',
            vus: 1500,
            iterations: 1500,
        },
    },
    rps: 40,
}

const poolIsNotFullCounter = new Counter('pool_is_not_full');
const poolIsFullCounter = new Counter('pool_is_full');
const unexpectedResponseCounter = new Counter('unexpected_response');

export function websocket_full_pool() {
    // Login
    const loginUrl = `http://${__ENV.BACKEND_HOSTNAME}/v0/auth/login`;
    const loginHeaders = {
        'Content-Type': 'application/json',
    }
    const loginBody = JSON.stringify({
        login: `test_user_${__VU}`,
        password: 'test'
    });
    const loginRes = http.post(loginUrl, loginBody, {
        headers: loginHeaders
    });
    check(loginRes, {
        'status is 200': (r) => r && r.status === 200,
        'response has token': (r) => r && r.body && r.json().accessToken,
    });
    const accessToken = loginRes.json().accessToken;

    // Connect websocket
    const wsUrl = `ws://${__ENV.BACKEND_HOSTNAME}/ws`;
    const wsParams = {
        headers: {'Authorization': `Bearer ${accessToken}`},
    }
    const wsRes = ws.connect(wsUrl, wsParams, (socket) => {
        socket.on('open', () => {
            console.log(`connected ${__VU}`);
            socket.setInterval(() => {
                socket.ping();
            }, 10000);
        });
        socket.on('error', (e) => {
            if (e.error() !== 'websocket: close sent') {
                console.log(`An unexpected error occurred ${__VU}: `, e.error());
            }
        });
        socket.on('close', () => console.log(`disconnected ${__VU}`));

        socket.setTimeout(function () {
            socket.close();
        }, 120000);
    });


    check(wsRes, {
        'status is 101 or 503': (r) => r && (r.status === 101 || r.status === 503),
    });

    if (wsRes.status === 101) {
        poolIsNotFullCounter.add(1)
    } else if (wsRes.status === 503) {
        poolIsFullCounter.add(1)
    } else {
        unexpectedResponseCounter.add(1)
    }
}
