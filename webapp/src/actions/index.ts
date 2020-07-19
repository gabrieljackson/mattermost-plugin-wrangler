import {GlobalState} from 'mattermost-redux/types/store';

import {RECEIVED_PLUGIN_SETTINGS} from '../types/wrangler';
import {OPEN_MOVE_THREAD_MODAL, CLOSE_MOVE_THREAD_MODAL} from '../types/ui';

import Client from '../client';

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

export function moveThread(postID: string, threadID: string, showRootMessage: boolean): ActionFunc {
    return async (dispatch: DispatchFunc, getState: GetStateFunc) => {
        const command = `/wrangler move thread ${postID} ${threadID} --show-root-message-in-summary=${showRootMessage}`;
        await Client.clientExecuteCommand(getState, command);

        return {data: null};
    };
}

export function copyThread(postID: string, threadID: string): ActionFunc {
    return async (dispatch: DispatchFunc, getState: GetStateFunc) => {
        const command = `/wrangler copy thread ${postID} ${threadID}`;
        await Client.clientExecuteCommand(getState, command);

        return {data: null};
    };
}
