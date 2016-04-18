import { combineReducers } from 'redux';

const folders = (state = null, action) => {
    if (!action.error && action.type === 'FOLDER_LIST_RESPONSE') {
        return action.payload.result;
    }
    return state;
};

export default combineReducers({ folders });