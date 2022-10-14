export function fetchData(url, refData, refLoaded, refError) {
  const method = `fetchLayersData(${url})`;
  console.log(`##--> ${method}`);
  refLoaded.value = false;
  return fetch(url, {
    method: 'get',
    headers: {
      'content-type': 'application/json',
    },
  })
    .then((response) => {
      if (!response.ok) {
        console.log(`#💥💥 ERROR in ${method}, status: ${response.status}`);
        const error = new Error(response.statusText);
        error.json = response.json();
        throw error;
      }
      console.log(`# ${method}, after fetch, status ${response.status}`);
      return response.json();
    })
    .then((json) => {
      console.log(`# ${method}, after json() ${typeof json}`);
      if (json instanceof Array) {
        refData.value = json;
      } else {
        refError.value = `# ${method}, json is not an instanceof Array`;
      }
    })
    .catch((err) => {
      console.log(`#💥💥 ERROR in ${method}, catch err:${err}`);
      refError.value = err;
      // In case a custom JSON error response was provided
      if (err.json) {
        return err.json.then((json) => {
          // set the JSON response message
          refError.value.message = json.message;
        });
      }
    })
    .then(() => {
      refLoaded.value = true;
    });
}
