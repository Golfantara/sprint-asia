import { useEffect, useState } from "react";

function App() {
  const [tasks, setTasks] = useState([]);
  const [completedTasks, setCompletedTasks] = useState([]);
  const [activeTab, setActiveTab] = useState("ongoing");
  const [form, setForm] = useState({
    id: null,
    user_id: 1,
    title: "",
    description: "",
    deadline: "",
  });
  const [loading, setLoading] = useState(false);
  const [editing, setEditing] = useState(false);
  const [expandedTask, setExpandedTask] = useState(null);
  const [subtaskForm, setSubtaskForm] = useState({ task_id: null, title: "" });
  const [showSubtaskModal, setShowSubtaskModal] = useState(false);

  const API_URL = "http://localhost:8000/tasks";

  const fetchTasks = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${API_URL}?page=1&size=100`);
      const json = await res.json();
      setTasks(json.data || []);
    } catch (err) {
      console.error(err);
    }
    setLoading(false);
  };

  const fetchCompletedTasks = async () => {
    setLoading(true);
    try {
      const res = await fetch(`${API_URL}/history?page=1&size=100`);
      const json = await res.json();
      setCompletedTasks(json.data || []);
    } catch (err) {
      console.error(err);
    }
    setLoading(false);
  };

  const fetchTaskDetails = async (taskId) => {
    try {
      const res = await fetch(`${API_URL}/${taskId}`);
      const json = await res.json();
      return json.data || null;
    } catch (err) {
      console.error(err);
      return null;
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    // ✅ Ensure full ISO format for backend
    let deadline = form.deadline ? new Date(form.deadline).toISOString() : null;

    try {
      if (editing) {
        await fetch(`${API_URL}/${form.id}`, {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            title: form.title,
            description: form.description,
            deadline,
          }),
        });
      } else {
        await fetch(API_URL, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            ...form,
            deadline,
          }),
        });
      }
      resetForm();
      fetchTasks();
    } catch (err) {
      console.error(err);
      alert("Failed to save task");
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm("Are you sure you want to delete this task?")) return;
    try {
      await fetch(`${API_URL}/${id}`, { method: "DELETE" });
      fetchTasks();
    } catch (err) {
      console.error(err);
      alert("Failed to delete task");
    }
  };

  const handleEdit = (task) => {
    setForm({
      id: task.id,
      user_id: task.user_id,
      title: task.title,
      description: task.description,
      deadline: task.deadline ? task.deadline.slice(0, 16) : "",
    });
    setEditing(true);
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  const handleCompleteTask = async (taskId) => {
    try {
      await fetch(`${API_URL}/${taskId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ is_completed: true }),
      });
      fetchTasks();
      if (activeTab === "completed") fetchCompletedTasks();
    } catch (err) {
      console.error(err);
      alert("Failed to complete task");
    }
  };

  const resetForm = () => {
    setForm({ id: null, user_id: 1, title: "", description: "", deadline: "" });
    setEditing(false);
  };

  const toggleExpandTask = async (taskId) => {
    if (expandedTask === taskId) {
      setExpandedTask(null);
    } else {
      const details = await fetchTaskDetails(taskId);
      if (details) {
        setTasks((prev) =>
          prev.map((t) => (t.id === taskId ? { ...t, ...details } : t))
        );
        setCompletedTasks((prev) =>
          prev.map((t) => (t.id === taskId ? { ...t, ...details } : t))
        );
      }
      setExpandedTask(taskId);
    }
  };

  const handleAddSubtask = (taskId) => {
    setSubtaskForm({ task_id: taskId, title: "" });
    setShowSubtaskModal(true);
  };

  const handleSubmitSubtask = async (e) => {
    e.preventDefault();
    try {
      await fetch(`${API_URL}/subtask`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(subtaskForm),
      });
      setShowSubtaskModal(false);
      setSubtaskForm({ task_id: null, title: "" });
      await toggleExpandTask(subtaskForm.task_id);
      fetchTasks();
    } catch (err) {
      console.error(err);
      alert("Failed to add subtask");
    }
  };

  const handleToggleSubtask = async (subtask, taskId) => {
    try {
      await fetch(`${API_URL}/subtask/${subtask.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          task_id: taskId,
          title: subtask.title,
          is_completed: !subtask.is_completed,
        }),
      });
      await toggleExpandTask(taskId);
      fetchTasks();
    } catch (err) {
      console.error(err);
      alert("Failed to update subtask");
    }
  };

  const handleDeleteSubtask = async (subtaskId, taskId) => {
    if (!window.confirm("Delete this subtask?")) return;
    try {
      await fetch(`${API_URL}/subtask/${subtaskId}`, { method: "DELETE" });
      await toggleExpandTask(taskId);
      fetchTasks();
    } catch (err) {
      console.error(err);
      alert("Failed to delete subtask");
    }
  };

  const isDue = (task) => task.is_due || false;

  const getProgressPercentage = (task) => {
    if (task.progress !== undefined && task.progress !== null)
      return task.progress;
    if (!task.subtasks || task.subtasks.length === 0) return 0;
    const completed = task.subtasks.filter((st) => st.is_completed).length;
    return Math.round((completed / task.subtasks.length) * 100);
  };

  useEffect(() => {
    fetchTasks();
  }, []);

  useEffect(() => {
    if (activeTab === "completed") {
      fetchCompletedTasks();
    }
  }, [activeTab]);

  const displayTasks = activeTab === "ongoing" ? tasks : completedTasks;

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-4 md:p-8">
      <div className="max-w-6xl mx-auto">
        {/* Tabs */}
        <div className="flex items-center justify-between mb-8 flex-wrap gap-4">
          <h1 className="text-4xl font-bold text-indigo-900">
            📋 Task Manager
          </h1>
          <div className="flex gap-2 bg-white rounded-lg p-1 shadow-sm">
            <button
              onClick={() => setActiveTab("ongoing")}
              className={`px-4 py-2 rounded-md font-medium transition-all ${
                activeTab === "ongoing"
                  ? "bg-indigo-600 text-white"
                  : "text-gray-600 hover:bg-gray-100"
              }`}
            >
              Ongoing ({tasks.length})
            </button>
            <button
              onClick={() => setActiveTab("completed")}
              className={`px-4 py-2 rounded-md font-medium transition-all ${
                activeTab === "completed"
                  ? "bg-indigo-600 text-white"
                  : "text-gray-600 hover:bg-gray-100"
              }`}
            >
              Completed ({completedTasks.length})
            </button>
          </div>
        </div>

        {/* Form */}
        {activeTab === "ongoing" && (
          <div className="bg-white shadow-lg rounded-xl p-6 mb-8 border border-gray-200">
            <h2 className="text-xl font-semibold mb-4 text-gray-800">
              {editing ? "✏️ Edit Task" : "➕ Create New Task"}
            </h2>
            <div className="grid gap-4">
              <input
                type="text"
                placeholder="Task Title *"
                className="border border-gray-300 p-3 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
                value={form.title}
                onChange={(e) => setForm({ ...form, title: e.target.value })}
              />
              <textarea
                placeholder="Task Description"
                className="border border-gray-300 p-3 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
                rows="3"
                value={form.description}
                onChange={(e) =>
                  setForm({ ...form, description: e.target.value })
                }
              />
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Deadline
                </label>
                <input
                  type="datetime-local"
                  className="border border-gray-300 p-3 rounded-lg w-full focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
                  value={form.deadline}
                  onChange={(e) =>
                    setForm({ ...form, deadline: e.target.value })
                  }
                />
              </div>
              <div className="flex gap-3">
                <button
                  onClick={handleSubmit}
                  className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-3 rounded-lg font-medium transition-colors shadow-md"
                >
                  {editing ? "💾 Update Task" : "➕ Add Task"}
                </button>
                {editing && (
                  <button
                    onClick={resetForm}
                    className="border-2 border-gray-300 text-gray-700 px-6 py-3 rounded-lg font-medium hover:bg-gray-50 transition-colors"
                  >
                    Cancel
                  </button>
                )}
              </div>
            </div>
          </div>
        )}

        {/* Task List */}
        {loading ? (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-4 border-indigo-600 border-t-transparent"></div>
            <p className="text-gray-600 mt-4">Loading tasks...</p>
          </div>
        ) : (
          <div className="grid gap-4">
            {displayTasks.length === 0 ? (
              <div className="bg-white rounded-xl p-12 text-center shadow-md">
                <div className="text-6xl mb-4">
                  {activeTab === "ongoing" ? "📝" : "✅"}
                </div>
                <p className="text-gray-500 text-lg">
                  {activeTab === "ongoing"
                    ? "No ongoing tasks. Create one to get started!"
                    : "No completed tasks yet."}
                </p>
              </div>
            ) : (
              displayTasks.map((task) => {
                const progress = getProgressPercentage(task);
                const overdue = isDue(task);
                const isExpanded = expandedTask === task.id;

                return (
                  <div
                    key={task.id}
                    className="bg-white shadow-md rounded-xl border border-gray-200 overflow-hidden hover:shadow-lg transition-shadow"
                  >
                    <div className="p-5">
                      <div className="flex items-start gap-4">
                        {activeTab === "ongoing" && (
                          <input
                            type="checkbox"
                            className="mt-1 h-5 w-5 text-indigo-600 rounded cursor-pointer"
                            onChange={() => handleCompleteTask(task.id)}
                          />
                        )}
                        <div className="flex-1">
                          <div className="flex items-start justify-between mb-2">
                            <div className="flex-1">
                              <h2 className="font-semibold text-xl text-gray-900 mb-1">
                                {task.title}
                              </h2>
                              {task.description && (
                                <p className="text-gray-600 text-sm mb-2">
                                  {task.description}
                                </p>
                              )}
                            </div>
                          </div>

                          <div className="flex flex-wrap items-center gap-3 text-sm mb-2">
                            {task.deadline && (
                              <span
                                className={`flex items-center gap-1 px-3 py-1 rounded-full font-medium ${
                                  overdue
                                    ? "bg-red-100 text-red-700"
                                    : "bg-blue-100 text-blue-700"
                                }`}
                              >
                                ⏰{" "}
                                {new Date(task.deadline).toLocaleString(
                                  "en-US",
                                  {
                                    month: "short",
                                    day: "numeric",
                                    year: "numeric",
                                    hour: "2-digit",
                                    minute: "2-digit",
                                  }
                                )}
                                {overdue && " 🔴 OVERDUE!"}
                              </span>
                            )}
                            {((task.subtasks && task.subtasks.length > 0) ||
                              (task.progress > 0 && !task.subtasks)) && (
                              <span className="flex items-center gap-1 px-3 py-1 bg-purple-100 text-purple-700 rounded-full font-medium">
                                📊 {progress}% Complete
                              </span>
                            )}
                          </div>

                          {/* Progress Bar */}
                          {((task.subtasks && task.subtasks.length > 0) ||
                            (task.progress > 0 && !task.subtasks)) && (
                            <div className="mt-3 bg-gray-200 rounded-full h-2 overflow-hidden">
                              <div
                                className="bg-indigo-600 h-full transition-all duration-300"
                                style={{ width: `${progress}%` }}
                              ></div>
                            </div>
                          )}

                          {/* ✅ Always show Add Subtask button */}
                          {activeTab === "ongoing" && (
                            <div className="mt-3">
                              <button
                                onClick={() => handleAddSubtask(task.id)}
                                className="text-indigo-600 hover:text-indigo-800 text-sm font-medium"
                              >
                                + Add Subtask
                              </button>
                            </div>
                          )}

                          {/* Details */}
                          {isExpanded && task.subtasks && (
                            <div className="mt-4 pl-4 border-l-2 border-indigo-300">
                              <div className="font-medium text-gray-700 mb-2">
                                Subtasks ({task.subtasks.length})
                              </div>
                              <div className="space-y-2">
                                {task.subtasks.map((subtask) => (
                                  <div
                                    key={subtask.id}
                                    className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50"
                                  >
                                    <input
                                      type="checkbox"
                                      checked={subtask.is_completed}
                                      onChange={() =>
                                        handleToggleSubtask(subtask, task.id)
                                      }
                                      className="h-4 w-4 text-indigo-600 rounded cursor-pointer"
                                      disabled={activeTab === "completed"}
                                    />
                                    <span
                                      className={`flex-1 text-sm ${
                                        subtask.is_completed
                                          ? "line-through text-gray-400"
                                          : "text-gray-700"
                                      }`}
                                    >
                                      {subtask.title}
                                    </span>
                                    {activeTab === "ongoing" && (
                                      <button
                                        onClick={() =>
                                          handleDeleteSubtask(
                                            subtask.id,
                                            task.id
                                          )
                                        }
                                        className="text-red-500 hover:text-red-700 text-xs"
                                      >
                                        Delete
                                      </button>
                                    )}
                                  </div>
                                ))}
                              </div>
                            </div>
                          )}

                          {/* Actions */}
                          <div className="flex gap-3 mt-4 pt-3 border-t border-gray-200">
                            <button
                              onClick={() => toggleExpandTask(task.id)}
                              className="text-indigo-600 hover:text-indigo-800 font-medium text-sm"
                            >
                              {isExpanded ? "▲ Hide Details" : "▼ Show Details"}
                            </button>
                            {activeTab === "ongoing" && (
                              <>
                                <button
                                  onClick={() => handleEdit(task)}
                                  className="text-blue-600 hover:text-blue-800 font-medium text-sm"
                                >
                                  ✏️ Edit
                                </button>
                                <button
                                  onClick={() => handleDelete(task.id)}
                                  className="text-red-600 hover:text-red-800 font-medium text-sm"
                                >
                                  🗑️ Delete
                                </button>
                              </>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                );
              })
            )}
          </div>
        )}

        {/* Subtask Modal */}
        {showSubtaskModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-xl p-6 max-w-md w-full shadow-2xl">
              <h3 className="text-xl font-semibold mb-4">Add Subtask</h3>
              <input
                type="text"
                placeholder="Subtask title *"
                className="border border-gray-300 p-3 rounded-lg w-full mb-4 focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none"
                value={subtaskForm.title}
                onChange={(e) =>
                  setSubtaskForm({ ...subtaskForm, title: e.target.value })
                }
                onKeyPress={(e) => {
                  if (e.key === "Enter" && subtaskForm.title.trim()) {
                    handleSubmitSubtask(e);
                  }
                }}
              />
              <div className="flex gap-3">
                <button
                  onClick={handleSubmitSubtask}
                  disabled={!subtaskForm.title.trim()}
                  className="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded-lg font-medium disabled:opacity-50"
                >
                  Add
                </button>
                <button
                  onClick={() => setShowSubtaskModal(false)}
                  className="border border-gray-300 text-gray-700 px-6 py-2 rounded-lg font-medium hover:bg-gray-50"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
