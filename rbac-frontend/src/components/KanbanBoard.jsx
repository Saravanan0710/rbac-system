import { useState } from 'react'
import './KanbanBoard.css'

function KanbanBoard({ tasks = [], onTaskUpdate, onTaskDelete }) {
  const [draggedTask, setDraggedTask] = useState(null)

  const columns = {
    'todo': { title: 'To Do', color: '--accent-pink', tasks: [] },
    'in-progress': { title: 'In Progress', color: '--accent-cyan', tasks: [] },
    'review': { title: 'Review', color: '--accent-orange', tasks: [] },
    'done': { title: 'Done', color: '--accent-green', tasks: [] }
  }

  // Map uppercase status from backend to lowercase column keys
  const statusMap = {
    'TODO': 'todo',
    'IN_PROGRESS': 'in-progress',
    'IN-PROGRESS': 'in-progress',
    'REVIEW': 'review',
    'DONE': 'done',
    'ARCHIVED': 'done'
  }

  tasks.forEach((task) => {
    // Convert backend status (uppercase) to column key (lowercase)
    const dbStatus = task.status?.toUpperCase() || 'TODO'
    const columnKey = statusMap[dbStatus] || 'todo'
    
    if (columns[columnKey]) {
      columns[columnKey].tasks.push(task)
    }
  })

  const handleDragStart = (task) => {
    setDraggedTask(task)
  }

  const handleDragOver = (e) => {
    e.preventDefault()
  }

  const handleDrop = (status) => {
    if (draggedTask && onTaskUpdate) {
      // Convert status to uppercase format required by backend
      // todo -> TODO, in-progress -> IN_PROGRESS, review -> REVIEW, done -> DONE
      const statusMap = {
        'todo': 'TODO',
        'in-progress': 'IN_PROGRESS',
        'review': 'REVIEW',
        'done': 'DONE'
      }
      const dbStatus = statusMap[status] || status.toUpperCase()
      
      onTaskUpdate({
        ...draggedTask,
        status: dbStatus,
        id: draggedTask.id
      })
      setDraggedTask(null)
    }
  }

  return (
    <div className="kanban-board">
      {Object.entries(columns).map(([status, column]) => (
        <div
          key={status}
          className="kanban-column"
          onDragOver={handleDragOver}
          onDrop={() => handleDrop(status)}
          style={{ borderTopColor: `var(${column.color})` }}
        >
          <h3 className="column-title" style={{ color: `var(${column.color})` }}>
            {column.title}
            <span className="task-count">{column.tasks.length}</span>
          </h3>
          <div className="tasks-list">
            {column.tasks.map((task) => (
              <div
                key={task.id}
                className="task-card"
                draggable
                onDragStart={() => handleDragStart(task)}
              >
                <h4 className="task-title">{task.title}</h4>
                <p className="task-description">{task.description}</p>
                {task.assignee && (
                  <div className="task-assignee">
                    <span>Assigned to: {task.assignee}</span>
                  </div>
                )}
                {task.dueDate && (
                  <div className="task-due-date">
                    <span>Due: {new Date(task.dueDate).toLocaleDateString()}</span>
                  </div>
                )}
                <div className="task-priority">
                  <span className={`priority-badge priority-${task.priority?.toLowerCase() || 'medium'}`}>
                    {task.priority || 'Medium'}
                  </span>
                </div>
                <div className="task-actions">
                  <button
                    className="btn btn-small btn-tertiary"
                    onClick={() => onTaskUpdate && onTaskUpdate(task)}
                  >
                    Edit
                  </button>
                  <button
                    className="btn btn-small btn-danger"
                    onClick={() => onTaskDelete && onTaskDelete(task.id)}
                  >
                    Delete
                  </button>
                </div>
              </div>
            ))}
            {column.tasks.length === 0 && (
              <div className="empty-state">
                <p>No tasks</p>
              </div>
            )}
          </div>
        </div>
      ))}
    </div>
  )
}

export default KanbanBoard
