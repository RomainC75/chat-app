import type { IAuthGateway } from "../../../core-logic/auth/interfaces/IAuthGateway";
import type { TLogin, TSignup } from "../../../core-logic/auth/types/auth.type";
import type { TVerifyResponse } from "../../../types/user.type";


export class FakeAuthGateway implements IAuthGateway {
    expectedId: string = "";
    expectedEmail: string = "";
    expectedToken: string = "";

    constructor(){}

      login(email: string, pass: string): Promise<TLogin>{
        return Promise.resolve({
            id: this.expectedId,
            email: this.expectedEmail,
            token: this.expectedToken
        })
      }
      signup(email: string, pass: string): Promise<TSignup>{
        return Promise.resolve({
            id: this.expectedId,
            email: this.expectedEmail
        })
      }
      verify(): Promise<TVerifyResponse>{
        return Promise.resolve({
            id: this.expectedId,
            email: this.expectedEmail
        })
      }
}