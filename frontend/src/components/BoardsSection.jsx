import { Plus, Users, Share2, Edit3, Trash2 } from 'lucide-react';

export default function BoardsSection() {
  return (
    <div className="mb-8">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-xl font-semibold text-gray-800">Boards</h2>
        <button className="flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
          <Plus className="h-4 w-4 mr-1" />
          Add Board
        </button>
      </div>

      <div className="bg-white rounded-xl shadow-sm overflow-hidden">
        <div className="p-5 border-b border-gray-100 flex justify-between items-center">
          <h3 className="text-lg font-medium text-gray-800">Marketing Board</h3>
          <div className="flex space-x-2">
            <button className="p-2 text-blue-600 hover:bg-blue-50 rounded-md transition-colors">
              <Users className="h-5 w-5" />
            </button>
            <button className="px-3 py-1 text-blue-600 bg-blue-50 rounded-md hover:bg-blue-100 transition-colors text-sm font-medium">
              Invite
            </button>
          </div>
        </div>
        
        <div className="p-5 border-b border-gray-100">
          <div className="flex justify-between items-center mb-4">
            <h4 className="font-medium text-gray-700">To Do</h4>
            <button className="p-1 text-green-600 hover:bg-green-50 rounded-full transition-colors">
              <Share2 className="h-4 w-4" />
            </button>
          </div>
          
          <div className="bg-gray-50 rounded-lg p-4 transform transition-all hover:shadow-md">
            <div className="flex justify-between items-start">
              <div>
                <h5 className="font-medium text-gray-800">Finish user onboarding</h5>
                <p className="text-sm text-gray-500 mt-1">Due: 2025-05-10</p>
              </div>
              <div className="flex space-x-1">
                <button className="p-1 text-blue-600 hover:bg-blue-50 rounded transition-colors">
                  <Edit3 className="h-4 w-4" />
                </button>
                <button className="p-1 text-red-600 hover:bg-red-50 rounded transition-colors">
                  <Trash2 className="h-4 w-4" />
                </button>
              </div>
            </div>
          </div>
        </div>
        
        <div className="p-4 bg-gray-50">
          <button className="w-full py-2 flex items-center justify-center text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors">
            <Plus className="h-4 w-4 mr-1" />
            Add Task
          </button>
        </div>
      </div>
    </div>
  );
}