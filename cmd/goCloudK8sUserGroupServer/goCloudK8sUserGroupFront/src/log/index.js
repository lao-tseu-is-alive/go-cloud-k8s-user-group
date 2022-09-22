/* eslint-disable no-underscore-dangle */
/* eslint-disable  no-console */
/**
 * Log library for javascript
 */
export const logLevel = {
  info: 4,
  trace: 3,
  warn: 2,
  err: 1,
  none: 0,
};
const logType = {
  linfo: 'background: #b3e5fc; color: #000',
  ltrace: 'background: #c8e6c9; color: #37474f',
  lwarn: 'background: #ff9800; color: #020202',
  lerr: 'background: #FF2325; color: #BFFF1A',
};

// eslint-disable-next-line func-names
const _log = function (moduleName, msg, logtype, ...args) {
  const prefix = `${moduleName} ➜`;
  switch (logtype) {
    case logtype.lerr:

      console.error(`%c ${prefix} ${msg}`, logtype);
      console.trace();
      break;
    case logtype.lwarn:
      console.warn(`%c ${prefix} ${msg}`, logtype);
      break;
    default:
      console.log(`%c ${prefix} ${msg}`, logtype);
      break;
  }
  if (args.length > 0) {
    // eslint-disable-next-line no-plusplus
    for (let i = 0; i < args.length; i++) {
      console.log(args[i]);
    }
  }
};
export class Log {
  constructor(moduleName = '', level = 4) {
    this._moduleName = moduleName;
    this._logLevel = level;
  }

  l(msg, ...args) {
    if (this._logLevel >= logLevel.info) {
      _log(this._moduleName, msg, logType.linfo, ...args);
    }
  }

  t(msg, ...args) {
    if (this._logLevel >= logLevel.trace) {
      _log(this._moduleName, msg, logType.ltrace, ...args);
    }
  }

  w(msg, ...args) {
    if (this._logLevel >= logLevel.warn) {
      _log(this._moduleName, msg, logType.lwarn, ...args);
    }
  }

  e(msg, ...args) {
    if (this._logLevel >= logLevel.err) {
      _log(this._moduleName, msg, logType.lerr, ...args);
    }
  }
}
