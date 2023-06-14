import { useNavigate } from "react-router-dom";

/**
 * alters a component to allow using navigation & history
 * @param {*} Component
 * @returns the component but with navigate.
 */
function withHistory(Component) {
    return props => <Component {...props} navigate={useNavigate()} />
}

export default withHistory;
