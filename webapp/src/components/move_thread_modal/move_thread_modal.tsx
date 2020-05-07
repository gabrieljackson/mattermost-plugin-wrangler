import React from 'react';

import {Modal} from 'react-bootstrap';

import {getTeam, getTeamMemberships} from 'mattermost-redux/selectors/entities/teams';
import {Team} from 'mattermost-redux/types/teams';
import {Channel} from 'mattermost-redux/types/channels';
import Form from 'react-bootstrap/Form';

import {GlobalState} from 'mattermost-redux/types/store';

interface Props {
    visible: boolean;
    postID: string;
    message: string;
    moveThread: Function;
    getChannelsForTeam: Function;
    closeMoveThreadModal: Function;
    channelsForTeam: Array<Channel>;
    state: GlobalState;
}

type State = {
    allTeams: Array<Team>;
    channelsInTeam: Array<Channel>;
    selectedTeam: string;
    selectedChannel: string;
}

export default class MoveThreadModal extends React.PureComponent<Props, State> {
    constructor(props: Props) {
        super(props);

        this.state = {
            allTeams: Array<Team>(),
            channelsInTeam: Array<Channel>(),
            selectedTeam: '',
            selectedChannel: '',
        };
    }

    private loadTeams = async () => {
        const myTeamMemberships = getTeamMemberships(this.props.state);
        const myTeams = Array<Team>();
        Object.keys(myTeamMemberships).forEach((id) => {
            const team = getTeam(this.props.state, id);
            myTeams.push(team);
        });

        this.setState({allTeams: myTeams});
    }

    componentDidMount() {
        this.loadTeams();
    }

    private handleTeamSelectChange = async (event: React.ChangeEvent<HTMLInputElement> | React.ChangeEvent<HTMLSelectElement>) => {
        const teamID = event.target.value;
        const channels = await this.props.getChannelsForTeam(teamID);

        this.setState({selectedTeam: teamID, channelsInTeam: channels.data});
    }

    private handleChannelSelectChange = (event: React.ChangeEvent<HTMLInputElement> | React.ChangeEvent<HTMLSelectElement>) => {
        this.setState({selectedChannel: event.target.value});
    }

    private handleOnClick = async (event: React.MouseEvent) => {
        if (event && event.preventDefault) {
            event.preventDefault();
        }

        await this.props.moveThread(this.props.postID, this.state.selectedChannel);
        this.props.closeMoveThreadModal();
    };

    private handleClose = async (event: React.MouseEvent) => {
        if (event && event.preventDefault) {
            event.preventDefault();
        }

        this.props.closeMoveThreadModal();
    };

    public render() {
        let disabled = false;
        if (this.props.postID === '' || this.state.selectedChannel === '') {
            disabled = true;
        }

        const style = {
            resize: 'none',
        };

        return (
            <Modal
                dialogClassName='modal--scroll'
                show={this.props.visible}
                onHide={this.handleClose}
                onExited={this.handleClose}
                bsSize='large'
                backdrop='static'
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title>
                        {'Wrangler - Move Message to Another Thread'}
                    </Modal.Title>
                </Modal.Header>
                <Modal.Body
                    ref='modalBody'
                >
                    <Form>
                        <Form.Group controlId='exampleForm.ControlSelect1'>
                            <Form.Label>{'Team'}</Form.Label>
                            <Form.Control
                                defaultValue='Select a Team'
                                as='select'
                                onChange={this.handleTeamSelectChange}
                            >
                                <option
                                    id='team-select'
                                    value='Select a Team'
                                    disabled={true}
                                >
                                    {'Select a Team'}
                                </option>
                                {this.state.allTeams.map((team) => (
                                    <option
                                        key={team.id}
                                        id={team.id}
                                        value={team.id}
                                    >
                                        {team.display_name}
                                    </option>
                                ))}
                            </Form.Control>
                        </Form.Group>
                        <Form.Group controlId='exampleForm.ControlSelect2'>
                            <Form.Label>{'Channel'}</Form.Label>
                            <Form.Control
                                defaultValue='Select a Channel'
                                disabled={this.state.selectedTeam === ''}
                                as='select'
                                onChange={this.handleChannelSelectChange}
                            >
                                <option
                                    id='channel-select'
                                    value='Select a Channel'
                                    disabled={true}
                                >
                                    {'Select a Channel'}
                                </option>
                                {this.state.channelsInTeam.map((channel) => (
                                    <option
                                        key={channel.id}
                                        id={channel.id}
                                        value={channel.id}
                                    >
                                        {channel.display_name}
                                    </option>
                                ))}
                            </Form.Control>
                        </Form.Group>
                        <Form.Group controlId='exampleForm.ControlTextarea1'>
                            <Form.Label>{'Thread Root Message'}</Form.Label>
                            <textarea
                                style={style}
                                className='form-control'
                                rows={5}
                                value={this.props.message}
                                disabled={true}
                                readOnly={true}
                            />
                        </Form.Group>
                    </Form>
                </Modal.Body>
                <Modal.Footer>
                    <button
                        id='footerClose'
                        className='btn btn-link'
                        onClick={this.handleClose}
                    >
                        {'Close'}
                    </button>
                    <button
                        id='saveSetting'
                        className='btn btn-primary'
                        onClick={this.handleOnClick}
                        disabled={disabled}
                    >
                        {'Move Thread'}
                    </button>
                </Modal.Footer>
            </Modal>
        );
    }
}
