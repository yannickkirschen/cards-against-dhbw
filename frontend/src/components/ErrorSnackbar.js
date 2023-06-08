import { Alert, Snackbar } from "@mui/material"


const ErrorSnackbar = ({ msg, open, setOpen }) => {
    return (
        <Snackbar
            anchorOrigin={{ vertical: "bottom", horizontal: "center" }}
            open={open}
            onClose={() => setOpen(false)}
            autoHideDuration={3000}
            key={"this is a unique key"}
        >
            <Alert severity="error" sx={{ width: "100%" }}>
                {msg}
            </Alert>
        </Snackbar>
    )
}

export default ErrorSnackbar;
