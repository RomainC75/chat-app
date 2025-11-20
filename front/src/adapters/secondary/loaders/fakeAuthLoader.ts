import type { TLogin, TSignup } from "../../../core-logic/auth/types/auth.type";
import type { TVerifyResponse } from "../../../types/user.type";
import type { AuthApiLoader } from "../gateways/loaders.api";

export class FakeAAuthApiLoader implements AuthApiLoader {
        expectedId: string = "";
        expectedEmail: string = "";
        expectedToken: string = "";
        shouldRaiseAnError: boolean= false;

        async login(email: string, password: string): Promise<TLogin>{
            if (this.shouldRaiseAnError){
                throw new Error("wrong email/password")
            }
            return Promise.resolve({
                id: this.expectedId,
                email: this.expectedEmail,
                token: this.expectedToken
            })
        }
        async signup(email: string, password: string): Promise<TSignup>{
            if (this.shouldRaiseAnError){
                throw new Error("wrong email/password")
            }
            return Promise.resolve({
                id: this.expectedId,
                email: this.expectedEmail
            })
        }
        
        async verify(): Promise<TVerifyResponse>{
            if (this.shouldRaiseAnError){
                throw new Error("wrong token")
            }
            return Promise.resolve({
                id: this.expectedId,
                email: this.expectedEmail
            })
        }
}