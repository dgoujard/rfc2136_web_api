import React from 'react';
import Container from '@material-ui/core/Container';
import {Zones} from "../features/zones/Zones";
import Typography from "@material-ui/core/Typography";
// Modules


const ZoneListPage = () => {
    return (
        <Container maxWidth="sm">
            <Typography variant="h5">DNS zones list</Typography>
            <Zones />
        </Container>
    );
};

export default ZoneListPage;