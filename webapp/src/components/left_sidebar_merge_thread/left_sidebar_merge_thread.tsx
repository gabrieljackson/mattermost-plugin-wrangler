import React from 'react';

import {RichPost} from 'src/types/post';

import LeftSidebarElement from '../left_sidebar_element';

interface Props {
    finishMergingThread: Function;
    mergeThreadPost: RichPost;
}

type State = {}

export default class LeftMergeThreadMessage extends React.PureComponent<Props, State> {
    private exit = async (event: React.MouseEvent) => {
        if (event && event.preventDefault) {
            event.preventDefault();
        }

        this.props.finishMergingThread();
    };

    public render() {
        const mergeThreadPost = this.props.mergeThreadPost;
        if (!mergeThreadPost) {
            return null;
        }

        const name = mergeThreadPost.user.first_name ? mergeThreadPost.user.first_name : mergeThreadPost.user.username;
        const originalMessage = mergeThreadPost.post.message;
        const trimmed = originalMessage.length > length ? originalMessage.substring(0, 75) + '...' : originalMessage;
        const tooltipContent = (<div>
            <p>{'Howdy Partner!'}</p>
            <p>{'It looks like you are merging a thread.'}</p>
            <p>{'Use the post dropdown on the thread you want to merge into or click the "X" right here to quit.'}</p>
            <hr/>
            <p>{name + '\'s message in ' + mergeThreadPost.channel.display_name + ':'}</p>
            <p>{'"' + trimmed + '"'}</p>
        </div>);

        return (
            <LeftSidebarElement
                id={'merge-thread'}
                text={'Merging Thread'}
                tooltip={tooltipContent}
                clickHandler={this.exit}
            />
        );
    }
}
