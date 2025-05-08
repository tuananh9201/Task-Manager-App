import { useState } from 'react';
import { CheckCircle, Clock, Layers, Edit3, Trash2, Share2, Users, Plus } from 'lucide-react';
import BoardList from '../components/BoardList'; // Added import for BoardList

export default function TaskManager() {
  const [activeTab, setActiveTab] = useState('dashboard');

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
                <button
                  onClick={() => setActiveTab('dashboard')}
                  className={`px-4 py-2 rounded-lg transition-colors ${
                    activeTab === 'dashboard'
                      ? 'bg-blue-100 text-blue-800'
                      : 'text-gray-600 hover:bg-gray-100'
                  }`}
                >
                  Dashboard
                </button>
                <button
                  onClick={() => setActiveTab('boards')}
                  className={`px-4 py-2 rounded-lg transition-colors ${
                    activeTab === 'boards'
                      ? 'bg-blue-100 text-blue-800'
                      : 'text-gray-600 hover:bg-gray-100'
                  }`}
                >
                  Boards
                </button>
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
        {activeTab === 'dashboard' && (
          <>
            {/* These components are placeholders, replace with actual implementations */}
            {/* <StatsCards /> */}
            {/* <ProgressOverview /> */}
            {/* <BoardsSection /> */}
            <div className="text-center p-10 bg-white rounded-lg shadow">
              <h2 className="text-xl font-semibold text-gray-700">Dashboard Content</h2>
              <p className="text-gray-500 mt-2">Stats, overview, and featured boards will be shown here.</p>
            </div>
          </>
        )}
        {activeTab === 'boards' && <BoardList />}
      </main>
    </div>
  );
}