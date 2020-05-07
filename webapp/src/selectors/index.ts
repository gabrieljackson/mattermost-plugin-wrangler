import {GlobalState} from 'mattermost-redux/types/store';

import id from '../plugin_id';

const pluginState = (state: GlobalState) => state['plugins-' + id] || {};

export const getPluginSettings = (state: GlobalState) => pluginState(state).pluginSettings;

export const isMoveModalVisable = (state: GlobalState) => pluginState(state).moveThreadModalVisable;

export const getMoveThreadPostID = (state: GlobalState) => pluginState(state).getMoveThreadPostID;

export const getChannelsForTeamSel = (state: GlobalState) => pluginState(state).getChannelsForTeam;
