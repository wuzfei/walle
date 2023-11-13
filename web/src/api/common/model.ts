export interface Version {
  version: string;
  hash: string;
  timestamp: string;
  release: boolean;
}

export interface ServerInfoOs {
  goos: string;
  goarch: string;
  cpu_num: number;
  compiler: string;
  go_version: string;
  goroutine_num: number;
}

export interface ServerInfoCpu {
  used_percent: number;
  cores: number;
}
export interface ServerInfoRam {
  used: string;
  total: string;
  used_percent: number;
}

export interface ServerInfoDisk {
  used: string;
  total: string;
  used_percent: number;
}

export interface ServerInfo {
  user: string;
  hostname: string;
  os: ServerInfoOs;
  cpu: ServerInfoCpu;
  ram: ServerInfoRam;
  disk: ServerInfoDisk;
}
