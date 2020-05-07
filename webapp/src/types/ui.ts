import id from '../plugin_id';

export const OPEN_MOVE_THREAD_MODAL = `${id}_open_move_thread_modal`;

export type OpenMoveThreadAction = {
    type: typeof OPEN_MOVE_THREAD_MODAL;
    post_id: string;
};

export const CLOSE_MOVE_THREAD_MODAL = `${id}_close_move_thread_modal`;

export type CloseMoveThreadAction = {
    type: typeof CLOSE_MOVE_THREAD_MODAL;
};

export type UIActionType = OpenMoveThreadAction | CloseMoveThreadAction;
