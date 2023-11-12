import {connect} from 'react-redux';
import {Dispatch, Action, bindActionCreators} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';
import {isCombinedUserActivityPost} from 'mattermost-redux/utils/post_list';
import {isSystemMessage} from 'mattermost-redux/utils/post_utils';
import {getPost} from 'mattermost-redux/selectors/entities/posts';
import {getUser} from 'mattermost-redux/selectors/entities/users';
import {getChannel} from 'mattermost-redux/selectors/entities/channels';

import {startMergingThread, finishMergingThread, mergeThread} from '../../actions';
import {getMergeThreadPost} from '../../selectors';

import MergeThreadDropdown from './merge_thread_dropdown';

interface Props {
    postId: string;
}

function mapStateToProps(state: GlobalState, props: Props) {
    const post = getPost(state, props.postId);
    const oldSystemMessageOrNull = post ? isSystemMessage(post) : true;
    const systemMessage = isCombinedUserActivityPost(post) || oldSystemMessageOrNull;

    let user = null;
    let channel = null;
    if (!systemMessage) {
        user = getUser(state, post.user_id);
        channel = getChannel(state, post.channel_id);
    }

    let validMerge = false;
    if (post) {
        if (state.entities.posts.postsInThread[post.id] && post.root_id === '') {
            validMerge = true;
        }
    }

    return {
        post,
        user,
        channel,
        isSystemMessage: systemMessage,
        isValidMergeMessage: validMerge,
        mergeThreadPost: getMergeThreadPost(state),
    };
}

function mapDispatchToProps(dispatch: Dispatch<Action>) {
    return bindActionCreators({
        startMergingThread,
        finishMergingThread,
        mergeThread,
    }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(MergeThreadDropdown);
