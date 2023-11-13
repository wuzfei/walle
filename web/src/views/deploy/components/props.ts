export type statusType = 'wait' | 'error' | 'finish' | 'process';

export interface ServerSteps {
  server_id: number;
  name: string;
  current: number;
  status: statusType;
  //status: PropType<'wait' | 'error' | 'finish' | 'process'>;
}
