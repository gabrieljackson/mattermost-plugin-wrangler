import {Post} from 'mattermost-redux/types/posts';
import {UserProfile} from 'mattermost-redux/types/users';
import {Channel} from 'mattermost-redux/types/channels';

export type RichPost = {
    post: Post;
    user: UserProfile;
    channel: Channel;
}
