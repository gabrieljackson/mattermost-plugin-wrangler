import React from 'react';

import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faHatCowboy} from '@fortawesome/free-solid-svg-icons';

import {Tooltip, OverlayTrigger} from 'react-bootstrap';

import {RichPost} from 'src/types/attach';

import './style.scss';

interface Props {
    finishAttachingPost: Function;
    postToBeAttached: RichPost;
    newSidebar: boolean;
}

type State = {}

export default class LeftSidebarHeader extends React.PureComponent<Props, State> {
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

        const iconStyle = {
            margin: '0 7px 0 1px',
        };

        const buttonClass = this.props.newSidebar ? 'SidebarChannelGroupHeader_addButton pull-right' : 'btn-old-sidebar';

        const name = postToBeAttached.user.first_name ? postToBeAttached.user.first_name : postToBeAttached.user.username;
        const originalMessage = postToBeAttached.post.message;
        const trimmed = originalMessage.length > length ? originalMessage.substring(0, 75) + '...' : originalMessage;
        const ttContent = (<div>
            <p>{'Howdy Partner!'}</p>
            <p>{'It looks like you are attaching a message to a thread.'}</p>
            <p>{'Use the post dropdown on the thread you want to attach it to or click the "X" right here to quit.'}</p>
            <hr/>
            <p>{name + '\'s message in ' + postToBeAttached.channel.display_name + ':'}</p>
            <p>{'"' + trimmed + '"'}</p>
        </div>);

        return (
            <OverlayTrigger
                key='githubAssignmentsLink'
                placement='right'
                overlay={<Tooltip id='reviewTooltip'>{ttContent}</Tooltip>}
            >
                <div className={'wrangler-left-sidebar'}>
                    <FontAwesomeIcon
                        className='MenuItem__icon'
                        style={iconStyle}
                        icon={faHatCowboy}
                    />
                    {'Attaching Message'}
                    <button
                        type='button'
                        className={buttonClass}
                        aria-label='Exit Wrangler Mode'
                        onClick={this.exit}
                    >
                        <i className='icon-close'/>
                    </button>
                </div>
            </OverlayTrigger>
        );
    }
}
