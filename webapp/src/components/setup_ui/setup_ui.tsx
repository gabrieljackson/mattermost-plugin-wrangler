import React from 'react';

interface Props {
    haveSetupUI: boolean;
    setupUI?: () => Promise<void>;
    finishedSetupUI: Function;
}

type State = {}

// SetupUI is a dummy Root component that we use to detect when the user has logged in
export default class SetupUI extends React.PureComponent<Props, State> {
    componentDidMount() {
        if (!this.props.haveSetupUI) {
            if (this.props.setupUI) {
                this.props.setupUI();
            }
            this.props.finishedSetupUI();
        }
    }

    render() {
        return null;
    }
}
