import {connect} from 'react-redux';
import {Dispatch, Action, bindActionCreators} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';
import {getPost} from 'mattermost-redux/selectors/entities/posts';

import {isMoveModalVisable, getMoveThreadPostID, getChannelsForTeamSel} from '../../selectors';
import {closeMoveThreadModal, moveThread, getChannelsForTeam} from '../../actions';

import MoveThreadModal from './move_thread_modal';

function mapStateToProps(state: GlobalState) {
    let postID = getMoveThreadPostID(state);
    let post = getPost(state, postID);
    const channels = getChannelsForTeamSel(state);
    let message = '';

    if (post) {
        if (post.root_id) {
            post = getPost(state, post.root_id);
        }

        postID = post.id;
        message = post.message;
    }

    return {
        visible: isMoveModalVisable(state),
        postID,
        message,
        channelsForTeam: channels,
        state,
    };
}

function mapDispatchToProps(dispatch: Dispatch<Action>) {
    return bindActionCreators({
        closeMoveThreadModal,
        moveThread,
        getChannelsForTeam,
    }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(MoveThreadModal);
