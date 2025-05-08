import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useParams } from 'react-router-dom';

const TaskList = () => {
  const { boardId } = useParams();
  const [tasks, setTasks] = useState([]);
  const [error, setError] = useState('');

  useEffect(() => {
    // Fetch initial tasks
    const fetchTasks = async () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          setError('Please log in');
          return;
        }
        const response = await axios.get(`http://localhost:8080/api/tasks?board_id=${boardId}`, {
          headers: { Authorization: token },
        });
        setTasks(response.data);
      } catch (err) {
        setError('Failed to fetch tasks');
      }
    };

    fetchTasks();

    // Set up WebSocket for real-time updates
    const ws = new WebSocket('ws://localhost:8080/ws');
    ws.onmessage = (event) => {
      const task = JSON.parse(event.data);
      // Ensure boardId from params is treated as a number for comparison
      if (task.board_id === parseInt(boardId)) { 
        setTasks((prev) => [...prev, task]);
      }
    };
    ws.onerror = () => setError('WebSocket connection failed');
    return () => ws.close();
  }, [boardId]);

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <h2 className="text-2xl font-bold mb-4">Tasks</h2>
      {error && <p className="text-red-500 mb-4">{error}</p>}
      {tasks.length === 0 ? (
        <p className="text-gray-500">No tasks found.</p>
      ) : (
        <div className="grid gap-4">
          {tasks.map((task) => (
            <div key={task.id} className="border p-4 rounded-lg shadow-sm">
              <h3 className="text-lg font-semibold">{task.title}</h3>
              <p className="text-gray-600">{task.description}</p>
              <p className="text-sm text-gray-500">Status: {task.status}</p>
              <p className="text-sm text-gray-500">Due: {task.due_date || 'N/A'}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default TaskList;