import {combineReducers} from 'redux';

import {RECEIVED_PLUGIN_SETTINGS, ReceivedPluginSettingsAction, RECEIVED_CHANNELS_FOR_TEAM, ReceivedChannelsForTeamAction, Settings, Channels} from '../types/wrangler';
import {OPEN_MOVE_THREAD_MODAL, CLOSE_MOVE_THREAD_MODAL, UIActionType, OpenMoveThreadAction} from '../types/ui';

function pluginSettings(state: Settings | null = null, action: ReceivedPluginSettingsAction) {
    switch (action.type) {
    case RECEIVED_PLUGIN_SETTINGS:
        return action.settings;
    default:
        return state;
    }
}

function moveThreadModalVisable(state = false, action: UIActionType) {
    switch (action.type) {
    case OPEN_MOVE_THREAD_MODAL:
        return true;
    case CLOSE_MOVE_THREAD_MODAL:
        return false;
    default:
        return state;
    }
}

function getMoveThreadPostID(state = '', action: OpenMoveThreadAction) {
    switch (action.type) {
    case OPEN_MOVE_THREAD_MODAL:
        return action.post_id;
    case CLOSE_MOVE_THREAD_MODAL:
        return '';
    default:
        return state;
    }
}

const rootReducer = combineReducers({
    pluginSettings,
    getMoveThreadPostID,
    moveThreadModalVisable,
});

export default rootReducer;

export type RootState = ReturnType<typeof rootReducer>;
