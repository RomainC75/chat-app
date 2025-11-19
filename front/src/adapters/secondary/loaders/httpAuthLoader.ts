import type { TLogin, TSignup } from "../../../core-logic/auth/types/auth.type";
import type { TVerifyResponse } from "../../../types/user.type";
import type { AuthApiLoader } from "../gateways/loaders.api";
import openApi from "./open.api";
import secureApi from "./secure.api";


const basePath = "/auth"

export class HttpAuthApiLoader implements AuthApiLoader {
        async login(email: string, password: string): Promise<TLogin>{
            const data = {
                email,
                password,
            }
            const response = await openApi.post<TLogin>(`${basePath}/login`, data);
            return response.data;

        }
        async signup(email: string, password: string): Promise<TSignup>{
            const data = {
                email,
                password,
            }
            const response = await openApi.post<TSignup>(`${basePath}/signup`, data);
            return response.data;
        }
        
        async verify(): Promise<TVerifyResponse>{
            const response = await secureApi.get<TVerifyResponse>(`${basePath}/verify`);
            return response.data;
        }
}