import axios from 'axios';
import { functionExist } from '../tools/utils';
import { BACKEND_URL, getLog } from '../config';
import { getLocalJwtTokenAuth } from './Login';

const log = getLog('User', 4, 1);
const getErrorMessage = (method, msg, err) => {
  let errMessage = msg;
  let lastErrorMsg = '';
  log.e(errMessage);
  if (err.response) {
    const errResponse = err.response;
    log.e(' -- The request was made, but the server responded with a status code > 2xx', errResponse, errResponse.data);
    lastErrorMsg = ` ${method} : La requete http a recu en retour un code status=${errResponse.status} >200 ! `;
    if (typeof errResponse.data === 'object') {
      errMessage += `${lastErrorMsg} <br> Message serveur : ${errResponse.data.message}`;
    } else {
      errMessage += `${lastErrorMsg} <br> Message serveur : ${errResponse.data}`;
    }
  } else if (err.request) {
    log.e(' -- The request was made, but no response was received from the server', err.request);
    lastErrorMsg = ` ${method} : La requete http n'a pas reçu de reponse du serveur ! `;
    errMessage += `${lastErrorMsg} <br> Message serveur : ${err.request}`;
  } else {
    log.e(' -- Something happened in setting up the request that triggered an Error', err.message);
    lastErrorMsg = ` ${method} : La requete http est erronée et n'a pas pu se faire! `;
    errMessage += `${lastErrorMsg} <br> Message serveur : ${err.message}`;
  }
  return errMessage;
};
// User Singleton and stateless class to get and persist data to backend
const user = {
  getList: (callbackLoaded) => {
    const method = 'getList';
    log.t(`## IN ${method}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.get(`${BACKEND_URL}/api/users`)
      .then((resp) => {
        log.t(`## IN ${method} axios get success resp.data :`, resp.data);
        if (functionExist(callbackLoaded)) {
          callbackLoaded(resp.data, 'SUCCESS');
        }
      })
      .catch((err) => {
        const errMessage = getErrorMessage(method, `## ERREUR RESEAU DANS ${method} PENDANT UN APPEL DISTANT axios.get`, err);
        if (functionExist(callbackLoaded)) callbackLoaded(null, errMessage);
      });
  },

  getUser: (idUser, callbackLoaded) => {
    const method = 'getUser';
    log.t(`## IN ${method}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.get(`${BACKEND_URL}/api/users/${idUser}`)
      .then((resp) => {
        log.t(`## IN ${method} axios get success resp.data :`, resp.data);
        if (functionExist(callbackLoaded)) {
          callbackLoaded(resp.data, 'SUCCESS');
        }
      })
      .catch((err) => {
        const errMessage = getErrorMessage(method, `## ERREUR RESEAU DANS ${method} PENDANT UN APPEL DISTANT axios.get`, err);
        if (functionExist(callbackLoaded)) callbackLoaded(null, errMessage);
      });
  },

  newUser: (data, callbackLoaded) => {
    const method = 'newUser';
    log.t(`## IN ${method}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.post(`${BACKEND_URL}/api/users`, data)
      .then((resp) => {
        log.t(`## IN ${method} axios get success resp.data :`, resp.data);
        if (functionExist(callbackLoaded)) {
          callbackLoaded(resp.data, 'SUCCESS');
        }
      })
      .catch((err) => {
        const errMessage = getErrorMessage(method, `## ERREUR RESEAU DANS ${method} PENDANT UN APPEL DISTANT axios.post`, err);
        if (functionExist(callbackLoaded)) callbackLoaded(null, errMessage);
      });
  },

  modifyUser: (data, callbackLoaded) => {
    const method = 'modifyUser';
    log.t(`## IN ${method}`);
    const { id } = data;
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.put(`${BACKEND_URL}/api/users/${id}`, data)
      .then((resp) => {
        log.t(`## IN ${method} axios put success resp.data :`, resp.data);
        if (functionExist(callbackLoaded)) {
          callbackLoaded(resp.data, 'SUCCESS');
        }
      })
      .catch((err) => {
        const errMessage = getErrorMessage(method, `## ERREUR RESEAU DANS ${method} PENDANT UN APPEL DISTANT axios.put`, err);
        if (functionExist(callbackLoaded)) callbackLoaded(null, errMessage);
      });
  },

  deleteUser: (idUserToDelete, callbackLoaded) => {
    const method = 'deleteUser';
    log.t(`## IN ${method} id:${idUserToDelete}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.delete(`${BACKEND_URL}/api/users/${idUserToDelete}`)
      .then((resp) => {
        log.t(`## IN ${method} axios delete success resp.data :`, resp.data);
        if (functionExist(callbackLoaded)) {
          callbackLoaded(resp.data, 'SUCCESS');
        }
      })
      .catch((err) => {
        const errMessage = getErrorMessage(method, `## ERREUR RESEAU DANS ${method} PENDANT UN APPEL DISTANT axios.delete`, err);
        if (functionExist(callbackLoaded)) callbackLoaded(null, errMessage);
      });
  },
};
// prevents modification to properties and values of the user singleton
Object.freeze(user);
export default user;
