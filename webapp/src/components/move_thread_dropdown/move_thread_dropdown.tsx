import React from 'react';

import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faHatCowboy} from '@fortawesome/free-solid-svg-icons';

interface Props {
    postId: string;
    isSystemMessage: boolean;
    openMoveThreadModal: Function;
}

type State = {}

export default class MoveThreadDropdown extends React.PureComponent<Props, State> {
    private handleOnClick = async (event: React.MouseEvent) => {
        if (event && event.preventDefault) {
            event.preventDefault();
        }

        this.props.openMoveThreadModal(this.props.postId);
    };

    public render() {
        if (this.props.isSystemMessage) {
            return null;
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
                        {'Wrangler - Move Thread'}
                    </button>
                </li>
            </React.Fragment>
        );
    }
}
