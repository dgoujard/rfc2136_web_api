import ky from 'ky';


let api = ky.extend({
    prefixUrl: "/",
});

export function setApiAuthToken(newToken) {
    api = api.extend({
        hooks: {
            beforeRequest: [
                request => {
                    request.headers.set('Authorization', 'Basic '+newToken);
                }
            ]
        }
    });
}

export async function getZones() {
    const response = await api.get('api/v1/zones');
    if (!response.ok) {
        console.log(response); //TODO return exception on error to handle it in the store
    }

    const parsed = await response.json();

    console.log(parsed);
    return parsed;
}

export async function getZone(name) {
    const response = await api.get('api/v1/zones/'+name);
    if (!response.ok) {
        console.log(response); //TODO return exception on error to handle it in the store
    }

    const parsed = await response.json();

    console.log(parsed);
    return parsed;
}

export async function deleteRecord(zoneName,recordId) {
    const response = await api.delete('api/v1/zones/'+zoneName+"/"+recordId);
    if (!response.ok) {
        console.log(response); //TODO return exception on error to handle it in the store
        return false;
    }
    return recordId
}

export async function createRecord(zoneName,record) {
    let payload = {
        "Fqdn": record.name+"."+zoneName+".",
        "TTL": parseInt(record.ttl),
        "Type": record.type,
        "Values": record.content
    }
    const response = await api.post('api/v1/zones/'+zoneName,{json:payload});
    if (!response.ok) {
        console.log(response); //TODO return exception on error to handle it in the store
        return false;
    }
    const parsed = await response.json();

    console.log(parsed)
    return parsed;
}


export async function updateRecord(zoneName,recordId,record) {
    await deleteRecord(zoneName, recordId)
    const newRecord = await createRecord(zoneName, record)
    return {oldRecordId:recordId,newRecord: newRecord}
}
