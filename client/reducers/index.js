import { combineReducers } from 'redux';
import { routerReducer } from 'react-router-redux';
import settings from './settings';

const entities = (state = {}, action) => {
    if (action.error || !action.payload || !action.payload.entities) {
        return state;
    }
    return Object.assign({}, state, action.payload.entities);
};

export default combineReducers({entities, settings, routing: routerReducer});