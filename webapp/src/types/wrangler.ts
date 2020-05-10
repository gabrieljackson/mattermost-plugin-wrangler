import {Channel} from 'mattermost-redux/types/channels';

import id from '../plugin_id';

export type Settings = {
    enable_web_ui: boolean;
}

export type Channels = Array<Channel>

export const RECEIVED_PLUGIN_SETTINGS = `${id}_plugin_settings`;

export type ReceivedPluginSettingsAction = {
    type: typeof RECEIVED_PLUGIN_SETTINGS;
    settings: Settings;
};

export const RECEIVED_CHANNELS_FOR_TEAM = `${id}_received_channels_for_team`;

export type ReceivedChannelsForTeamAction = {
    type: typeof RECEIVED_CHANNELS_FOR_TEAM;
    channels: Channels;
};

export type WranglerActionType = ReceivedPluginSettingsAction | ReceivedChannelsForTeamAction;
