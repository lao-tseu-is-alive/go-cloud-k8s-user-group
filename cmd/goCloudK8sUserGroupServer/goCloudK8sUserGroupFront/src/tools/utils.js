export const isNullOrUndefined = (variable) => typeof variable === 'undefined' || variable === null;
export const isEmpty = (variable) => typeof variable === 'undefined' || variable === null || variable === '';
export const isFunction = (f) => typeof f === 'function';
export const functionExist = (functionName) => typeof functionName !== 'undefined' && functionName !== null && isFunction(functionName);
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
    lastErrorMsg = ` ${method} : La requete http a recu en retour un code status=${errResponse.status} >200 ! `;
    if (typeof errResponse.data === 'object') {
      errMessage += `${lastErrorMsg} <br> Message serveur : ${errResponse.data.message}`;
    } else {
      errMessage += `${lastErrorMsg} <br> Message serveur : ${errResponse.data}`;
    }
  } else if (err.request) {
    if (!isNullOrUndefined(log)) {
      log.e(' -- The request was made, but no response was received from the server', err.request);
    }
    lastErrorMsg = ` ${method} : La requete http n'a pas reçu de reponse du serveur ! `;
    errMessage += `${lastErrorMsg} <br> Message serveur : ${err.request}`;
  } else {
    if (!isNullOrUndefined(log)) {
      log.e(' -- Something happened in setting up the request that triggered an Error', err.message);
    }
    lastErrorMsg = ` ${method} : La requete http est erronée et n'a pas pu se faire! `;
    errMessage += `${lastErrorMsg} <br> Message serveur : ${err.message}`;
  }
  return errMessage;
};
