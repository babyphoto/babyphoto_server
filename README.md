# babyphoto_server

## React-Native File upload

```javascript
var formData = new FormData();
formData.append('file', {
    uri: data.file.uri,
    name: data.file.name,
    type: data.file.type,
});

fetch('http://112.169.11.118:38080/api/files/upload', {
    method: 'POST',
    body: formData,
})
.then(response => {
if (response.status === 200) {
    return response.json();
}
})
.then(responseJson => {
console.log(responseJson);
if (res) {
    res(responseJson);
}
})
.catch(error => {
console.error(error);
if (err) {
    err(error);
}
});
```