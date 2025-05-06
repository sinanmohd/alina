export type ServerConfig = {
  file_size_limit: number
  chunk_size: number
}

export const useServerConfig = () => {
  return useState<ServerConfig>('serverConfig');
}
