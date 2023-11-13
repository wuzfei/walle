import { ListReq, CreateReq, ListItemRes, ListItem } from './model';
import { defHttp } from '/@/utils/http/axios';
import { getWebsocketApiUrl } from '/@/utils/common';

enum Api {
  Deploy = '/deploy',
  DeployId = '/deploy/{id}',
  DeployRelease = '/deploy/{id}/release',
  DeployStopRelease = '/deploy/{id}/stop_release',
  DeployAudit = '/deploy/{id}/audit',
  DeployConsoleWs = '/deploy/{id}/console',
}

export const getDeployListByPage = (params?: ListReq) =>
  defHttp.get<ListItemRes>({ url: Api.Deploy, params });

export const createDeploy = (params: CreateReq) =>
  defHttp.post({ url: Api.Deploy, params: params });

export const deleteDeploy = (id: number) =>
  defHttp.delete({ url: Api.DeployId.replace('{id}', id.toString()) });

export const detailDeploy = (id: number, notAlertErrMsg: boolean | undefined) =>
  defHttp.get<ListItem>(
    { url: Api.DeployId.replace('{id}', id.toString()) },
    notAlertErrMsg ? { errorMessageMode: 'none' } : {},
  );

export const startDeploy = (id: number, notAlertErrMsg: boolean | undefined) =>
  defHttp.get<ListItem>(
    { url: Api.DeployRelease.replace('{id}', id.toString()) },
    notAlertErrMsg ? { errorMessageMode: 'none' } : {},
  );

export const stopDeploy = (id: number, notAlertErrMsg: boolean | undefined) =>
  defHttp.get<ListItem>(
    { url: Api.DeployStopRelease.replace('{id}', id.toString()) },
    notAlertErrMsg ? { errorMessageMode: 'none' } : {},
  );

export const auditDeploy = (id: number, audit: boolean) =>
  defHttp.post<ListItem>({
    url: Api.DeployAudit.replace('{id}', id.toString()),
    params: { audit: audit },
  });

export const getDeployConsoleWs = (id: number) =>
  getWebsocketApiUrl(Api.DeployConsoleWs.replace('{id}', id.toString()));
