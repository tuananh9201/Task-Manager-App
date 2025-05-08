export default function ProgressOverview() {
  return (
    <div className="bg-white rounded-xl shadow-sm p-6 mb-8">
      <h2 className="text-lg font-semibold text-gray-800 mb-4">Progress Overview</h2>
      <div className="w-full bg-gray-200 rounded-full h-4 mb-6">
        <div className="bg-blue-600 h-4 rounded-full" style={{ width: '70%' }}></div>
      </div>
      <div className="flex justify-between text-sm">
        <div className="flex items-center">
          <div className="h-3 w-3 rounded-full bg-green-500 mr-2"></div>
          <span className="text-gray-600">Completed (12)</span>
        </div>
        <div className="flex items-center">
          <div className="h-3 w-3 rounded-full bg-yellow-500 mr-2"></div>
          <span className="text-gray-600">In Progress (5)</span>
        </div>
        <div className="flex items-center">
          <div className="h-3 w-3 rounded-full bg-red-500 mr-2"></div>
          <span className="text-gray-600">Not Started (3)</span>
        </div>
      </div>
    </div>
  );
}