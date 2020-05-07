import {GlobalState} from 'mattermost-redux/types/store';

import {RECEIVED_PLUGIN_SETTINGS, RECEIVED_CHANNELS_FOR_TEAM} from '../types/wrangler';
import {OPEN_MOVE_THREAD_MODAL, CLOSE_MOVE_THREAD_MODAL} from '../types/ui';

import Client from '../client';
import {GetChannelsForTeamRequest} from '../types/api';

export type GetStateFunc = () => GlobalState;
export type ActionResult = {
    data: any; //eslint-disable-line @typescript-eslint/no-explicit-any
} | {
    error: any; //eslint-disable-line @typescript-eslint/no-explicit-any
};
export type DispatchFunc = (action: Action, getState?: GetStateFunc | null) => Promise<ActionResult>;
export type ActionFunc = (dispatch: DispatchFunc, getState: GetStateFunc) => Promise<ActionResult|ActionResult[]> | ActionResult;
export type Action = ActionFunc | GenericAction;
export type GenericAction = {
    type: string;
    [extraProps: string]: any; //eslint-disable-line @typescript-eslint/no-explicit-any
};

export function openMoveThreadModal(postID: string): ActionFunc {
    return async (dispatch: DispatchFunc) => {
        dispatch({
            type: OPEN_MOVE_THREAD_MODAL,
            post_id: postID,
        });

        return {data: postID};
    };
}

export function closeMoveThreadModal(): ActionFunc {
    return async (dispatch: DispatchFunc) => {
        dispatch({
            type: CLOSE_MOVE_THREAD_MODAL,
        });

        return {data: null};
    };
}

export function getSettings(): ActionFunc {
    return async (dispatch: DispatchFunc) => {
        const {data: settings, error} = await Client.getSettings();
        if (error) {
            return {error};
        }

        dispatch({
            type: RECEIVED_PLUGIN_SETTINGS,
            settings,
        });

        return {data: settings};
    };
}

export function getChannelsForTeam(teamID: string): ActionFunc {
    return async (dispatch: DispatchFunc) => {
        const {data: channels, error} = await Client.getChannelsForTeam({
            team_id: teamID,
        } as GetChannelsForTeamRequest);
        if (error) {
            return {error};
        }

        dispatch({
            type: RECEIVED_CHANNELS_FOR_TEAM,
            channels,
        });

        return {data: channels};
    };
}

export function moveThread(postID: string, threadID: string): ActionFunc {
    return async (dispatch: DispatchFunc, getState: GetStateFunc) => {
        const command = `/wrangler move thread ${postID} ${threadID}`;
        await Client.clientExecuteCommand(getState, command);

        return {data: null};
    };
}
