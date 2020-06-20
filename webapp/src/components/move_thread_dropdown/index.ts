import {connect} from 'react-redux';
import {Dispatch, Action, bindActionCreators} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';
import {isCombinedUserActivityPost} from 'mattermost-redux/utils/post_list';
import {isSystemMessage} from 'mattermost-redux/utils/post_utils';
import {getPost as getPostSel} from 'mattermost-redux/selectors/entities/posts';
import {getPostThread} from 'mattermost-redux/actions/posts';

import {openMoveThreadModal} from '../../actions';

import MoveThreadDropdown from './move_thread_dropdown';

interface Props {
    postId: string;
}

function mapStateToProps(state: GlobalState, props: Props) {
    const post = getPostSel(state, props.postId);
    const oldSystemMessageOrNull = post ? isSystemMessage(post) : true;
    const systemMessage = isCombinedUserActivityPost(post) || oldSystemMessageOrNull;
    let needRootMessage = false;
    let rootPostID = props.postId;
    let threadCount = 1;

    if (post) {
        threadCount = post.reply_count + 1;
        if (post.root_id) {
            rootPostID = post.root_id;
            const rootPost = getPostSel(state, post.root_id);
            if (rootPost) {
                threadCount = rootPost.reply_count + 1;
            } else {
                needRootMessage = true;
            }
        }
    }

    return {
        postID: props.postId,
        isSystemMessage: systemMessage,
        threadCount,
        needRootMessage,
        rootPostID,
    };
}

function mapDispatchToProps(dispatch: Dispatch<Action>) {
    return bindActionCreators({
        openMoveThreadModal,
        getPostThread,
    }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(MoveThreadDropdown);
