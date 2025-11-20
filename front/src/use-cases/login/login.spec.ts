import { HttpAuthGateway } from "../../adapters/secondary/gateways/httpAuthGateway";
import { FakeAAuthApiLoader } from "../../adapters/secondary/loaders/fakeAuthLoader";
import type { AppState } from "../../store/appState";
import { initReduxStore, type ReduxStore } from "../../store/store";
import { login } from "./login";

describe("login use-case", () => {
  let store: ReduxStore;
  let authGateway: HttpAuthGateway;
  let fakeAuthApiLoader: FakeAAuthApiLoader;

  beforeEach(() => {
    fakeAuthApiLoader = new FakeAAuthApiLoader();
    authGateway = new HttpAuthGateway(fakeAuthApiLoader);
    store = initReduxStore({
      gateways: {
        authGateway,
      },
    });
  });

  it("should login", async () => {
    const email = "bob@email.com";
    const id = "123";
    const token = "TOKEN";
    fakeAuthApiLoader.expectedEmail = email;
    fakeAuthApiLoader.expectedId = id;
    fakeAuthApiLoader.expectedToken = token;

    await store.dispatch(login("john", "pass"));
    expect(store.getState().authManagement).toEqual<AppState["authManagement"]>(
      { data: { email, id, token }, error: null }
    );
  });

  it("should raise an error if the email/pass is wrong", async () => {
    fakeAuthApiLoader.shouldRaiseAnError = true;
    await store.dispatch(login("john", "pass"));
    expect(store.getState().authManagement).toEqual<AppState["authManagement"]>(
      { data: { email: "", id: "", token: "" }, error: "wrong email/password" }
    );
  });
});
