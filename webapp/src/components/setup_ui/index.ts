import {connect} from 'react-redux';
import {Dispatch, Action, bindActionCreators} from 'redux';

import SetupUI from './setup_ui';

function mapDispatchToProps(dispatch: Dispatch<Action>) {
    return bindActionCreators({}, dispatch);
}

export default connect(null, mapDispatchToProps)(SetupUI);
