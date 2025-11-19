
import type { IAuthGateway } from "../../core-logic/auth/interfaces/IAuthGateway";
import { HttpAuthGateway } from "./gateways/httpAuthGateway";
import { HttpAuthApiLoader } from "./loaders/httpAuthLoader";


export type Gateways = {
  authGateway: IAuthGateway;
};


const authGateway = new HttpAuthGateway(
  new HttpAuthApiLoader(),
);

export const gateways: Gateways = {
    authGateway,
};
