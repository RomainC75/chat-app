import { FakeAuthGateway } from "../../adapters/secondary/gateways/fakeAuthGateway";
import type { AppState } from "../../store/appState";
import { initReduxStore, type ReduxStore } from "../../store/store"
import { login } from "./login";

describe("login use-case",()=>{
    let store: ReduxStore;
    let authGateway: FakeAuthGateway;

    beforeEach(()=>{
        authGateway = new FakeAuthGateway();
        store = initReduxStore({gateways: {
            authGateway
        }})
    })

    it('should login', async ()=>{
        const email = "bob@email.com"
        const id = "123"
        const token = "TOKEN"
        authGateway.expectedEmail= email
        authGateway.expectedId = id
        authGateway.expectedToken = token

        await store.dispatch(login("john", "pass"));
        console.log("-> store.getState().authManagement")
        expect(store.getState().authManagement).toEqual<
      AppState["authManagement"]
    >({ data: {email, id, token}, error: null });
    })
})