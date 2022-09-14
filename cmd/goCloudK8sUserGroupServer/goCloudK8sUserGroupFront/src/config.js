import { Log } from './log/index';

export const APP = 'GoCloudK8sUserGroupFront';
export const APP_TITLE = 'Goéland';
export const VERSION = '0.1.2';
export const BUILD_DATE = '2022-09-13';
// eslint-disable-next-line no-undef
export const DEV = process.env.NODE_ENV === 'development';
export const HOME = DEV ? 'http://localhost:5173/' : '/';
export const DEFAULT_BASE_SERVER_URL = DEV ? 'http://localhost:5173/' : 'https://goeland.io/';
export const getLog = (ModuleName, verbosityDev, verbosityProd) => (
  (DEV) ? new Log(ModuleName, verbosityDev) : new Log(ModuleName, verbosityProd)
);
