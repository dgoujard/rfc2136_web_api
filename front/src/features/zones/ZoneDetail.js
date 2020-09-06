import {useDispatch, useSelector} from "react-redux";
import React, {Fragment, useEffect, useRef} from "react";
import {deleteRecordByZoneAndId, fetchZoneByName} from "./zonesSlices";
import { useNavigate } from 'react-router-dom';
import TableContainer from "@material-ui/core/TableContainer";
import Table from "@material-ui/core/Table";
import {Paper} from "@material-ui/core";
import TableCell from "@material-ui/core/TableCell";
import TableRow from "@material-ui/core/TableRow";
import TableHead from "@material-ui/core/TableHead";
import TableBody from "@material-ui/core/TableBody";
import Chip from "@material-ui/core/Chip";
import moment from "moment";
import Tooltip from "@material-ui/core/Tooltip";
import DeleteOutlinedIcon from '@material-ui/icons/DeleteOutlined';
import EditOutlinedIcon from '@material-ui/icons/EditOutlined';
import AddOutlinedIcon from '@material-ui/icons/AddOutlined';
import IconButton from "@material-ui/core/IconButton";
import Fab from "@material-ui/core/Fab";
import { makeStyles } from '@material-ui/core/styles';
import Collapse from "@material-ui/core/Collapse";
import Box from "@material-ui/core/Box";
import Typography from "@material-ui/core/Typography";
import CloseIcon from '@material-ui/icons/Close';
import RecordForm from "../../composent/recordform";

const scrollToRef = (ref) => window.scrollTo(0, ref.current.offsetTop)

const RelativeTTLTime = React.forwardRef(function RelativeTTLTime(props, ref) {
    //  Spread the props to the underlying DOM element.
    return <div {...props} ref={ref}>{props.ttl}</div>
});

const useStyles = makeStyles((theme) => ({
    fab: {
        position: 'fixed',
        bottom: theme.spacing(2),
        right: theme.spacing(2),
    }
}));

const useStylesRow = makeStyles((theme) => ({
    actionCol : {
        minWidth: '8rem'
    },
    actionBtn: {
        padding: '5px'
    }
}));

function RecordRow(props) {
    const {record, onDeleteRecordRequest, zoneName} = props;
    const [open, setOpen] = React.useState(false);
    const classes = useStylesRow();

    return (
        <React.Fragment>
            <TableRow>
                <TableCell component="th" scope="row">
                    {record.Fqdn}
                </TableCell>
                <TableCell align="right"><Chip variant="outlined" label={record.Type}/></TableCell>
                <TableCell align="right"><Tooltip arrow title={moment.duration(record.TTL, "seconds").humanize()}><RelativeTTLTime ttl={record.TTL}/></Tooltip></TableCell>
                <TableCell align="right">{record.Values.join(" ")}</TableCell>
                <TableCell align="right" className={classes.actionCol}>
                    <IconButton color="primary" className={classes.actionBtn} onClick={() => setOpen(!open)} aria-label="Edit record" component="span">
                        {!open ? <EditOutlinedIcon />:<CloseIcon/>}
                    </IconButton>
                    &nbsp;
                    <IconButton className={classes.actionBtn} onClick={()=>{onDeleteRecordRequest(record)}} color="secondary" aria-label="Delete record" component="span">
                        <DeleteOutlinedIcon />
                    </IconButton>
                </TableCell>
            </TableRow>
            <TableRow>
                <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6} >
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box margin={1}>
                            <Typography variant="h6" gutterBottom component="div">
                                Edit
                            </Typography>
                            <RecordForm zoneName={zoneName} id={record.GeneratedId} ttl={record.TTL} content={record.Values} fqdn={record.Fqdn} type={record.Type} />
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow>
        </React.Fragment>
    );
}
export function ZoneDetail() {
    const classes = useStyles();
    const dispatch = useDispatch();
    let navigate = useNavigate();
    const [addRecordOpen, setAddRecordOpen] = React.useState(false);
    const addRecordFormRef = useRef(null)
    const bottomTableFormRef = useRef(null)
    const executeScroll = () => scrollToRef(addRecordFormRef)
    const executeScrollBottom = () => scrollToRef(bottomTableFormRef)

    const { selectedZone:{name,records}, loading } =  useSelector(state => state.zones)
    useEffect(() => {
        if(typeof name == 'undefined' || name.length === 0)
            return;
        dispatch(fetchZoneByName(name)).then((e) => { //TODO Move the error handler to an composant in the App main component to avoid duplicate code
            console.log(e)
            // eslint-disable-next-line
            if(e.error !== undefined && (e.error.message === "Unauthorized" || e.error.message == 401 ) ){
                navigate('/login', { replace: true });
            }
        })
    },[dispatch,navigate,name])

    const onDeleteRecordRequest = record => {
        if(window.confirm("Do you want to delete "+record.Type+" "+record.Fqdn+" record ?")){
            dispatch(deleteRecordByZoneAndId({zoneName:name, recordId: record.GeneratedId})).then((e)=>{
                console.log(e);
            });
        }
    };

    return (
        <>
            {loading === "pending" &&
                <div>Loading...</div>
            }
            {loading !== "pending" && records !== undefined &&
                <Fragment >
                    <Typography variant="h5" ref={addRecordFormRef}>{name} records</Typography>
                    <Collapse in={addRecordOpen} timeout="auto" unmountOnExit >
                        <Box margin={1}>
                            <RecordForm onCreated={()=>{executeScrollBottom();setAddRecordOpen(!addRecordOpen)}} content={[]} zoneName={name}/>
                        </Box>
                    </Collapse>
                    <Fab color="primary" onClick={()=>{executeScroll();setAddRecordOpen(!addRecordOpen)}} aria-label="add" className={classes.fab}>
                        <AddOutlinedIcon />
                    </Fab>
                    <TableContainer component={Paper}>
                        <Table style={{minWidth: 650}} aria-label="simple table">
                            <TableHead>
                                <TableRow>
                                    <TableCell>Name</TableCell>
                                    <TableCell align="right">Type</TableCell>
                                    <TableCell align="right">TTL</TableCell>
                                    <TableCell align="right">Value</TableCell>
                                    <TableCell align="right">Action</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {records.map((data) => (
                                    <RecordRow key={data.GeneratedId} record={data} zoneName={name} onDeleteRecordRequest={onDeleteRecordRequest}  />
                                ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
                    <div ref={bottomTableFormRef}/>
                </Fragment>
            }
        </>
    )
}