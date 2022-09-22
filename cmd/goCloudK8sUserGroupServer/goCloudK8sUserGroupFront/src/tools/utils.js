export const isNullOrUndefined = (variable) => typeof variable === 'undefined' || variable === null;
export const isEmpty = (variable) => typeof variable === 'undefined' || variable === null || variable === '';
export const isFunction = (f) => typeof f === 'function';
export const functionExist = (functionName) => typeof functionName !== 'undefined' && functionName !== null && isFunction(functionName);
