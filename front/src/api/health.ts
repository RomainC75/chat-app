import globalApi from "./global.api";

const basePath = '/';

export const isHealthy = async (): Promise<boolean> => {
    const response = await globalApi.get<any>(
      `${basePath}`
    );
    return response.data.isTaken;
  };