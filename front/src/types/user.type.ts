import type { TDefaultSlice } from "./store-slice";

export type TUserSlice = TDefaultSlice & {
    user: TUser | null;
}

export type TUser = TSignupUser & {
    id: string;
}

export type TSignupUser = {
    email: string;
    password: string;
}