import id from '../plugin_id';

import {RichPost} from './post';

export const INITIALIZE_MERGE_THREAD = `${id}_init_merge_thread`;

export type MergeTheadInitializeAction = {
    type: typeof INITIALIZE_MERGE_THREAD;
    post: RichPost;
};

export const FINALIZE_MERGE_THREAD = `${id}_finalize_merge_thread`;

export type MergeTheadFinalizeAction = {
    type: typeof FINALIZE_MERGE_THREAD;
};

export type MergeAction = MergeTheadInitializeAction | MergeTheadFinalizeAction;
