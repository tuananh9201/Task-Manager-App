import React from 'react';
import { Outlet } from 'react-router-dom';
import { CheckCircle, Clock, Layers } from 'lucide-react';

export default function AppLayout() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-50">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <div className="flex items-center">
              <h1 className="text-2xl font-bold text-blue-600">TaskManager</h1>
            </div>
            <div className="flex items-center space-x-4">
              <nav className="flex space-x-1">
                <a
                  href="/"
                  className="px-4 py-2 rounded-lg transition-colors bg-blue-100 text-blue-800"
                >
                  Dashboard
                </a>
                <a
                  href="/boards"
                  className="px-4 py-2 rounded-lg transition-colors text-gray-600 hover:bg-gray-100"
                >
                  Boards
                </a>
              </nav>
              <button className="ml-2 px-4 py-2 bg-red-50 text-red-600 font-medium rounded-lg hover:bg-red-100 transition-colors">
                Logout
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <Outlet />
      </main>
    </div>
  );
}