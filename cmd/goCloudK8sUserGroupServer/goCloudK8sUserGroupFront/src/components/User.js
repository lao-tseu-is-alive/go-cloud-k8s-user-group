import axios from 'axios';
import { functionExist, getErrorMessage } from '../tools/utils';
import {apiRestrictedUrl, BACKEND_URL, getLog} from '../config';
import { getLocalJwtTokenAuth } from './Login';

const log = getLog('User', 4, 1);

// User Singleton and stateless class to get and persist data to backend
const user = {
  getList: (callbackLoaded) => {
    const method = 'getList';
    log.t(`## IN ${method}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.get(`${BACKEND_URL}/${apiRestrictedUrl}/users`)
      .then((resp) => {
        log.t(`## IN ${method} axios get success resp.data :`, resp.data);
        if (functionExist(callbackLoaded)) {
          callbackLoaded(resp.data, 'SUCCESS');
        }
      })
      .catch((err) => {
        const errMessage = getErrorMessage(method, `## ERREUR RESEAU DANS ${method} PENDANT UN APPEL DISTANT axios.get`, err, log);
        if (functionExist(callbackLoaded)) callbackLoaded(err, errMessage);
      });
  },

  getUser: (idUser, callbackLoaded) => {
    const method = 'getUser';
    log.t(`## IN ${method}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.get(`${BACKEND_URL}/${apiRestrictedUrl}/users/${idUser}`)
      .then((resp) => {
        log.t(`## IN ${method} axios get success resp.data :`, resp.data);
        if (functionExist(callbackLoaded)) {
          callbackLoaded(resp.data, 'SUCCESS');
        }
      })
      .catch((err) => {
        const errMessage = getErrorMessage(method, `## ERREUR RESEAU DANS ${method} PENDANT UN APPEL DISTANT axios.get`, err, log);
        if (functionExist(callbackLoaded)) callbackLoaded(err, errMessage);
      });
  },

  newUser: (data, callbackLoaded) => {
    const method = 'newUser';
    log.t(`## IN ${method}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.post(`${BACKEND_URL}/${apiRestrictedUrl}/users`, data)
      .then((resp) => {
        log.t(`## IN ${method} axios get success resp.data :`, resp.data);
        if (functionExist(callbackLoaded)) {
          callbackLoaded(resp.data, 'SUCCESS');
        }
      })
      .catch((err) => {
        const errMessage = getErrorMessage(method, `## ERREUR RESEAU DANS ${method} PENDANT UN APPEL DISTANT axios.post`, err, log);
        if (functionExist(callbackLoaded)) callbackLoaded(err, errMessage);
      });
  },

  modifyUser: (data, callbackLoaded) => {
    const method = 'modifyUser';
    log.t(`## IN ${method}`);
    const { id } = data;
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.put(`${BACKEND_URL}/${apiRestrictedUrl}/users/${id}`, data)
      .then((resp) => {
        log.t(`## IN ${method} axios put success resp.data :`, resp.data);
        if (functionExist(callbackLoaded)) {
          callbackLoaded(resp.data, 'SUCCESS');
        }
      })
      .catch((err) => {
        const errMessage = getErrorMessage(method, `## ERREUR RESEAU DANS ${method} PENDANT UN APPEL DISTANT axios.put`, err, log);
        if (functionExist(callbackLoaded)) callbackLoaded(err, errMessage);
      });
  },

  deleteUser: (idUserToDelete, callbackLoaded) => {
    const method = 'deleteUser';
    log.t(`## IN ${method} id:${idUserToDelete}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.delete(`${BACKEND_URL}/${apiRestrictedUrl}/users/${idUserToDelete}`)
      .then((resp) => {
        log.t(`## IN ${method} axios delete success resp.data :`, resp.data);
        if (functionExist(callbackLoaded)) {
          callbackLoaded(resp.data, 'SUCCESS');
        }
      })
      .catch((err) => {
        const errMessage = getErrorMessage(method, `## ERREUR RESEAU DANS ${method} PENDANT UN APPEL DISTANT axios.delete`, err, log);
        if (functionExist(callbackLoaded)) callbackLoaded(err, errMessage);
      });
  },
};
// prevents modification to properties and values of the user singleton
Object.freeze(user);
export default user;
