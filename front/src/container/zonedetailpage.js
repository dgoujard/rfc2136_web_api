import React, {useEffect} from 'react';
import Container from '@material-ui/core/Container';
import {useDispatch} from "react-redux";
import { selectZone} from "../features/zones/zonesSlices";
import { useParams } from 'react-router-dom';
import {ZoneDetail} from "../features/zones/ZoneDetail";

// Modules


const ZoneDetailPage = () => {
    const dispatch = useDispatch();
    let { zoneName } = useParams();
    useEffect(() => {
        dispatch(selectZone({zoneName : zoneName}))
    },[dispatch, zoneName])

    return (
        <Container fixed>
            <ZoneDetail/>
        </Container>
    );
};

export default ZoneDetailPage;