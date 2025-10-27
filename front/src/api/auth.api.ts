import type { TLoginResponse, TLoginsUser, TSignupUser, TVerifyResponse } from "../types/user.type";
import openApi from "./open.api";
import globalApi from "./secure.api";

const basePath = "/auth"

export const fetchSignupUser = async (data: TSignupUser): Promise<string> => {
    const response = await openApi.post<string>(`${basePath}/signup`, data);
    return response.data;
}


export const fetchLoginUser = async (data: TLoginsUser): Promise<TLoginResponse> => {
    const response = await openApi.post<TLoginResponse>(`${basePath}/login`, data);
    return response.data;
}


export const fetchVerify = async (): Promise<TVerifyResponse> => {
    const response = await globalApi.get<TVerifyResponse>(`${basePath}/verify`);
    return response.data;
}