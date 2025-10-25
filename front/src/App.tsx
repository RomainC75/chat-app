import Signup from "./components/Signup";
import "./App.css";
import { Navigate, Route, Routes, BrowserRouter as Router, } from "react-router-dom";
import Login from "./components/Login";
import ProtectedRoute from "./components/ProtectedRoute/ProtectedRoute";

function App() {
  return (
    <Router>
      <div className="app">
        <Routes>
          <Route
            path="/signup"
            element={
                <Signup />
            }
          />
          <Route
            path="/login"
            element={
                <Login />
            }
          /> 
          <Route
            path="/*"
            element={
              <ProtectedRoute>
                <div>
                  <h1>Application</h1>
                </div>
              </ProtectedRoute>
            }
          />

          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
