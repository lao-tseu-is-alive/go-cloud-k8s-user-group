import sha256 from 'crypto-js/sha256';
import axios from 'axios';
import { getLog, APP, BACKEND_URL, apiRestrictedUrl } from '../config';

const log = getLog('Login', 4, 1);

export const getPasswordHash = (password) => sha256(password);

export const parseJwt = (token) => {
  const base64Url = token.split('.')[1];
  const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
  const jsonPayload = decodeURIComponent(atob(base64).split('').map((c) => `%${(`00${c.charCodeAt(0).toString(16)}`).slice(-2)}`).join(''));

  return JSON.parse(jsonPayload);
};

export const getToken = async (baseServerUrl, username, passwordHash) => {
  const data = {
    username,
    password_hash: `${passwordHash}`,
  };
  log.t('# IN getToken() data:', data);
  let response = null;
  try {
    response = await axios.post(`${baseServerUrl}/login`, data); // .then((response) => {
    log.l('getToken() axios.post Success ! response :', response.data);
    const jwtValues = parseJwt(response.data.token);
    log.l('getToken() token values : ', jwtValues);
    const dExpires = new Date(0);
    dExpires.setUTCSeconds(jwtValues.exp);
    log.l(`getToken() JWT token expiration : ${dExpires}`);
    if (response.status === 200) {
      if (typeof Storage !== 'undefined') {
        // Code for localStorage/sessionStorage.
        sessionStorage.setItem(`${APP}_goapi_jwt_session_token`, response.data.token);
        sessionStorage.setItem(`${APP}_goapi_idgouser`, jwtValues.User.user_id);
        sessionStorage.setItem(`${APP}_goapi_user_external_id`, jwtValues.User.external_id);
        sessionStorage.setItem(`${APP}_goapi_name`, jwtValues.User.name);
        sessionStorage.setItem(`${APP}_goapi_username`, jwtValues.User.login);
        sessionStorage.setItem(`${APP}_goapi_email`, jwtValues.User.email);
        sessionStorage.setItem(`${APP}_goapi_isadmin`, jwtValues.User.is_admin);
        sessionStorage.setItem(`${APP}_goapi_date_expiration`, jwtValues.exp);
      }
      return response.data;
    }
    log.w('axios get a bad status ! response was:', response);
    return response;
  } catch (e) {
    log.e('getToken() ## Try Catch ERROR ## error :', e);
    log.e('axios response was:', e.response);
    return e;
  }
};

export const getTokenStatus = async (baseServerUrl = BACKEND_URL) => {
  log.t('# IN getTokenStatus() ');
  axios.defaults.headers.common.Authorization = `Bearer ${sessionStorage.getItem(`${APP}_goapi_jwt_session_token`)}`;
  try {
    const res = await axios.get(`${baseServerUrl}/${apiRestrictedUrl}/status`);
    log.l('getTokenStatus() axios.get Success ! response :', res);
    const dExpires = new Date(0);
    dExpires.setUTCSeconds(res.data.exp);
    const msg = `getTokenStatus() JWT token expiration : ${dExpires}`;
    log.w(msg);
    const { data } = res;
    return {
      msg, err: null, status: res.status, data,
    };
  } catch (error) {
    const msg = `Error: in getTokenStatus() ## axios.get(${baseServerUrl}/${apiRestrictedUrl}/status) ERROR ## error :${error}`;
    log.w(msg);
    if (error.response) {
      const errResponse = error.response;
      return {
        msg, err: error, status: errResponse.status, data: null,
      };
    }
    return {
      msg, err: error, status: null, data: null,
    };
  }
};

export const clearSessionStorage = () => {
  // Code for localStorage/sessionStorage.
  sessionStorage.removeItem(`${APP}_goapi_jwt_session_token`);
  sessionStorage.removeItem(`${APP}_goapi_idgouser`);
  sessionStorage.removeItem(`${APP}_goapi_user_external_id`);
  sessionStorage.removeItem(`${APP}_goapi_name`);
  sessionStorage.removeItem(`${APP}_goapi_username`);
  sessionStorage.removeItem(`${APP}_goapi_email`);
  sessionStorage.removeItem(`${APP}_goapi_isadmin`);
  sessionStorage.removeItem(`${APP}_goapi_groups`);
  sessionStorage.removeItem(`${APP}_goapi_date_expiration`);
};

