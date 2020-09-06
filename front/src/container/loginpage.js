import React from 'react';
import Button from '@material-ui/core/Button';
import Container from '@material-ui/core/Container';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import { makeStyles } from '@material-ui/core/styles';
import { useForm } from 'react-hook-form';
import { updateAuthCredentials } from "../features/zones/zonesSlices";
import {useDispatch, useSelector} from "react-redux";
import { Navigate } from 'react-router-dom';

// Modules

const useStyles = makeStyles((theme) => ({
    container: {
        padding: theme.spacing(3),
    },
}));

const LoginPage = () => {
    const { handleSubmit, register } = useForm();
    const dispatch = useDispatch();
    const { authIsValid } =  useSelector(state => state.zones)

    const classes = useStyles();

    const onSubmit = handleSubmit((data) => {
        dispatch(updateAuthCredentials({user: data.user, password: data.password}))
    });
    return (
        <Container className={classes.container} maxWidth="xs">
            {authIsValid && <Navigate to="/" replace={true} />}
            <form onSubmit={onSubmit}>
                <Grid container spacing={3}>
                    <Grid item xs={12}>
                        <Grid container spacing={2}>
                            <Grid item xs={12}>
                                <TextField
                                    fullWidth
                                    inputRef={register}
                                    label="Username"
                                    name="user"
                                    size="small"
                                    variant="outlined"
                                />
                            </Grid>
                            <Grid item xs={12}>
                                <TextField
                                    fullWidth
                                    inputRef={register}
                                    label="Password"
                                    name="password"
                                    size="small"
                                    type="password"
                                    variant="outlined"
                                />
                            </Grid>
                        </Grid>
                    </Grid>
                    <Grid item xs={12}>
                        <Button color="secondary" fullWidth type="submit" variant="contained">
                            Log in
                        </Button>
                    </Grid>
                </Grid>
            </form>
        </Container>
    );
};

export default LoginPage;