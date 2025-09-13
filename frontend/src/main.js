import './style.css';
import {
    CreateTask,
    CreateTaskWithDetails,
    DeleteTask,
    GetAllTasks,
    GetFilteredTasks,
    ToggleTaskStatus,
    SetTaskPriority,
    SetTaskDueDate,
    UpdateTask
} from '../wailsjs/go/main/App';

class TodoApp {
    constructor() {
        this.tasks = [];
        this.currentFilters = {
            status: 'all',
            priority: 'all',
            dateType: 'all',
            sortField: 'created',
            sortOrder: 'desc'
        };
        this.taskToDelete = null;
        this.init();
    }

    async init() {
        this.setupEventListeners();
        this.setupTheme();
        await this.loadTasks();
        this.render();
    }

    setupEventListeners() {
        document.getElementById('task-form').addEventListener('submit', this.handleTaskSubmit.bind(this));
        document.getElementById('toggle-details').addEventListener('click', this.toggleTaskDetails.bind(this));

        document.getElementById('status-filter').addEventListener('change', this.handleFilterChange.bind(this));
        document.getElementById('priority-filter').addEventListener('change', this.handleFilterChange.bind(this));
        document.getElementById('date-filter').addEventListener('change', this.handleFilterChange.bind(this));
        document.getElementById('sort-field').addEventListener('change', this.handleFilterChange.bind(this));
        document.getElementById('sort-order').addEventListener('change', this.handleFilterChange.bind(this));

        document.getElementById('theme-toggle').addEventListener('click', this.toggleTheme.bind(this));

        document.getElementById('cancel-delete').addEventListener('click', this.hideDeleteModal.bind(this));
        document.getElementById('confirm-delete').addEventListener('click', this.confirmDelete.bind(this));

        document.getElementById('delete-modal').addEventListener('click', (e) => {
            if (e.target.id === 'delete-modal') {
                this.hideDeleteModal();
            }
        });
    }

    setupTheme() {
        const savedTheme = localStorage.getItem('theme') || 'light';
        document.documentElement.setAttribute('data-theme', savedTheme);
        this.updateThemeIcon(savedTheme);
    }

    toggleTheme() {
        const currentTheme = document.documentElement.getAttribute('data-theme');
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';

        document.documentElement.setAttribute('data-theme', newTheme);
        localStorage.setItem('theme', newTheme);
        this.updateThemeIcon(newTheme);
    }

    updateThemeIcon(theme) {
        const icon = document.querySelector('.theme-icon');
        icon.textContent = theme === 'dark' ? '☀️' : '🌙';
    }

    toggleTaskDetails() {
        const details = document.getElementById('task-details');
        const button = document.getElementById('toggle-details');
        const isVisible = details.style.display !== 'none';

        details.style.display = isVisible ? 'none' : 'flex';
        button.textContent = isVisible ? '+ More Options' : '- Less Options';
    }

    async handleTaskSubmit(e) {
        e.preventDefault();

        const title = document.getElementById('task-title').value.trim();
        if (!title) return;

        const description = document.getElementById('task-description').value.trim();
        const priority = document.getElementById('task-priority').value;
        const dueDateValue = document.getElementById('task-due-date').value;

        try {
            let newTask;
            if (description || priority !== 'medium' || dueDateValue) {
                const dueDate = dueDateValue ? new Date(dueDateValue) : null;
                newTask = await CreateTaskWithDetails(title, description, priority, dueDate);
            } else {
                newTask = await CreateTask(title, description);
            }

            document.getElementById('task-form').reset();
            this.toggleTaskDetails();

            await this.loadTasks();
            this.render();

        } catch (error) {
            console.error('Error creating task:', error);
            this.showError('Failed to create task');
        }
    }

    async handleFilterChange() {
        this.currentFilters = {
            status: document.getElementById('status-filter').value,
            priority: document.getElementById('priority-filter').value,
            dateType: document.getElementById('date-filter').value,
            sortField: document.getElementById('sort-field').value,
            sortOrder: document.getElementById('sort-order').value
        };

        await this.loadTasks();
        this.render();
    }

    async loadTasks() {
        try {
            const { status, priority, dateType, sortField, sortOrder } = this.currentFilters;
            this.tasks = await GetFilteredTasks(status, priority, dateType, sortField, sortOrder);
        } catch (error) {
            console.error('Error loading tasks:', error);
            this.tasks = [];
        }
    }

