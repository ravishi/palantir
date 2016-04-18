import { push } from 'react-router-redux';
import { call, Schemas } from '../middelware/palantir-api';

export const goToHome = () => push('/');

export const goToSettings = () => push('/settings');

export const getFolders = () => {
    return call({
        types: ['FOLDER_LIST_REQUEST', 'FOLDER_LIST_RESPONSE'],
        schema: Schemas.FOLDER_ARRAY,
        endpoint: `/library/folders`
    });
};
