package usecase

import (
	"case-study-2/features/task"
	"case-study-2/features/task/dtos"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
)

type service struct {
	model task.Repository
}

func New(model task.Repository) task.Usecase {
	return &service {
		model: model,
	}
}

func (svc *service) FindAll(page, size int) []dtos.ResTask {
	var tasks []dtos.ResTask
	tasksEnt := svc.model.Paginate(page, size)

	for _, t := range tasksEnt {
		var data dtos.ResTask
		_ = smapping.FillStruct(&data, smapping.MapFields(t))

		total := len(t.Subtasks)
		done := 0
		for _, sub := range t.Subtasks {
			if sub.IsCompleted {
				done++
			}
		}
		if total > 0 {
			data.Progress = (done * 100) / total
		} else {
			data.Progress = 100
		}

		if t.Deadline != nil && time.Now().After(*t.Deadline) && !t.IsCompleted {
			data.IsDue = true
		}

		tasks = append(tasks, data)
	}

	return tasks
}

func (svc *service) FindHistory(page, size int) []dtos.ResTask {
	var tasks []dtos.ResTask
	tasksEnt := svc.model.GetHistoryTasks(page, size)

	for _, t := range tasksEnt {
		var data dtos.ResTask
		_ = smapping.FillStruct(&data, smapping.MapFields(t))

		total := len(t.Subtasks)
		done := 0
		for _, sub := range t.Subtasks {
			if sub.IsCompleted {
				done++
			}
		}
		if total > 0 {
			data.Progress = (done * 100) / total
		} else {
			data.Progress = 100
		}
		tasks = append(tasks, data)
	}
	return tasks
}


func (svc *service) FindByID(taskID int) *dtos.ResTask {
	res := dtos.ResTask{}
	task := svc.model.SelectByID(taskID)

	if task == nil {
		return nil
	}

	// Map field utama
	err := smapping.FillStruct(&res, smapping.MapFields(task))
	if err != nil {
		log.Error(err)
		return nil
	}

	for _, sub := range task.Subtasks {
		var resSub dtos.ResSubtask
		if err := smapping.FillStruct(&resSub, smapping.MapFields(sub)); err != nil {
			log.Error(err)
			continue
		}
		res.Subtask = append(res.Subtask, resSub)
	}

	return &res
}


func (svc *service) Create(newTask dtos.InputTask) *dtos.ResTask {
	task := task.Task{}

	err := smapping.FillStruct(&task, smapping.MapFields(newTask))
	if err != nil {
		log.Error(err)
		return nil
	}

	if newTask.Deadline != "" {
		parsedTime, err := time.Parse(time.RFC3339, newTask.Deadline)
		if err != nil {
			log.Error("invalid deadline format:", err)
			return nil
		}
		task.Deadline = &parsedTime
	} else {
		task.Deadline = nil // supaya NULL di database
	}

	taskID := svc.model.Insert(task)
	if taskID == -1 {
		return nil
	}

	resTask := dtos.ResTask{}
	errRes := smapping.FillStruct(&resTask, smapping.MapFields(task))
	if errRes != nil {
		log.Error(errRes)
		return nil
	}

	return &resTask
}


func (svc *service) Modify(taskData dtos.InputTask, taskID int) bool {
	newTask := task.Task{}

	err := smapping.FillStruct(&newTask, smapping.MapFields(taskData))
	if err != nil {
		log.Error(err)
		return false
	}

	newTask.ID = taskID
	rowsAffected := svc.model.Update(newTask)

	if rowsAffected <= 0 {
		log.Error("There is No Task Updated!")
		return false
	}
	
	return true
}

func (svc *service) Remove(taskID int) bool {
	rowsAffected := svc.model.DeleteByID(taskID)

	if rowsAffected <= 0 {
		log.Error("There is No Task Deleted!")
		return false
	}

	return true
}

func (svc *service) FindByTaskID(taskID int) []dtos.ResSubtask {
	var result []dtos.ResSubtask
	subtasks := svc.model.GetByTaskID(taskID)

	for _, sub := range subtasks {
		var data dtos.ResSubtask
		if err := smapping.FillStruct(&data, smapping.MapFields(sub)); err != nil {
			log.Error(err)
		}
		result = append(result, data)
	}
	return result
}

func (svc *service) CreateSubtask(input dtos.InputSubtask) *dtos.ResSubtask {
	subtask := task.SubTask{}
	if err := smapping.FillStruct(&subtask, smapping.MapFields(input)); err != nil {
		log.Error(err)
		return nil
	}

	id := svc.model.InsertSubtask(subtask)
	if id == -1 {
		return nil
	}

	res := dtos.ResSubtask{}
	_ = smapping.FillStruct(&res, smapping.MapFields(subtask))
	return &res
}

func (svc *service) ModifySubtask(input dtos.InputSubtask, subtaskID int) bool {
	subtask := task.SubTask{}

	err := smapping.FillStruct(&subtask, smapping.MapFields(input))
	if err != nil {
		log.Error(err)
		return false
	}

	subtask.ID = subtaskID
	rows := svc.model.UpdateSubtask(subtask)
	if rows <= 0 {
		log.Error("No subtask updated")
		return false
	}

	subtasks := svc.model.GetByTaskID(subtask.TaskID)
	allDone := true
	for _, s := range subtasks {
		if !s.IsCompleted {
			allDone = false
			break
		}
	}

	if allDone {
		mainTask := svc.model.SelectByID(subtask.TaskID)
		mainTask.IsCompleted = true
		svc.model.Update(*mainTask)
	}

	return true
}
