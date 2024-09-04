import { Log } from './log/index';

export const APP = 'GoCloudK8sUserGroupFront';
export const APP_TITLE = 'Goéland';
export const VERSION = '0.3.1';
export const BUILD_DATE = '2024-09-04';
// eslint-disable-next-line no-undef
export const DEV = process.env.NODE_ENV === 'development';
export const HOME = DEV ? 'http://localhost:5173/' : '/';
// eslint-disable-next-line no-restricted-globals
const url = new URL(location.toString());
export const BACKEND_URL = DEV ? 'http://localhost:8888' : url.origin;
export const apiRestrictedUrl = `api/v1`;
export const getLog = (ModuleName, verbosityDev, verbosityProd) => (
  (DEV) ? new Log(ModuleName, verbosityDev) : new Log(ModuleName, verbosityProd)
);
