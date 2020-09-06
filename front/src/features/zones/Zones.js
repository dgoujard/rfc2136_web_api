import {useDispatch, useSelector} from "react-redux";
import React, {useEffect} from "react";
import {fetchZones} from "./zonesSlices";
import {Link, useNavigate} from 'react-router-dom';
import Button from "@material-ui/core/Button";

export function Zones() {
    const dispatch = useDispatch();
    let navigate = useNavigate();
    const { entities, loading } =  useSelector(state => state.zones)

    useEffect(() => {
        dispatch(fetchZones()).then((e) => { //TODO Move the error handler to an composant in the App main component to avoid duplicate code
            console.log(e)
            // eslint-disable-next-line
            if(e.error !== undefined && (e.error.message === "Unauthorized" || e.error.message == 401 ) ){
                navigate('/login', { replace: true });
            }
        })
    },[dispatch,navigate])

    return (
        <>
            {loading === "pending" ? <div>Loading...</div> : (entities.map(data => <Button key={data} component={Link} to={"/zones/"+data} variant="contained">{data}</Button>))}
        </>
    )
}