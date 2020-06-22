export const MessageActionTypeMove = 'move';
export const MessageActionTypeCopy = 'copy';

export type MessageActionType = typeof MessageActionTypeMove | typeof MessageActionTypeCopy;
