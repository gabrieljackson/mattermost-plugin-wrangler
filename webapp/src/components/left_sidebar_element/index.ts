import {connect} from 'react-redux';

import {GlobalState} from 'mattermost-redux/types/store';
import {getNewSidebarPreference} from 'mattermost-redux/selectors/entities/preferences';

import LeftSidebarElement from './left_sidebar_element';

function mapStateToProps(state: GlobalState) {
    return {
        newSidebar: getNewSidebarPreference(state),
    };
}

export default connect(mapStateToProps, null)(LeftSidebarElement);
