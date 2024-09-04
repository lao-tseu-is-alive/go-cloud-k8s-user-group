export const isFunction = (f) => typeof f === "function";
export const functionExist = (f) => typeof f !== "undefined" && f !== null &&  isFunction(f);
export const isNullOrUndefined = (v) =>  typeof v === "undefined" || v === null;
export const isEmpty = (v) => isNullOrUndefined(v) || v === "";
export const getErrorMessage = (method, msg, err, log = null) => {
  let errMessage = msg;
  let lastErrorMsg = '';
  if (!isNullOrUndefined(log)) {
    log.e(errMessage);
  }
  if (err.response) {
    const errResponse = err.response;
    if (!isNullOrUndefined(log)) {
      log.e(' -- The request was made, but the server responded with a status code > 2xx', errResponse, errResponse.data);
    }
    lastErrorMsg = ` ${method} : La requête http a reçu en retour un code status=${errResponse.status} >200 ! `;
    if (typeof errResponse.data === 'object') {
      errMessage += `${lastErrorMsg} <br> Message serveur : ${errResponse.data.message}`;
    } else {
      errMessage += `${lastErrorMsg} <br> Message serveur : ${errResponse.data}`;
    }
  } else if (err.request) {
    if (!isNullOrUndefined(log)) {
      log.e(' -- The request was made, but no response was received from the server', err.request);
    }
    lastErrorMsg = ` ${method} : La requête http n'a pas reçu de réponse du serveur ! `;
    errMessage += `${lastErrorMsg} <br> Message serveur : ${err.request}`;
  } else {
    if (!isNullOrUndefined(log)) {
      log.e(' -- Something happened in setting up the request that triggered an Error', err.message);
    }
    lastErrorMsg = ` ${method} : La requête http est erronée et n'a pas pu se faire! `;
    errMessage += `${lastErrorMsg} <br> Message serveur : ${err.message}`;
  }
  return errMessage;
};
