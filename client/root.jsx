import React from 'react';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';
import promise from 'redux-promise';
import createLogger from 'redux-logger';
import { Router, Route, hashHistory } from 'react-router';
import { syncHistoryWithStore, routerMiddleware } from 'react-router-redux';
import reducers from './reducers';
import Application from './app.jsx';
import Home from './containers/Home.jsx';
import Settings from './containers/Settings.jsx';
import ApiMiddleware from './middelware/palantir-api';

const Root = () => {
    const logger = createLogger();
    const store = createStore(
        reducers,
        applyMiddleware(ApiMiddleware, thunk, promise, routerMiddleware(hashHistory), logger)
    );

    const history = syncHistoryWithStore(hashHistory, store);
    
    return (
        <Provider store={store}>
            <Router history={history}>
                <Route path="/" component={Application}>
                    <Route path="" component={Home} />
                    <Route path="/settings" component={Settings} />
                </Route>
            </Router>
        </Provider>
    );
};

export default Root;
