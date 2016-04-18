import React from 'react';
import { connect } from 'react-redux'
import MuiAppBar from 'material-ui/AppBar';
import IconButton from 'material-ui/IconButton';
import IconMenu from 'material-ui/IconMenu';
import MenuItem from 'material-ui/MenuItem';
import MoreVertIcon from 'material-ui/svg-icons/navigation/more-vert';
import { goToSettings, goToHome } from '../actions';

const customAppBar = (props) => (
    <MuiAppBar iconElementRight={
            <IconMenu
                iconButtonElement={<IconButton><MoreVertIcon /></IconButton>}
                targetOrigin={{horizontal: 'right', vertical: 'top'}}
                anchorOrigin={{horizontal: 'right', vertical: 'top'}}>
                <MenuItem primaryText="Settings" onTouchTap={props.onSettingsTouchTap} />
            </IconMenu>
        } {...props} onLeftIconButtonTouchTap={props.onHomeTouchTap} />
);

const mapStateToProps = (state, ownProps) => {
    return {
        title: state.windowTitle || ownProps.title
    };
};

const mapDispatchToProps = (dispatch, ownProps) => {
    return {
        onHomeTouchTap: () => dispatch(goToHome()),
        onSettingsTouchTap: () => dispatch(goToSettings())
    };
};

const AppBar = connect(
    mapStateToProps,
    mapDispatchToProps
)(customAppBar);

export default AppBar;
