import {createAsyncThunk, createSlice} from '@reduxjs/toolkit';
import {createRecord, deleteRecord, getZone, getZones, setApiAuthToken, updateRecord} from "../../api";

export const fetchZones = createAsyncThunk(
    'zones/fetchZones',
    async () => {
        return await getZones()
    }
);

export const fetchZoneByName = createAsyncThunk(
    'zones/fetchZoneByName',
    async (name, thunkAPI) => {
        console.log(name)
        return await getZone(name)
    }
);

export const deleteRecordByZoneAndId = createAsyncThunk(
    'zones/deleteRecordByZoneAndId',
    async (data, thunkAPI) => {
        return await deleteRecord(data.zoneName, data.recordId)
    }
);

export const createRecordWithZoneNameAndRecord = createAsyncThunk(
    'zones/createRecord',
    async (data, thunkAPI) => {
        return await createRecord(data.zoneName, data.record)
    }
);

export const updateRecordWithZoneNameAndRecord = createAsyncThunk(
    'zones/updateRecord',
    async (data, thunkAPI) => {
        return await updateRecord(data.zoneName, data.recordId, data.record)
    }
);

export const updateAuthCredentials = createAsyncThunk(
    'api/updateAuthCredentials',
    async (payload, thunkAPI) => {
        const { user, password } = payload
        const authBasic =  btoa(user+":"+password)
        localStorage.setItem('authBasic',authBasic);
        setApiAuthToken(authBasic)
        return await getZones().then(function () {
          return authBasic;
        })
    }
);


const initialAuthToken = localStorage.getItem("authBasic")||"";
setApiAuthToken(initialAuthToken);

export const zonesSlices = createSlice({
    name: 'zones',
    initialState: {
        entities: [], loading: 'idle', selectedZone: {}, authToken : initialAuthToken, authIsValid:false //TODO voir pour recuper bone valeur
    },
    reducers: {
        selectZone(state, action) {
            const { zoneName } = action.payload
            state.selectedZone = {
                name: zoneName,
                records: [],
            }
        }
    },
    extraReducers: {
        [fetchZoneByName.fulfilled]: (state, action) => {
            state.selectedZone.records = action.payload
            state.loading = 'idle'
        },
        [fetchZoneByName.pending]: (state, action) => {
            if (state.loading === 'idle') {
                state.loading = 'pending'
            }
        },
        [fetchZones.pending]: (state, action) => {
            if (state.loading === 'idle') {
                state.loading = 'pending'
            }
        },
        // Add reducers for additional action types here, and handle loading state as needed
        [fetchZones.fulfilled]: (state, action) => {
            // Add user to the state array
            state.entities = action.payload
            state.loading = 'idle'
        },
        [updateAuthCredentials.fulfilled]: (state, action) => {
            state.authToken = action.payload
            state.authIsValid = true
        },
        [updateAuthCredentials.rejected]: (state, action) => {
            state.authToken = ""
            state.authIsValid = false
        },
        [deleteRecordByZoneAndId.fulfilled]: (state, action) => {
            state.selectedZone.records = state.selectedZone.records.filter((record, i) => record.GeneratedId !== action.payload)
        },
        [createRecordWithZoneNameAndRecord.fulfilled]: (state, action) => {
            state.selectedZone.records.push(action.payload)
        },
        [updateRecordWithZoneNameAndRecord.fulfilled]: (state, action) => {
            state.selectedZone.records = state.selectedZone.records.filter((record, i) => record.GeneratedId !== action.payload.oldRecordId)
            state.selectedZone.records.push(action.payload.newRecord)
        }

    }
})


export const { selectZone } = zonesSlices.actions
export default zonesSlices.reducer;
