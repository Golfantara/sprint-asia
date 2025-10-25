package repository

import (
	"case-study-2/features/task"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) task.Repository {
	return &model {
		db: db,
	}
}

func (mdl *model) Paginate(page, size int) []task.Task {
	var tasks []task.Task
	offset := (page - 1) * size

	result := mdl.db.Preload("Subtasks").Offset(offset).Limit(size).Find(&tasks)
	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return tasks
}
func (mdl *model) Insert(newTask task.Task) int64 {
	result := mdl.db.Create(&newTask)

	if result.Error != nil {
		log.Error(result.Error)
		return -1
	}

	return int64(newTask.ID)
}

func (mdl *model) SelectByID(taskID int) *task.Task {
	var task task.Task
	result := mdl.db.Preload("Subtasks").First(&task, taskID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &task
}

func (mdl *model) Update(task task.Task) int64 {
	result := mdl.db.Updates(&task)

	if result.Error != nil {
		log.Error(result.Error)
	}

	return result.RowsAffected
}

func (mdl *model) DeleteByID(taskID int) int64 {
	result := mdl.db.Delete(&task.Task{}, taskID)
	
	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
}

func (mdl *model) GetByTaskID(taskID int) []task.SubTask {
	var subtasks []task.SubTask
	result := mdl.db.Where("task_id = ?", taskID).Find(&subtasks)
	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}
	return subtasks
}

func (mdl *model) InsertSubtask(newSubtask task.SubTask) int64 {
	result := mdl.db.Create(&newSubtask)
	if result.Error != nil {
		log.Error(result.Error)
		return -1
	}
	return int64(newSubtask.ID)
}

func (mdl *model) GetHistoryTasks(page, size int) []task.Task {
	var tasks []task.Task
	offset := (page - 1) * size

	result := mdl.db.Preload("Subtasks").Where("is_completed = ?", true).
		Offset(offset).Limit(size).Find(&tasks)
	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return tasks
}

func (mdl *model) UpdateSubtask(subtask task.SubTask) int64 {
	result := mdl.db.Model(&task.SubTask{}).Where("id = ?", subtask.ID).Updates(subtask)
	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}
	return result.RowsAffected
}
