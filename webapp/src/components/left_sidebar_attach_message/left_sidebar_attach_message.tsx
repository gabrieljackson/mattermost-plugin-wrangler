import React from 'react';

import {RichPost} from 'src/types/attach';

import LeftSidebarElement from '../left_sidebar_element';

interface Props {
    finishAttachingPost: Function;
    postToBeAttached: RichPost;
}

type State = {}

export default class LeftSidebarAttachMessage extends React.PureComponent<Props, State> {
    private exit = async (event: React.MouseEvent) => {
        if (event && event.preventDefault) {
            event.preventDefault();
        }

        this.props.finishAttachingPost();
    };

    public render() {
        const postToBeAttached = this.props.postToBeAttached;
        if (!postToBeAttached) {
            return null;
        }

        const name = postToBeAttached.user.first_name ? postToBeAttached.user.first_name : postToBeAttached.user.username;
        const originalMessage = postToBeAttached.post.message;
        const trimmed = originalMessage.length > length ? originalMessage.substring(0, 75) + '...' : originalMessage;
        const tooltipContent = (<div>
            <p>{'Howdy Partner!'}</p>
            <p>{'It looks like you are attaching a message to a thread.'}</p>
            <p>{'Use the post dropdown on the thread you want to attach it to or click the "X" right here to quit.'}</p>
            <hr/>
            <p>{name + '\'s message in ' + postToBeAttached.channel.display_name + ':'}</p>
            <p>{'"' + trimmed + '"'}</p>
        </div>);

        return (
            <LeftSidebarElement
                id={'attach-message'}
                text={'Attaching Message'}
                tooltip={tooltipContent}
                clickHandler={this.exit}
            />
        );
    }
}
