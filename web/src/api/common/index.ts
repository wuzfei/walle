import { defHttp } from '/@/utils/http/axios';
import { Version, ServerInfo } from './model';

enum Api {
  Version = '/version',
  ServerInfo = '/server_info',
}

export function getVersion() {
  return defHttp.get<Version>({ url: Api.Version });
}

export function getServerInfo() {
  return defHttp.get<ServerInfo>({ url: Api.ServerInfo });
}
