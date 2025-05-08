import { CheckCircle, Clock, Layers } from 'lucide-react';

export default function StatsCards() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div className="bg-white rounded-xl shadow-sm p-6 border border-green-100 transform transition-all hover:shadow-md">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-gray-500">Completed Tasks</p>
            <h2 className="text-3xl font-bold text-green-600 mt-1">12</h2>
          </div>
          <div className="bg-green-100 p-3 rounded-full">
            <CheckCircle className="h-6 w-6 text-green-600" />
          </div>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm p-6 border border-yellow-100 transform transition-all hover:shadow-md">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-gray-500">Pending Tasks</p>
            <h2 className="text-3xl font-bold text-yellow-600 mt-1">5</h2>
          </div>
          <div className="bg-yellow-100 p-3 rounded-full">
            <Clock className="h-6 w-6 text-yellow-600" />
          </div>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm p-6 border border-blue-100 transform transition-all hover:shadow-md">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-gray-500">Boards</p>
            <h2 className="text-3xl font-bold text-blue-600 mt-1">3</h2>
          </div>
          <div className="bg-blue-100 p-3 rounded-full">
            <Layers className="h-6 w-6 text-blue-600" />
          </div>
        </div>
      </div>
    </div>
  );
}