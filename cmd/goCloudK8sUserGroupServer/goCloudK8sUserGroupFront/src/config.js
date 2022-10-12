import { Log } from './log/index';

export const APP = 'GoCloudK8sUserGroupFront';
export const APP_TITLE = 'Goéland';
export const VERSION = '0.1.3';
export const BUILD_DATE = '2022-10-12';
// eslint-disable-next-line no-undef
export const DEV = process.env.NODE_ENV === 'development';
export const HOME = DEV ? 'http://localhost:5173/' : '/';
// eslint-disable-next-line no-restricted-globals
const url = new URL(location.toString());
export const BACKEND_URL = DEV ? 'http://localhost:8888' : url.origin;
export const getLog = (ModuleName, verbosityDev, verbosityProd) => (
  (DEV) ? new Log(ModuleName, verbosityDev) : new Log(ModuleName, verbosityProd)
);
