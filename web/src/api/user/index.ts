import { defHttp } from '/@/utils/http/axios';
import { CreateReq, ListItemRes, ListReq, UpdateReq } from './model';
import { GetOptionItemsModel } from '../baseModel';

enum Api {
  User = '/user',
  UserId = '/user/{id}',
  UserOptions = '/user/options',
}

export const getUserListByPage = (params?: ListReq) =>
  defHttp.get<ListItemRes>({ url: Api.User, params });

export const createUser = (params: CreateReq) => defHttp.post({ url: Api.User, params: params });

export const updateUser = (params: UpdateReq) => defHttp.put({ url: Api.User, params: params });

export const deleteUser = (id: number) =>
  defHttp.delete({ url: Api.UserId.replace('{id}', id.toString()) });

export const getUserOptions = (params?: ListReq) =>
  defHttp.get<GetOptionItemsModel>({ url: Api.UserOptions, params });
