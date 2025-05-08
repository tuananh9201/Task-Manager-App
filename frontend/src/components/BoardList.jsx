import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Plus, ChevronLeft, ChevronRight } from 'lucide-react';
import ModalCreateBoard from './ModalCreateBoard';

const BoardList = () => {
  const navigate = useNavigate();
  const [boards, setBoards] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalItems, setTotalItems] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const fetchBoards = async (page) => {
    setLoading(true);
    setError(null);
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        setError('Authentication token not found.');
        setLoading(false);
        return;
      }

      const response = await fetch(`http://localhost:8080/api/boards?page=${page}&limit=10`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      setBoards(data.boards || []);
      setCurrentPage(data.current_page || 1);
      setTotalPages(data.total_pages || 1);
      setTotalItems(data.total_items || 0);
    } catch (err) {
      setError(err.message);
      setBoards([]); // Clear boards on error
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchBoards(currentPage);
  }, [currentPage]);

  const handlePreviousPage = () => {
    if (currentPage > 1) {
      setCurrentPage(currentPage - 1);
    }
  };

  const handleNextPage = () => {
    if (currentPage < totalPages) {
      setCurrentPage(currentPage + 1);
    }
  };

  if (loading) {
    return <div className="text-center p-4">Loading boards...</div>;
  }

  if (error) {
    return <div className="text-center p-4 text-red-500">Error: {error}</div>;
  }

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-semibold text-gray-800">My Boards</h1>
        <button 
          className="flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          onClick={() => setIsModalOpen(true)}
        >
          <Plus className="h-5 w-5 mr-2" />
          Create New Board
        </button>
      </div>

      {boards.length === 0 && !loading && (
        <div className="text-center text-gray-500 p-4 bg-white rounded-lg shadow">
          No boards found. Create one to get started!
        </div>
      )}

      {boards.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {boards.map((board) => (
            <div 
              key={board.id} 
              className="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow cursor-pointer"
              onClick={() => navigate(`/boards/${board.id}/tasks`)}
            >
              <h2 className="text-xl font-semibold text-blue-700 mb-2">{board.name}</h2>
              <p className="text-gray-600 mb-3 h-12 overflow-hidden text-ellipsis">{board.description || 'No description available.'}</p>
              <div className="text-sm text-gray-500">
                <p>Created by: User ID {board.created_by}</p>
                <p>Created at: {new Date(board.created_at).toLocaleDateString()}</p>
              </div>
            </div>
          ))}
        </div>
      )}

      {totalPages > 1 && (
        <div className="mt-8 flex justify-center items-center space-x-4">
          <button
            onClick={handlePreviousPage}
            disabled={currentPage === 1}
            className="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
          >
            <ChevronLeft className="h-5 w-5 mr-1" />
            Previous
          </button>
          <span className="text-gray-700">
            Page {currentPage} of {totalPages}
          </span>
          <button
            onClick={handleNextPage}
            disabled={currentPage === totalPages}
            className="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 disabled:opacity-50 disabled:cursor-not-allowed flex items-center"
          >
            Next
            <ChevronRight className="h-5 w-5 ml-1" />
          </button>
        </div>
      )}
      <div className="mt-2 text-center text-sm text-gray-500">
        Total Boards: {totalItems}
      </div>
      <ModalCreateBoard
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onError={setError}
      />
    </div>
  );
};

export default BoardList;
