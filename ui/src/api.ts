import { ITile } from 'contracts';

const endPoint = "http://5.200.55.80:8000/data.json"

class ResponseError extends Error {
    response: any
}

export function getData<T>(): Promise<T> {
    return fetch(
        endPoint,
        {
            mode: 'no-cors'
        },
    ).then(checkStatus)
    .then(respToJson);
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