export const logoutAndResetToken = (baseServerUrl) => {
  log.t('# IN logoutAndResetToken()');
  axios.defaults.headers.common.Authorization = `Bearer ${sessionStorage.getItem(`${APP}_goapi_jwt_session_token`)}`;
  axios.get(`${baseServerUrl}/${apiRestrictedUrl}/logout`)
    .then((response) => {
      log.l('logoutAndResetToken() axios.get Success ! response :', response);
      clearSessionStorage();
    })
    .catch((error) => {
      log.e('logoutAndResetToken() ## axios.get ERROR ## error :', error);
    });
};

export const doesCurrentSessionExist = () => {
  log.t('# IN doesCurrentSessionExist() ');
  if (sessionStorage.getItem(`${APP}_goapi_jwt_session_token`) == null) return false;
  if (sessionStorage.getItem(`${APP}_goapi_idgouser`) == null) return false;
  if (sessionStorage.getItem(`${APP}_goapi_isadmin`) == null) return false;
  if (sessionStorage.getItem(`${APP}_goapi_email`) == null) return false;
  if (sessionStorage.getItem(`${APP}_goapi_date_expiration`) !== null) {
    const dateExpire = new Date(sessionStorage.getItem(`${APP}_goapi_date_expiration`));
    const now = new Date();
    if (now > dateExpire) {
      clearSessionStorage();
      log.w('# IN doesCurrentSessionExist() SESSION EXPIRED');
      return false;
    }
    // attention qu'une session existe en local veut pas dire que le jwt token est encore valide !
    return true;
  }
  log.w('# IN doesCurrentSessionExist() goapi_date_expiration was null ');
  return false;
};

export const getLocalJwtTokenAuth = () => {
  if (doesCurrentSessionExist()) {
    return `Bearer ${sessionStorage.getItem(`${APP}_goapi_jwt_session_token`)}`;
  }
  return '';
};

export const getUserEmail = () => {
  if (doesCurrentSessionExist()) {
    return `${sessionStorage.getItem(`${APP}_goapi_email`)}`;
  }
  return '';
};

export const getUserId = () => {
  if (doesCurrentSessionExist()) {
    return parseInt(`${sessionStorage.getItem(`${APP}_goapi_idgouser`)}`, 10);
  }
  return 0;
};

export const getUserName = () => {
  if (doesCurrentSessionExist()) {
    return `${sessionStorage.getItem(`${APP}_goapi_name`)}`;
  }
  return '';
};

export const getUserLogin = () => {
  if (doesCurrentSessionExist()) {
    return `${sessionStorage.getItem(`${APP}_goapi_username`)}`;
  }
  return '';
};

export const getUserIsAdmin = () => {
  if (doesCurrentSessionExist()) {
    return (sessionStorage.getItem(`${APP}_goapi_isadmin`) === 'true');
  }
  return false;
};

export const getUserFirstGroups = () => {
  if (doesCurrentSessionExist()) {
    if (sessionStorage.getItem(`${APP}_goapi_groups`) == null) return null;
    if (sessionStorage.getItem(`${APP}_goapi_groups`) === 'null') return null;
    // let's clone it and converting to an array of integers
    const tmpArr = sessionStorage.getItem(`${APP}_goapi_groups`);
    if (tmpArr.indexOf(',') > 0) {
      const firstFiltered = tmpArr.split(',').map((e) => +e);
      return firstFiltered[0];
    }
    return parseInt(tmpArr, 10);
  }
  return null;
};

export const getUserGroupsArray = () => {
  if (doesCurrentSessionExist()) {
    if (sessionStorage.getItem(`${APP}_goapi_groups`) == null) return null;
    if (sessionStorage.getItem(`${APP}_goapi_groups`) === 'null') return null;
    // let's clone it and converting to an array of integers
    const tmpArr = sessionStorage.getItem(`${APP}_goapi_groups`);
    if (tmpArr.indexOf(',') > 0) {
      return tmpArr.split(',').map((i) => parseInt(i, 10));
    }
    return [parseInt(tmpArr, 10)];
  }
  return null;
};

export const isUserHavingGroups = () => {
  if (doesCurrentSessionExist()) {
    if (sessionStorage.getItem(`${APP}_goapi_groups`) == null) return false;
    if (sessionStorage.getItem(`${APP}_goapi_groups`) === 'null') return false;
    // let's clone it and converting to an array of integers
    const tmp = sessionStorage.getItem(`${APP}_goapi_groups`);
    if (tmp.indexOf(',') > 0) {
      return true;
    }
    if (parseInt(tmp, 10) > 0) {
      return true;
    }
    return false;
  }
  return false;
};
