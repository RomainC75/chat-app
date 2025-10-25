import { createAsyncThunk } from "@reduxjs/toolkit";
import { fetchLoginUser, fetchSignupUser, fetchVerify } from "../../api/user.api";
import type { TLoginResponse, TLoginsUser, TSignupUser, TVerifyResponse } from "../../types/user.type";


export const signupUser = createAsyncThunk<string, { signupUser: TSignupUser }>(
    "/user/signup",
    async ({ signupUser }) => {
        return fetchSignupUser(signupUser);
    }
);

export const loginUser = createAsyncThunk<TLoginResponse, { loginUser: TLoginsUser }>(
    "/user/login",
    async ({ loginUser }) => {
        return fetchLoginUser(loginUser);
    }
);


export const verify = createAsyncThunk<TVerifyResponse>(
    "/user/verify",
    async () => {
        return fetchVerify();
    }
);