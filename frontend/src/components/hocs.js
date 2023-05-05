import { useNavigate } from "react-router-dom";

function withHistory(Component) {
    return props => <Component {...props} navigate={useNavigate()} />
}

export default withHistory;
