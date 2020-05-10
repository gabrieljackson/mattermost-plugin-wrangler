import React from 'react';

import {Action, Store} from 'redux';
import {PluginRegistry} from 'mattermost-webapp/plugins/registry';

import {getSettings} from './actions';
import reducer from './reducers';

import SetupUI from './components/setup_ui';
import MoveThreadModal from './components/move_thread_modal';
import MoveThreadDropdown from './components/move_thread_dropdown';

const setupUILater = (registry: PluginRegistry, store: Store<object, Action<object>>): () => Promise<void> => async () => {
    registry.registerReducer(reducer);

    const settings = await store.dispatch(getSettings());

    if (settings.data.enable_web_ui) {
        registry.registerRootComponent(MoveThreadModal);
        registry.registerPostDropdownMenuComponent(MoveThreadDropdown);
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
