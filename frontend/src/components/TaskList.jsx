import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useParams } from 'react-router-dom';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';
import { AnimatePresence } from 'framer-motion';
import TaskItem from './TaskItem';
import ModalCreateTask from './ModalCreateTask';
import ModalEditTask from './ModalEditTask';

const TaskList = () => {
  const { boardId } = useParams();
  const [tasks, setTasks] = useState([]);
  const [error, setError] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [editTaskId, setEditTaskId] = useState(null);
  const [editTitle, setEditTitle] = useState('');
  const [editDescription, setEditDescription] = useState('');
  const [editDueDate, setEditDueDate] = useState('');

  const updateTaskStatus = (taskId, newStatus) => {
    setTasks(tasks.map(task =>
      task.id === taskId ? { ...task, status: newStatus } : task
    ));
  };

  useEffect(() => {
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
        setTasks(response.data.tasks);
      } catch (err) {
        setError('Failed to fetch tasks');
      }
    };

    fetchTasks();

    const ws = new WebSocket('ws://localhost:8080/ws');
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);

      if (data.Type === 'create' && data.Task.board_id === parseInt(boardId)) {
        setTasks((prev) => [...prev, data.Task]);
      } else if (data.Type === 'update' && data.Task.board_id === parseInt(boardId)) {
        setTasks((prev) =>
          prev.map((task) =>
            task.id === data.Task.id ? data.Task : task
          )
        );
      } else if (data.Type === 'delete') {

        setTasks((prev) =>
          prev.filter((task) => task.id !== data.Task.id)
        );
      }
    };
    ws.onerror = () => setError('WebSocket connection failed');
    ws.onclose = () => setError('WebSocket connection closed');
    ws.onopen = () => setError('');
    return () => ws.close();
  }, [boardId]);

  const onEditTask = (task) => {

    setEditTaskId(task.id);
    setEditTitle(task.title);
    setEditDescription(task.description);
    setEditDueDate(task.due_date || '');
    setIsEditModalOpen(true);
  };

  const onDeleteTask = async (taskId) => {
    try {
      const token = localStorage.getItem('token');
      await axios.delete(`http://localhost:8080/api/tasks/${taskId}`, {
        headers: { Authorization: token },
      });
      setTasks(tasks.filter(task => task.id !== taskId));
    } catch (err) {
      setError('Failed to delete task');
    }
  };

  const onDragEnd = (result) => {
    if (!result.destination) {
      return;
    }

    const { source, destination } = result;

    if (source.droppableId !== destination.droppableId) {
      const sourceListId = source.droppableId;
      const destListId = destination.droppableId;
      const sourceIndex = source.index;
      const destIndex = destination.index;

      // Implement reordering logic here
      console.log(`Moved task from list ${sourceListId} at index ${sourceIndex} to list ${destListId} at index ${destIndex}`);
    }
  };

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">Tasks</h2>
        <button
          onClick={() => setIsModalOpen(true)}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
        >
          Add Task
        </button>
      </div>
      {error && <p className="text-red-500 mb-4">{error}</p>}
      {tasks.length === 0 ? (
        <p className="text-gray-500">No tasks found.</p>
      ) : (
        <DragDropContext onDragEnd={onDragEnd}>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {Object.entries(
              tasks.reduce((acc, task) => {
                if (!acc[task.list_id]) {
                  acc[task.list_id] = {
                    list: task.list_name,
                    tasks: []
                  };
                }
                acc[task.list_id].tasks.push(task);
                return acc;
              }, {})
            ).map(([listId, { list, tasks }]) => (
              <Droppable droppableId={listId} key={listId}>
                {(provided) => (
                  <div
                    {...provided.droppableProps}
                    ref={provided.innerRef}
                    className="bg-white p-4 rounded-lg shadow-sm border border-gray-200"
                  >
                    <h3 className="text-lg font-semibold text-gray-800 mb-3">{list}</h3>
                    <div className="space-y-3">
                      <AnimatePresence>
                        {tasks.map((task, index) => (
                          <Draggable key={task.id} draggableId={String(task.id)} index={index}>
                            {(provided) => (
                              <div
                                ref={provided.innerRef}
                                {...provided.draggableProps}
                                {...provided.dragHandleProps}
                              >
                                <TaskItem key={task.id} task={task} onStatusUpdate={updateTaskStatus} onEditTask={onEditTask} onDeleteTask={onDeleteTask} />
                              </div>
                            )}
                          </Draggable>
                        ))}
                      </AnimatePresence>
                      {provided.placeholder}
                    </div>
                  </div>
                )}
              </Droppable>
            ))}
          </div>
        </DragDropContext>
      )}
      <ModalCreateTask
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        boardId={boardId}
        onError={setError}
      />
      <ModalEditTask
        isOpen={isEditModalOpen}
        onClose={() => setIsEditModalOpen(false)}
        taskId={editTaskId}
        task={{
          title: editTitle,
          description: editDescription,
          due_date: editDueDate
        }}
        boardId={boardId}
        onError={setError}
      />
    </div>
  );
};

export default TaskList;
