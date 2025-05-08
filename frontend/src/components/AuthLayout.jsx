import React from 'react';
import { Outlet } from 'react-router-dom';

export default function AuthLayout() {
  return (
    <div className="min-h-screen min-w-full bg-gray-50">
      <main className="flex items-center justify-center min-h-screen">
        <Outlet />
      </main>
    </div>
  );
}