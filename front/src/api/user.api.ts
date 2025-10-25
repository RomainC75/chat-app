import type { TLoginResponse, TLoginsUser, TSignupUser, TVerifyResponse } from "../types/user.type";
import globalApi from "./global.api";

const basePath = "/user"

export const fetchSignupUser = async (data: TSignupUser): Promise<string> => {
    const response = await globalApi.post<string>(`${basePath}/signup`, data);
    return response.data;
}


export const fetchLoginUser = async (data: TLoginsUser): Promise<TLoginResponse> => {
    const response = await globalApi.post<TLoginResponse>(`${basePath}/login`, data);
    return response.data;
}


export const fetchVerify = async (): Promise<TVerifyResponse> => {
    const response = await globalApi.get<TVerifyResponse>(`${basePath}/verify`);
    return response.data;
}