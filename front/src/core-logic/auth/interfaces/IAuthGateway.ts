import type { TVerifyResponse } from "../../../types/user.type";
import type { TLogin, TSignup } from "../types/auth.type";

export interface IAuthGateway {
  login(email: string, pass: string): Promise<TLogin>;
  signup(email: string, pass: string): Promise<TSignup>;
  verify(): Promise<TVerifyResponse>;
}