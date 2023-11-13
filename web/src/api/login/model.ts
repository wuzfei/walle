export interface LoginParams {
  email: string;
  password: string;
  remember?: boolean;
}

export interface RoleInfo {
  roleName: string;
  value: string;
}

export interface SpaceInfo {
  space_name: string;
  space_id: number;
  status: number;
  role: string;
}

/**
 * @description: Login interface return value
 */
export interface LoginResultModel {
  user_id: string | number;
  token: string;
  token_expire: number;
  refresh_token: string;
  refresh_token_expire: number;
  role?: RoleInfo;
}

/**
 * @description: Get user information return value
 */
export interface GetUserInfoModel {
  // 用户id
  user_id: string | number;
  // 用户名
  email: string;

  // 真实名字
  username: string;

  role: string;
  current_space_id: number;

  roles?: RoleInfo[];

  spaces: SpaceInfo[];
}
