import PropTypes from "prop-types";
import React, {memo, useEffect} from "react";
import TextField from "@material-ui/core/TextField";
import Grid from "@material-ui/core/Grid";



const RecordFormContentInput = memo(({ register, unregister, setValue, name, type, defaultValue, watch, zoneName }) => {

    useEffect(() => {
        register({ name });
        setValue(name, defaultValue);
        return () => unregister(name);
    }, [name, register, unregister,setValue]);

    let currentValue = watch("content", defaultValue)
    let helpMessage = ""

    console.log(currentValue)
    switch (type){
        case "CNAME":
            if(currentValue.length > 0 && currentValue[0].length > 0 && currentValue[0][currentValue[0].length-1] !== "."){
                helpMessage = "Real target : « "+currentValue[0]+"."+zoneName+". » \n"
            }
        case "A":
        case "AAAA":
        case "NS":
        case "TXT":
            let inputLabel = "IPv4 address"
            if(type === "AAAA")
                inputLabel = "IPv6 address"
            else if(type === "CNAME" || type === "NS" )
                inputLabel = "Target"
            else if(type === "TXT")
                inputLabel = "Txt content"
            return (<Grid item xs={4}><TextField
                fullWidth
                required={true}
                helperText={helpMessage}
                label={inputLabel}
                name="content"
                size="small"
                variant="outlined"
                defaultValue={currentValue[0]}
                onChange={e => {setValue(name, [e.target.value]);}}
            /></Grid>)
        case "MX":
            const setMxServerValue = (e)=>{
                setValue(name, [currentValue[0], e.target.value])
            }
            const setMxPriorityValue = (e)=>{
                setValue(name, [e.target.value, currentValue[1]])
            }
            if(currentValue.length > 0 && currentValue[1].length > 0 && currentValue[1][currentValue[1].length-1] !== "."){
                helpMessage = "Real target : « "+currentValue[1]+"."+zoneName+". » \n"
            }
            return (<>
                <Grid item  xs={4}>
                <TextField
                    fullWidth
                    required={true}
                    label={"Mail server"}
                    name="mailserver"
                    helperText={helpMessage}
                    size="small"
                    variant="outlined"
                    defaultValue={currentValue[1]}
                    onChange={setMxServerValue}
                />
                </Grid>
                <Grid item  xs={1}>
                <TextField
                    required={true}
                    InputLabelProps={{ shrink: true }}
                    type="number"
                    label={"Priority"}
                    name="priority"
                    size="small"
                    variant="outlined"
                    helperText={"0 - 65535"}
                    defaultValue={currentValue[0]}
                    onChange={setMxPriorityValue}
                />
                </Grid>
            </>
            )
    }
    return <></>
});

RecordFormContentInput.propTypes = {
    type: PropTypes.string.isRequired,
    content: PropTypes.array,

}
export default RecordFormContentInput;