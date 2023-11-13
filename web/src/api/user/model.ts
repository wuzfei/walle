import { BasicPageParams, BasicFetchResult } from '../baseModel';

export type ListReq = BasicPageParams & {
  username?: string;
};

export interface CreateReq {
  username: string;
  email: string;
  password: string;
  status: number;
}

export interface UpdateReq extends CreateReq {
  id: number;
}

export interface ListItem {
  id: number;
  username: string;
  email: string;
  password: string;
  status: number;
  created_at: string;
}

export type ListItemRes = BasicFetchResult<ListItem>;
