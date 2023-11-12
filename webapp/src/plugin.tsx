import React from 'react';

import {Action, Store} from 'redux';
import {PluginRegistry} from 'mattermost-webapp/plugins/registry';

import {getChannel} from 'mattermost-redux/selectors/entities/channels';

import {getSettings, startCopyToChannel} from './actions';
import reducer from './reducers';

import SetupUI from './components/setup_ui';
import MoveThreadModal from './components/move_thread_modal';
import MoveThreadDropdown from './components/move_thread_dropdown';
import AttachMessageDropdown from './components/attach_message_dropdown';
import CopyToChannelDropdown from './components/copy_to_channel_dropdown';
import MergeThreadDropdown from './components/merge_thread_dropdown';
import LeftSidebarAttachMessage from './components/left_sidebar_attach_message';
import LeftSidebarCopyToChannel from './components/left_sidebar_copy_to_channel';
import LeftSidebarMergeThread from './components/left_sidebar_merge_thread';

const setupUILater = (registry: PluginRegistry, store: Store<object, Action<object>>): () => Promise<void> => async () => {
    registry.registerReducer(reducer);

    const settings = await store.dispatch(getSettings());

    if (settings.data.enable_web_ui) {
        registry.registerRootComponent(MoveThreadModal);
        registry.registerLeftSidebarHeaderComponent(LeftSidebarAttachMessage);
        registry.registerLeftSidebarHeaderComponent(LeftSidebarCopyToChannel);
        registry.registerPostDropdownMenuComponent(MoveThreadDropdown);
        registry.registerPostDropdownMenuComponent(AttachMessageDropdown);
        registry.registerPostDropdownMenuComponent(CopyToChannelDropdown);
        registry.registerChannelHeaderMenuAction(
            'Copy Messages to Channel',
            (channelId: string) => store.dispatch(startCopyToChannel(getChannel(store.getState(), channelId))),
        );
        if (settings.data.enable_merge_thread) {
            registry.registerLeftSidebarHeaderComponent(LeftSidebarMergeThread);
            registry.registerPostDropdownMenuComponent(MergeThreadDropdown);
        }
    }
};

export default class Plugin {
    private haveSetupUI = false;
    private setupUI?: () => Promise<void>;

    private finishedSetupUI = () => {
        this.haveSetupUI = true;
    };

    public async initialize(registry: PluginRegistry, store: Store<object, Action<object>>) {
        this.setupUI = setupUILater(registry, store);
        this.haveSetupUI = false;

        // Register the dummy component, which will call setupUI when it is activated (i.e., when the user logs in)
        registry.registerRootComponent(
            () => {
                return (
                    <SetupUI
                        setupUI={this.setupUI}
                        haveSetupUI={this.haveSetupUI}
                        finishedSetupUI={this.finishedSetupUI}
                    />
                );
            },
        );
    }
}
