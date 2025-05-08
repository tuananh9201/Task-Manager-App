import React, {useState} from 'react';
import { CheckCircle, Circle, MoreHorizontal } from 'lucide-react';
import { motion } from 'framer-motion';
import axios from 'axios';

const TaskItem = ({ task, onStatusUpdate, onEditTask, onDeleteTask }) => {
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const handleTaskComplete = async () => {
    try {
      const token = localStorage.getItem('token');
      const newStatus = task.status === 'completed' ? 'pending' : 'completed';
      await axios.patch(`http://localhost:8080/api/tasks/${task.id}`, {
        status: newStatus
      }, {
        headers: { Authorization: token }
      });
      onStatusUpdate(task.id, newStatus);
    } catch (err) {
      console.error('Failed to update task status:', err);
    }
  };
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, x: -10 }}
      transition={{ duration: 0.2 }}
      className="bg-white p-4 rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow"
    >
      <div className="flex items-start gap-3">
        <button
          onClick={handleTaskComplete}
          className="mt-1 text-gray-400 hover:text-green-500 transition-colors"
        >
          {task.status === 'completed' ? (
            <CheckCircle className="h-5 w-5 text-green-500" />
          ) : (
            <Circle className="h-5 w-5 text-gray-400" />
          )}
        </button>
        
        <div className="flex-1">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-semibold text-gray-800">{task.title}</h3>
            <div className="relative">
            <button
            onClick={() => {
              setIsPopoverOpen(!isPopoverOpen)
            }}
            className="mt-1 text-gray-400 hover:text-blue-500 transition-colors"
            aria-label="More options"
          >
            <MoreHorizontal className="h-5 w-5" />
          </button>
          
          {isPopoverOpen && (
            <div className="absolute right-0 mt-2 w-40 bg-white border border-gray-200 rounded-md shadow-lg z-10">
              <button
                onClick={() => {
                  setIsPopoverOpen(false);
                  onEditTask(task);
                }}
                className="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
              >
                Update
              </button>
              <button
                onClick={() => {
                  setIsPopoverOpen(false);
                  onDeleteTask(task.id);
                }}
                className="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50"
              >
                Delete
              </button>
            </div>
          )}
            </div>
          </div>
          <p className="text-gray-600 mt-1">{task.description}</p>
          <div className="flex items-center mt-2 text-sm text-gray-500 gap-4">
            <span className="flex items-center gap-1">
              {task.status === 'completed' ? 'Completed' : 'Pending'}
            </span>
            <span>Due: {task.due_date || 'No deadline'}</span>
          </div>
        </div>
      </div>
    </motion.div>
  );
};

export default TaskItem;
