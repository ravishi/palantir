import React from 'react';
import ReactDOM from 'react-dom';
import { AppContainer } from 'react-hot-loader';
import injectTapEventPlugin from 'react-tap-event-plugin';
import RootComponent from './root.jsx';

injectTapEventPlugin();

const rootContainer = 'application';

ReactDOM.render(
    <AppContainer
        component={RootComponent} />,
    document.getElementById(rootContainer));

if (module.hot) {
    const _module = './root.jsx';
    module.hot.accept(_module, () => {
        ReactDOM.render(
            <AppContainer
                component={require(_module).default} />,
            document.getElementById(rootContainer));
    });
}
