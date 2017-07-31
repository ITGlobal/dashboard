import { ITile } from 'contracts';

import * as data from './data.json'

declare var ENDPOINT;

class ResponseError extends Error {
    response: any
}

export function getData<T>(): Promise<T> {
    if (ENDPOINT) {
        return fetch(
            ENDPOINT,
            {
                mode: 'no-cors'
            },
        ).then(checkStatus)
            .then(respToJson);
    }

    return Promise.resolve(data as any);
}

const respToJson = response => response.json();
const checkStatus = response => {
    if (response.status >= 200 && response.status < 300) {
        return response
    } else {
        const error = new ResponseError(response.statusText)
        error.response = response
        throw error
    }
}
