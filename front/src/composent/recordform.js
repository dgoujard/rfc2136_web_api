import React from 'react';
import {Controller, useForm} from "react-hook-form";
import PropTypes from 'prop-types';
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";
import Select from "@material-ui/core/Select";
import MenuItem from "@material-ui/core/MenuItem";
import FormControl from "@material-ui/core/FormControl";
import FormHelperText from "@material-ui/core/FormHelperText";
import RecordFormContentInput from "./recordformcontentinput";
import InputAdornment from "@material-ui/core/InputAdornment";
import Grid from "@material-ui/core/Grid";
import {useDispatch} from "react-redux";
import {createRecordWithZoneNameAndRecord, updateRecordWithZoneNameAndRecord} from "../features/zones/zonesSlices";
// Modules

const RecordForm = (props) => {
    const { handleSubmit, register, unregister, setValue, control, errors, watch } = useForm();
    const { id, fqdn, type, ttl, content, zoneName, onCreated } = props;
    const dispatch = useDispatch();

    const onSubmit = handleSubmit((data) => {
        if(id !== undefined  && id.length > 0){
            dispatch(updateRecordWithZoneNameAndRecord({
                zoneName: zoneName,
                record: data,
                recordId: id
            })).then((e)=>{
                console.log(e);
                console.log("done")
            })
        }else{
            dispatch(createRecordWithZoneNameAndRecord({
                zoneName: zoneName,
                record: data
            })).then((e)=>{
                console.log(e);
                console.log("done")
                onCreated()
            })
        }

    });
    const defaultRecordType = "A"
    const selectedType = watch("type", (id !== undefined)?type:defaultRecordType);
    const defaultName = (fqdn !== undefined)?fqdn.substring( 0, fqdn.indexOf( "."+zoneName ) ):"";

    return (
        <form onSubmit={onSubmit} style={{flexGrow: 1}}>
            <Grid container spacing={2}>
                <Grid item xs={1}>
                { id !== undefined  && id.length > 0 &&
                <TextField
                    disabled={true}
                    inputRef={register}
                    label="Type"
                    name="type"
                    size="small"
                    variant="outlined"
                    defaultValue={type}
                />
                }
                { id === undefined &&
                <FormControl error={Boolean(errors.type)}  >
                    <Controller
                        as={
                            <Select name="type" size="small" >
                                <MenuItem value={"A"}>A</MenuItem>
                                <MenuItem value={"AAAA"}>AAAA</MenuItem>
                                <MenuItem value={"CNAME"}>CNAME</MenuItem>
                                <MenuItem value={"TXT"}>TXT</MenuItem>
                                <MenuItem value={"MX"}>MX</MenuItem>
                                <MenuItem value={"NS"}>NS</MenuItem>
                            </Select>
                        }
                        defaultValue={defaultRecordType}
                        name="type"
                        native={false}
                        control={control}
                    />
                    <FormHelperText>
                    {errors.type && errors.type.message}
                    </FormHelperText>
                    </FormControl>
                }
                </Grid>
                <Grid item  xs={4}>
                <TextField
                    inputRef={register}
                    InputProps={{
                        endAdornment: <InputAdornment position="end">.{zoneName}</InputAdornment>
                    }}
                    label="Name"
                    name="name"
                    size="small"
                    variant="outlined"
                    defaultValue={defaultName}
                />
                </Grid>

                <RecordFormContentInput {...{ register, unregister, setValue, watch, name: "content", defaultValue: content, type:selectedType, zoneName }} />

                <Grid item xs={2}>
            <TextField
                fullWidth
                required={true}
                InputLabelProps={{ shrink: true }}
                type="number"
                inputRef={register}
                label="TTL"
                name="ttl"
                size="small"
                variant="outlined"
                defaultValue={ttl}
            />
                </Grid>

            </Grid>

            <Button color="secondary"  type="submit" variant="contained">
                Save
            </Button>
        </form>
    );
};

RecordForm.propTypes = {
    id: PropTypes.string,
    fqdn: PropTypes.string,
    type: PropTypes.string,
    ttl: PropTypes.number,
    content: PropTypes.array,
    zoneName: PropTypes.string.isRequired,
    onCreated: PropTypes.func,
}
export default RecordForm;