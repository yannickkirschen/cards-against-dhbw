import { Alert, Snackbar } from "@mui/material"

/**
 * Component that provides an error-snackbar
 * @param {Object} props the props of the ErrorSnackbar, containing the message, the open-state and a method to set open
 * @returns an error-snackbar that shows for 4 seconds
 */
const ErrorSnackbar = ({ msg, open, setOpen }) => {
    return (
        <Snackbar
            anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
            open={open}
            onClose={() => setOpen(false)}
            autoHideDuration={4000}
            key={"this is a unique key"}
        >
            <Alert severity="error" sx={{ width: "100%" }}>
                {msg}
            </Alert>
        </Snackbar>
    )
}

export default ErrorSnackbar;
