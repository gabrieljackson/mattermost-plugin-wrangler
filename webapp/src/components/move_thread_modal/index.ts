import {connect} from 'react-redux';
import {Dispatch, Action, bindActionCreators} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';
import {getTeam, getTeamMemberships} from 'mattermost-redux/selectors/entities/teams';
import {Team} from 'mattermost-redux/types/teams';
import {getPost} from 'mattermost-redux/selectors/entities/posts';

import {isMoveModalVisable, getMoveThreadPostID} from '../../selectors';
import {closeMoveThreadModal, moveThread, getChannelsForTeam} from '../../actions';

import MoveThreadModal from './move_thread_modal';

function mapStateToProps(state: GlobalState) {
    let postID = getMoveThreadPostID(state);
    let post = getPost(state, postID);
    let message = '';
    let threadCount = 1;

    if (post) {
        if (post.root_id) {
            post = getPost(state, post.root_id);
        }

        postID = post.id;
        message = post.message;

        const postsInThread = state.entities.posts.postsInThread[postID];
        if (postsInThread) {
            threadCount = postsInThread.length + 1;
        }
    }

    const getMyTeamsFunc = () => {
        const myTeamMemberships = getTeamMemberships(state);
        const myTeams = Array<Team>();
        Object.keys(myTeamMemberships).forEach((id) => {
            const team = getTeam(state, id);
            myTeams.push(team);
        });

        return myTeams;
    };

    return {
        visible: isMoveModalVisable(state),
        getMyTeams: getMyTeamsFunc,
        postID,
        message,
        threadCount,
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
