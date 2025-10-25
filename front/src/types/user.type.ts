import type { TDefaultSlice } from "./store-slice";

export type TUserSlice = TDefaultSlice & {
    user: TUser | null;
    isConnected: boolean;
}

export type TUser = TSignupUser & {
    id: string;
}

export type TSignupUser = {
    email: string;
    password: string;
}

export type TLoginsUser = {
    email: string;
    password: string;
}

export type TLoginResponse = {
    id: number,
    token: string
}

export type TVerifyResponse =  {
    id: string;
    email: string;
}