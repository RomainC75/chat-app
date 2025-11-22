import { render, screen } from "@testing-library/react";
import { BrowserRouter as Router } from "react-router-dom";
import * as router from 'react-router-dom';
import { userEvent } from "@testing-library/user-event";
import { initReduxStore, type ReduxStore } from "../../../../../store/store";
import { HttpAuthGateway } from "../../../../secondary/gateways/httpAuthGateway";
import { FakeAAuthApiLoader } from "../../../../secondary/loaders/fakeAuthLoader";
import Signup from "./Signup";
import { Provider } from "react-redux";
import { vi } from 'vitest';
import type { AppState } from "../../../../../store/appState";

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: vi.fn(),
  };
});

describe("", () => {
  let store: ReduxStore;
  let authGateway: HttpAuthGateway;
  let authLoader: FakeAAuthApiLoader;

  beforeEach(() => {
    authLoader = new FakeAAuthApiLoader();
    authGateway = new HttpAuthGateway(authLoader);
    store = initReduxStore({
      gateways: { authGateway },
    });
  });

  it("should signup, then navigate to login", async()=>{
    authLoader.shouldRaiseAnError=false
    const user = userEvent.setup();
    const mockNavigate = vi.fn();
    
    vi.spyOn(router, 'useNavigate').mockReturnValue(mockNavigate);
    renderSignup();

    const emailInput = screen.getByLabelText(/email address/i);
    await userEvent.type(emailInput, 'test@example.com');

    const passwordInput = screen.getByLabelText(/New Password/i);
    await userEvent.type(passwordInput, 'pass');
    
    const confirmPasswordInput = screen.getByLabelText(/confirm password/i);
    await userEvent.type(confirmPasswordInput, 'pass');

    const submitButton = screen.getByRole('confirm');
    await user.click(submitButton);

    expect(store.getState().authManagement).toEqual<AppState["authManagement"]>(
          { data: null, error: null, isLoading: false }
        );

    expect(mockNavigate).toHaveBeenCalledWith('/login');
  })

  it("should not signup then display an error message", async()=>{
    authLoader.shouldRaiseAnError=true
    const user = userEvent.setup();
    const mockNavigate = vi.fn();
    
    vi.spyOn(router, 'useNavigate').mockReturnValue(mockNavigate);
    renderSignup();

    const emailInput = screen.getByLabelText(/email address/i);
    await userEvent.type(emailInput, 'test@example.com');

    const passwordInput = screen.getByLabelText(/New Password/i);
    await userEvent.type(passwordInput, 'pass');
    
    const confirmPasswordInput = screen.getByLabelText(/confirm password/i);
    await userEvent.type(confirmPasswordInput, 'pass');

    const submitButton = screen.getByRole('confirm');
    await user.click(submitButton);

    expect(store.getState().authManagement).toEqual<AppState["authManagement"]>(
          { data: null, error: "email already used", isLoading: false }
    );
    expect(mockNavigate).not.toHaveBeenCalled();
    
    expect(screen.getByText(/email already used/i)).toBeInTheDocument()
  })

  const renderSignup = () => {
    render(
      <Provider store={store}>
        <Router>
        <Signup/>
        </Router>
      </Provider>,
    );
  };

});

