import { defHttp } from '/@/utils/http/axios';
import { LoginParams, LoginResultModel, GetUserInfoModel } from './model';
import { ErrorMessageMode } from '/#/axios';

enum Api {
  Login = '/login',
  Logout = '/logout',
  RefreshToken = '/refresh_token',
  GetUserInfo = '/user_info',
}

/**
 * @description: user login api
 */
export function loginApi(params: LoginParams, mode: ErrorMessageMode = 'modal') {
  return defHttp.post<LoginResultModel>(
    {
      url: Api.Login,
      params,
    },
    {
      errorMessageMode: mode,
    },
  );
}

/**
 * @description: getUserInfo
 */
export function getUserInfo() {
  return defHttp.get<GetUserInfoModel>({ url: Api.GetUserInfo }, { errorMessageMode: 'none' });
}

export function refreshToken(refresh_token: string) {
  return defHttp.post<LoginResultModel>({
    url: Api.RefreshToken,
    params: { refresh_token: refresh_token },
  });
}

export function doLogout() {
  return defHttp.post({ url: Api.Logout });
}
