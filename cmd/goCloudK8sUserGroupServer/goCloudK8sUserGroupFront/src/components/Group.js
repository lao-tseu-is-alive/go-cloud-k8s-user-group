import axios from 'axios';
import { functionExist, getErrorMessage } from '../tools/utils';
import { apiRestrictedUrl, BACKEND_URL, getLog } from '../config';
import { getLocalJwtTokenAuth } from './Login';

const log = getLog('Group', 4, 1);

// Group Singleton and stateless class to get and persist data to backend
const group = {
  getList: (callbackLoaded) => {
    const method = 'getList';
    log.t(`## IN ${method}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.get(`${BACKEND_URL}/${apiRestrictedUrl}/groups`)
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

  getGroup: (idGroup, callbackLoaded) => {
    const method = 'getGroup';
    log.t(`## IN ${method}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.get(`${BACKEND_URL}/${apiRestrictedUrl}/groups/${idGroup}`)
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

  newGroup: (data, callbackLoaded) => {
    const method = 'newGroup';
    log.t(`## IN ${method}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.post(`${BACKEND_URL}/${apiRestrictedUrl}/groups`, data)
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

  modifyGroup: (data, callbackLoaded) => {
    const method = 'modifyGroup';
    log.t(`## IN ${method}`);
    const { id } = data;
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.put(`${BACKEND_URL}/${apiRestrictedUrl}/groups/${id}`, data)
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

  deleteGroup: (idGroupToDelete, callbackLoaded) => {
    const method = 'deleteGroup';
    log.t(`## IN ${method} id:${idGroupToDelete}`);
    axios.defaults.headers.common.Authorization = getLocalJwtTokenAuth();
    axios.delete(`${BACKEND_URL}/${apiRestrictedUrl}/groups/${idGroupToDelete}`)
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
// prevents modification to properties and values of the group singleton
Object.freeze(group);
export default group;
