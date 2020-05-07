export type MoveThreadRequest = {
    post_id: string;
    thread_id: string;
    original_channel_id: string;
}

export type GetChannelsForTeamRequest = {
    team_id: string;
}
