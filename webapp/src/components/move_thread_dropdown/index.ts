import {connect} from 'react-redux';
import {Dispatch, Action, bindActionCreators} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';
import {isCombinedUserActivityPost} from 'mattermost-redux/utils/post_list';
import {isSystemMessage} from 'mattermost-redux/utils/post_utils';
import {getPost} from 'mattermost-redux/selectors/entities/posts';

import {openMoveThreadModal} from '../../actions';

import MoveThreadDropdown from './move_thread_dropdown';

interface Props {
    postId: string;
}

function mapStateToProps(state: GlobalState, props: Props) {
    let post = getPost(state, props.postId);
    const oldSystemMessageOrNull = post ? isSystemMessage(post) : true;
    const systemMessage = isCombinedUserActivityPost(post) || oldSystemMessageOrNull;

    let threadCount = 1;
    if (post) {
        if (post.root_id) {
            post = getPost(state, post.root_id);
        }

        const postsInThread = state.entities.posts.postsInThread[post.id];
        if (postsInThread) {
            threadCount = postsInThread.length + 1;
        }
    }

    return {
        postId: props.postId,
        threadCount,
        isSystemMessage: systemMessage,
    };
}

function mapDispatchToProps(dispatch: Dispatch<Action>) {
    return bindActionCreators({
        openMoveThreadModal,
    }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(MoveThreadDropdown);
