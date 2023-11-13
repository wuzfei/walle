import { getAppEnvConfig, isDevMode } from '/@/utils/env';

export const getWebsocketApiUrl = (path: Nullable<string>): string => {
  let host = window.location.host;
  const env = getAppEnvConfig();
  let preApi = '';
  if (env.VITE_GLOB_API_URL_PREFIX) {
    preApi = env.VITE_GLOB_API_URL_PREFIX;
  }

  if (isDevMode()) {
    host = 'localhost:8989';
  }

  return `ws://${host}${preApi}${path}`;
};
