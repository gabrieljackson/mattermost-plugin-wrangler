import {Post} from 'mattermost-redux/types/posts';
import {UserProfile} from 'mattermost-redux/types/users';
import {Channel} from 'mattermost-redux/types/channels';

import id from '../plugin_id';

export type RichPost = {
    post: Post;
    user: UserProfile;
    channel: Channel;
}

export const INITIALIZE_ATTACH_POST = `${id}_init_attach_post`;

export type AttachPostInitializeAction = {
    type: typeof INITIALIZE_ATTACH_POST;
    post: RichPost;
};

export const FINALIZE_ATTACH_POST = `${id}_finalize_attach_post`;

export type AttachPostFinalizeAction = {
    type: typeof FINALIZE_ATTACH_POST;
};

export type AttachAction = AttachPostInitializeAction | AttachPostFinalizeAction;
