import type { IAuthGateway } from "../../../core-logic/auth/interfaces/IAuthGateway";
import type { TLogin, TSignup } from "../../../core-logic/auth/types/auth.type";
import type { TVerifyResponse } from "../../../types/user.type";
import type { AuthApiLoader } from "./loaders.api";

export class HttpAuthGateway implements IAuthGateway {
  constructor(private readonly authLoader: AuthApiLoader) {}

  login(email: string, pass: string): Promise<TLogin> {
    return this.authLoader.login(email, pass);
  }

  signup(email: string, pass: string): Promise<TSignup> {
    return this.authLoader.signup(email, pass);
  }

  async verify(): Promise<TVerifyResponse> {
    return this.authLoader.verify();
  }
}
