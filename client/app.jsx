import React from 'react';
import getMuiTheme from 'material-ui/styles/getMuiTheme';
import AppBar from './containers/AppBar.jsx';

class Application extends React.Component {
    getChildContext() {
        return {muiTheme: getMuiTheme()};
    }

    render() {
        return (
            <div style={{margin: '16px 8px'}}>
                <AppBar title="Palantir" />
                {this.props.children}
            </div>
        );
    }
}

Application.childContextTypes = {
    muiTheme: React.PropTypes.object.isRequired
};

export default Application;