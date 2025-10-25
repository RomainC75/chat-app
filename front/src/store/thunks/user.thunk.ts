import { createAsyncThunk } from "@reduxjs/toolkit";
import { fetchSignupUser } from "../../api/user.api";
import type { TSignupUser } from "../../types/user.type";


export const signupUser = createAsyncThunk<string, { signupUser: TSignupUser }>(
    "/user/signup",
    async ({ signupUser }) => {
        return fetchSignupUser(signupUser);
    }
);