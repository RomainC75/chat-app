import type { TLogin, TSignup } from "../../../core-logic/auth/types/auth.type";
import type { TVerifyResponse } from "../../../types/user.type";

export interface AuthApiLoader {
    login(email: string, pass: string): Promise<TLogin>;
    signup(email: string, pass: string): Promise<TSignup>;
    verify(): Promise<TVerifyResponse>;
}