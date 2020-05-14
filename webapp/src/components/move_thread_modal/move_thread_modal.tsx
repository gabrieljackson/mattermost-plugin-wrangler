import React from 'react';

import {Modal} from 'react-bootstrap';
import Form from 'react-bootstrap/Form';

import {Team} from 'mattermost-redux/types/teams';
import {Channel} from 'mattermost-redux/types/channels';

interface Props {
    visible: boolean;
    postID: string;
    message: string;
    moveThread: Function;
    getMyTeams: Function;
    getChannelsForTeam: Function;
    closeMoveThreadModal: Function;
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
        const myTeams = this.props.getMyTeams();

        let firstTeamID = '';
        let firstChannelID = '';
        let channels = Array<Channel>();
        if (myTeams.length > 0) {
            const firstTeam = myTeams[0];
            firstTeamID = firstTeam.id;
            const channelResponse = await this.props.getChannelsForTeam(firstTeamID);
            channels = channelResponse.data;
            if (channels.length > 0) {
                const firstChannel = channels[0];
                firstChannelID = firstChannel.id;
            }
        }

        this.setState({
            allTeams: myTeams,
            channelsInTeam: channels,
            selectedTeam: firstTeamID,
            selectedChannel: firstChannelID,
        });
    }

    componentDidMount() {
        this.loadTeams();
    }

    private handleTeamSelectChange = async (event: React.ChangeEvent<HTMLInputElement> | React.ChangeEvent<HTMLSelectElement>) => {
        const teamID = event.target.value;
        const channelResponse = await this.props.getChannelsForTeam(teamID);
        const channels = channelResponse.data;
        let firstChannelID = '';
        if (channels.length > 0) {
            const firstChannel = channels[0];
            firstChannelID = firstChannel.id;
        }

        this.setState({
            selectedTeam: teamID,
            selectedChannel: firstChannelID,
            channelsInTeam: channels,
        });
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
                centered={true}
            >
                <Modal.Header closeButton={true}>
                    <Modal.Title>
                        {'Wrangler - Move Thread to Another Channel'}
                    </Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <Form>
                        <Form.Group>
                            <Form.Label>{'Team'}</Form.Label>
                            <Form.Control
                                as='select'
                                onChange={this.handleTeamSelectChange}
                                value={this.state.selectedTeam}
                            >
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
                        <Form.Group>
                            <Form.Label>{'Channel'}</Form.Label>
                            <Form.Control
                                as='select'
                                onChange={this.handleChannelSelectChange}
                                value={this.state.selectedChannel}
                                disabled={this.state.selectedTeam === ''}
                            >
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
                        <Form.Group>
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
