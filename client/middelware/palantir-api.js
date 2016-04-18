import { Schema, arrayOf, normalize } from 'normalizr';
import { camelizeKeys } from 'humps';

export const API_CALL = 'API_CALL';
export const API_ROOT = 'http://localhost:5000';

const folderSchema = new Schema('folders');

export const Schemas = {
    FOLDER: folderSchema,
    FOLDER_ARRAY: arrayOf(folderSchema)
};

function fetchApi(endpoint) {
    const url = endpoint.startsWith(API_ROOT) ? endpoint : API_ROOT.concat(endpoint);

    const request = {
        headers: {
            'Accepts': 'application/json'
        }
    };

    return fetch(url, request)
        .then(response => (response.ok ? response : Promise.reject(response.statusText)))
        .then(response => response.json());
}

export function call({ endpoint, types, schema }) {
    return {
        type: API_CALL,
        payload: {
            types,
            schema,
            endpoint
        }
    };
}

export default store => next => action => {
    if (action.type !== API_CALL) {
        return next(action);
    }
    
    let { endpoint } = action.payload;
    const { schema, types } = action.payload;
    
    if (typeof endpoint === 'function') {
        endpoint = endpoint(store.getState());
    }
    
    if (typeof endpoint !== 'string') {
        throw new Error("`endpoint` should be either a string or a function which return a string.");
    }
    
    if (!Array.isArray(types) || types.length != 2) {
        throw new Error("`types` should be an array of two action types.")
    }
    
    if (!types.every(type => typeof type === 'string')) {
        throw new Error("`types` should be an array of strings.")
    }
    
    const [ requestType, responseType ] = types;
    
    next(derivedAction({type: requestType}));
    
    return fetchApi(endpoint)
        .then(camelizeKeys)
        .then(result => normalize(result, schema))
        .then(
            result => next(derivedAction({type: responseType, payload: result})),
            error => next(derivedAction({type: responseType, payload: error, error: true}))
    );
    
    function derivedAction({ type, error, payload }) {
        return { type, error, payload };
    }
};