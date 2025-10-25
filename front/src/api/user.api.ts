import type { TSignupUser } from "../types/user.type";
import globalApi from "./global.api";

const basePath = "/user"

export const fetchSignupUser = async (data: TSignupUser): Promise<string> => {
    const response = await globalApi.post<string>(`${basePath}/signup`, data);
    return response.data;
}
