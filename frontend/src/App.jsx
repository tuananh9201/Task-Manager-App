import './App.css'
import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Login from './components/Login';
import TaskList from './components/TaskList';
import SignupForm from './components/SignupForm';
import BoardList from './components/BoardList';

import AppLayout from './components/AppLayout';
import AuthLayout from './components/AuthLayout';

const App = () => {
  return (
    <Router>
      <Routes>
        <Route element={<AuthLayout />}>
          <Route path="/login" element={<Login />} />
          <Route path="/signup" element={<SignupForm />} />
          <Route path="/" element={<Login />} />
        </Route>
        <Route element={<AppLayout />}>
          <Route path="/boards/:boardId/tasks" element={<TaskList />} />
          <Route path="/boards" element={<BoardList />} />
        </Route>
      </Routes>
    </Router>
  );
};

export default App;