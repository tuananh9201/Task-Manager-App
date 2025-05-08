import React, { useState } from 'react';
import axios from 'axios';

const ModalCreateTask = ({ isOpen, onClose, boardId, onError }) => {
  const [newTask, setNewTask] = useState({
    title: '',
    description: '',
    due_date: new Date().toISOString().split('T')[0],
    status: 'pending'
  });

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const token = localStorage.getItem('token');
      await axios.post('http://localhost:8080/api/tasks', {
        ...newTask,
        board_id: parseInt(boardId)
      }, {
        headers: { Authorization: token }
      });
      setNewTask({
        title: '',
        description: '',
        due_date: new Date().toISOString().split('T')[0],
        status: 'pending'
      });
      onClose();
    } catch (err) {
      onError('Failed to create task');
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded-lg shadow-lg w-full max-w-md">
        <h3 className="text-xl font-bold mb-4">Create New Task</h3>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-gray-700 mb-2">Title</label>
            <input
              type="text"
              value={newTask.title}
              onChange={(e) => setNewTask({...newTask, title: e.target.value})}
              className="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>
          <div className="mb-4">
            <label className="block text-gray-700 mb-2">Description</label>
            <textarea
              value={newTask.description}
              onChange={(e) => setNewTask({...newTask, description: e.target.value})}
              className="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              rows="3"
            />
          </div>
          <div className="mb-4">
            <label className="block text-gray-700 mb-2">Due Date</label>
            <input
              type="date"
              value={newTask.due_date}
              onChange={(e) => setNewTask({...newTask, due_date: e.target.value})}
              className="w-full p-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div className="flex justify-end space-x-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
            >
              Create Task
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ModalCreateTask;