    async toggleTask(taskId) {
        try {
            await ToggleTaskStatus(taskId);
            await this.loadTasks();
            this.render();
        } catch (error) {
            console.error('Error toggling task:', error);
        }
    }

    showDeleteModal(taskId) {
        this.taskToDelete = taskId;
        document.getElementById('delete-modal').style.display = 'flex';
    }

    hideDeleteModal() {
        this.taskToDelete = null;
        document.getElementById('delete-modal').style.display = 'none';
    }

    async confirmDelete() {
        if (!this.taskToDelete) return;

        try {
            await DeleteTask(this.taskToDelete);
            await this.loadTasks();
            this.render();
            this.hideDeleteModal();
        } catch (error) {
            console.error('Error deleting task:', error);
        }
    }

    formatDate(dateString) {
        const date = new Date(dateString);
        const now = new Date();
        const diffTime = date - now;
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

        const formatter = new Intl.DateTimeFormat('en-US', {
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });

        let formatted = formatter.format(date);

        if (diffDays < 0) {
            formatted += ' (overdue)';
        } else if (diffDays === 0) {
            formatted += ' (today)';
        } else if (diffDays === 1) {
            formatted += ' (tomorrow)';
        }

        return formatted;
    }

    isOverdue(task) {
        if (!task.due_date || task.status === 'completed') return false;
        return new Date(task.due_date) < new Date();
    }

    renderTask(task) {
        const isOverdue = this.isOverdue(task);
        const isCompleted = task.status === 'completed';

        return `
            <div class="task-item ${isCompleted ? 'completed' : ''} ${isOverdue ? 'overdue' : ''}">
                <div class="task-header">
                    <input
                        type="checkbox"
                        class="task-checkbox"
                        ${isCompleted ? 'checked' : ''}
                        onchange="todoApp.toggleTask('${task.id}')"
                    >
                    <div class="task-content">
                        <div class="task-title">${task.title}</div>
                        ${task.description ? `<div class="task-description">${task.description}</div>` : ''}
                        <div class="task-meta">
                            <span class="priority-badge priority-${task.priority}">${task.priority}</span>
                            ${task.due_date ? `<span class="due-date ${isOverdue ? 'overdue' : ''}">${this.formatDate(task.due_date)}</span>` : ''}
                            <span class="created-date">Created ${this.formatDate(task.created_at)}</span>
                        </div>
                    </div>
                    <div class="task-actions">
                        <button class="btn btn-danger" onclick="todoApp.showDeleteModal('${task.id}')">
                            Delete
                        </button>
                    </div>
                </div>
            </div>
        `;
    }

    render() {
        const activeTasks = this.tasks.filter(task => task.status === 'active');
        const completedTasks = this.tasks.filter(task => task.status === 'completed');

        document.getElementById('active-count').textContent = `${activeTasks.length} active`;
        document.getElementById('completed-count').textContent = `${completedTasks.length} completed`;
        document.getElementById('total-count').textContent = `${this.tasks.length} total`;
        document.getElementById('active-counter').textContent = activeTasks.length;
        document.getElementById('completed-counter').textContent = completedTasks.length;

        const activeTasksContainer = document.getElementById('active-tasks');
        if (activeTasks.length > 0) {
            activeTasksContainer.innerHTML = activeTasks.map(task => this.renderTask(task)).join('');
        } else {
            activeTasksContainer.innerHTML = '<div class="empty-state"><p>No active tasks</p></div>';
        }

        const completedTasksContainer = document.getElementById('completed-tasks');
        if (completedTasks.length > 0) {
            completedTasksContainer.innerHTML = completedTasks.map(task => this.renderTask(task)).join('');
        } else {
            completedTasksContainer.innerHTML = '<div class="empty-state"><p>No completed tasks</p></div>';
        }

        const emptyState = document.getElementById('empty-state');
        if (this.tasks.length === 0) {
            emptyState.style.display = 'block';
        } else {
            emptyState.style.display = 'none';
        }

        document.getElementById('active-tasks-section').style.display =
            this.currentFilters.status === 'completed' ? 'none' : 'block';
        document.getElementById('completed-tasks-section').style.display =
            this.currentFilters.status === 'active' ? 'none' : 'block';
    }

    showError(message) {
        console.error(message);
        alert(message);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    window.todoApp = new TodoApp();
});
