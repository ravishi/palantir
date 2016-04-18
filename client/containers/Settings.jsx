import React from 'react';
import { connect } from 'react-redux'
import Paper from 'material-ui/Paper';
import CircularProgress from 'material-ui/CircularProgress';
import Subheader from 'material-ui/Subheader';
import { List, ListItem } from 'material-ui/List';
import { getFolders } from '../actions';

class FolderListComponent extends React.Component {
    componentDidMount() {
        this.props.loadFolders();
    }

    render() {
        const { ready, folders } = this.props;
        const renderContent = () => {
            if (!ready) {
                return <CircularProgress mode="indeterminate" size={2} />;
            } else {
                return folders.map(folder => <ListItem key={folder.id} primaryText={folder.path} />);
            }
        };
        return (
            <List>
                <Subheader>Folders</Subheader>
                {renderContent()}
            </List>
        );
    }
}

const mapStateToProps = (state, ownProps) => {
    const { folders } = state.settings;
    const { entities } = state;
    return {
        ready: Boolean(folders),
        folders: (folders || []).map(i => entities.folders[i])
    };
};

const mapDispatchToProps = (dispatch, ownProps) => {
    return {
        loadFolders: () => dispatch(getFolders())
    };
};

const FolderList = connect(mapStateToProps, mapDispatchToProps)(FolderListComponent);

const Settings = (props) => (
    <div style={{margin: '36px 0'}}>
        <Paper zDepth={1}>
            <FolderList />
        </Paper>
    </div>
);

export default Settings;