import {connect} from 'react-redux';
import {Dispatch, Action, bindActionCreators} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';
import {getNewSidebarPreference} from 'mattermost-redux/selectors/entities/preferences';

import {finishAttachingPost} from '../../actions';
import {getPostToBeAttached} from '../../selectors';

import LeftSidebarHeader from './left_sidebar_header';

function mapStateToProps(state: GlobalState) {
    return {
        postToBeAttached: getPostToBeAttached(state),
        newSidebar: getNewSidebarPreference(state),
    };
}

function mapDispatchToProps(dispatch: Dispatch<Action>) {
    return bindActionCreators({
        finishAttachingPost,
    }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(LeftSidebarHeader);
