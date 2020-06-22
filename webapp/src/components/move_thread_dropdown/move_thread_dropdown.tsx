import React from 'react';

import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faHatCowboy} from '@fortawesome/free-solid-svg-icons';

interface Props {
    postID: string;
    threadCount: number;
    isSystemMessage: boolean;
    rootPostID: string;
    needRootMessage: boolean;
    getPostThread: Function;
    openMoveThreadModal: Function;
}

type State = {}

export default class MoveThreadDropdown extends React.PureComponent<Props, State> {
    private handleOnClick = async (event: React.MouseEvent) => {
        if (event && event.preventDefault) {
            event.preventDefault();
        }

        this.props.openMoveThreadModal(this.props.rootPostID);
    };

    private getRootMessage() {
        this.props.getPostThread(this.props.rootPostID);
    }

    public render() {
        if (this.props.isSystemMessage) {
            return null;
        }

        if (this.props.needRootMessage) {
            this.getRootMessage();
            return null;
        }

        let content = 'Move/Copy Message';
        if (this.props.threadCount > 1) {
            content = 'Move/Copy Thread';
        }

        return (
            <React.Fragment>
                <li
                    className='MenuItem'
                    role='menuitem'
                >
                    <button
                        className='style--none'
                        role='presentation'
                        onClick={this.handleOnClick}
                    >
                        <FontAwesomeIcon
                            className='MenuItem__icon'
                            icon={faHatCowboy}
                        />
                        {content}
                    </button>
                </li>
            </React.Fragment>
        );
    }
}
