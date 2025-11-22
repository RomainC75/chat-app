import { HttpAuthGateway } from "../../../adapters/secondary/gateways/httpAuthGateway";
import { FakeAAuthApiLoader } from "../../../adapters/secondary/loaders/fakeAuthLoader";
import type { AppState } from "../../../store/appState";
import { initReduxStore, type ReduxStore } from "../../../store/store";
import { signup } from "./signup";

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

  it("should signup", async () => {
    const email = "bob@email.com";
    const id = "123";
    fakeAuthApiLoader.expectedEmail = email;
    fakeAuthApiLoader.expectedId = id;

    await store.dispatch(signup("email@email.com", "pass"));
    expect(store.getState().authManagement).toEqual<AppState["authManagement"]>(
      { data: null, error: null, isLoading: false }
    );
  });

  it("should raise an error if the email/pass is wrong", async () => {
    fakeAuthApiLoader.shouldRaiseAnError = true;
    
    await expect(store.dispatch(signup("bad", "bad")))
      .rejects.toThrow("email already used");
    expect(store.getState().authManagement).toEqual<AppState["authManagement"]>(
      { data: null, error: "email already used", isLoading: false }
    );
  });

  it("should be in loading mode when trying to signup", async () => {
    
    store.dispatch(signup("john", "pass"));
    
    expect(store.getState().authManagement).toEqual<AppState["authManagement"]>(
      { data: null, error: null, isLoading: true }
    );

  });
});
