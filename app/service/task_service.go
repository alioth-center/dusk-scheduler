package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/alioth-center/dusk-scheduler/app/config"
	"github.com/alioth-center/dusk-scheduler/app/domain"
	"github.com/alioth-center/dusk-scheduler/app/repository"
	"github.com/alioth-center/dusk-scheduler/infra/logger"
	"time"
)

type taskService struct {
	taskDao          repository.TaskDao
	outcomeDao       repository.OutcomeDao
	taskContentCache repository.TaskContentCache
	sysLogger        logger.Logger
	appConfig        *config.AppConfig
}

func (srv *taskService) CreateTask(ctx context.Context, task *domain.Task, base64Content string) (taskID uint64, err error) {
	// store task content
	if storeErr := srv.taskContentCache.StoreTaskContent(ctx, taskID, bytes.NewBufferString(base64Content)); storeErr != nil {
		srv.sysLogger.ErrorCtx(ctx, "Failed to store task content", storeErr)

		return 0, storeErr
	}

	// initialize task entity and insert into database
	taskEntity := domain.Task{
		Submitter:   task.Submitter,
		ContentHash: task.ContentHash,
		Width:       task.Width,
		Height:      task.Height,
		Type:        task.Type,
		Priority:    task.Priority,
		Format:      task.Format,
		DelayRender: task.DelayRender,
		CreatedAt:   time.Now(),
	}
	createdID, createErr := srv.taskDao.CreateTask(ctx, &taskEntity)
	if createErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to create task for submitter: %d", task.Submitter), createErr)

		return 0, createErr
	}

	return createdID, nil
}

func (srv *taskService) GetTaskByID(ctx context.Context, taskID uint64) (task *domain.Task, exist bool, err error) {
	task, exist, err = srv.taskDao.GetTaskByID(ctx, taskID)
	if err != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to query task by id: %d", taskID), err)

		return nil, false, err
	}

	return task, exist, nil
}

func (srv *taskService) GetCompletedTasksByClientID(ctx context.Context, clientID uint64, statusFilter []string, offsetTaskID uint64) (tasks []*domain.Task, hasMore bool, err error) {
	pageLimit := uint32(srv.appConfig.TaskOptions.ListPageLimit)
	tasks, err = srv.taskDao.GetTaskListByClientID(ctx, clientID, statusFilter, offsetTaskID, pageLimit+1, true)
	if err != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to query tasks by client_id: %d", clientID), err)

		return nil, false, err
	}

	if len(tasks) > int(pageLimit) {
		return tasks[0:pageLimit], true, nil
	}

	return tasks, false, nil
}

func (srv *taskService) CompleteTask(ctx context.Context, taskID uint64) (err error) {
	// update task in database
	if updateErr := srv.taskDao.UpdateTaskAsCompleted(ctx, taskID); updateErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to update task as completed: %s", taskID), updateErr)

		return updateErr
	}

	// delete cached task content
	if deleteErr := srv.taskContentCache.DeleteTaskContent(ctx, taskID); deleteErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to delete task content: %d", taskID), deleteErr)

		return deleteErr
	}

	return nil
}

func (srv *taskService) ArchiveTaskByOutcomeReference(ctx context.Context, outcomeReference string, archiveReason domain.TaskArchiveReason) (exist bool, err error) {
	// get outcome by outcome reference
	outcome, existOutcome, queryOutcomeErr := srv.outcomeDao.GetOutcomeByReference(ctx, outcomeReference)
	if queryOutcomeErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to query outcome by reference: %s", outcomeReference), queryOutcomeErr)

		return false, queryOutcomeErr
	}
	if !existOutcome {
		return false, nil
	}

	// check task existence
	task, existTask, queryTaskErr := srv.taskDao.GetTaskByID(ctx, outcome.TaskID)
	if queryTaskErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to query task by id: %d", outcome.TaskID), err)

		return false, queryTaskErr
	}
	if !existTask {
		return false, nil
	}

	// cannot update task already archived
	if !task.ArchivedAt.IsZero() && task.ArchiveReason != domain.TaskArchiveReasonUnarchived {
		return false, nil
	}
	if updateErr := srv.taskDao.UpdateTaskAsArchived(ctx, task.ID, archiveReason); updateErr != nil {
		srv.sysLogger.ErrorCtx(ctx, fmt.Sprintf("failed to update task as archived: %s", task.ID), updateErr)

		return false, updateErr
	}

	return true, nil
}

func (srv *taskService) GetScheduledTaskListByPainterName(ctx context.Context, painterName string) (list []*domain.Task, content []string, err error) {
	//TODO implement me
	panic("implement me")
}

func (srv *taskService) GetNextScheduledTaskListByPainterName(ctx context.Context, painterName string) (list []*domain.Task, content []string, err error) {
	//TODO implement me
	panic("implement me")
}

func (srv *taskService) FlushPainterScheduler(ctx context.Context, painterName string) {
	//TODO implement me
	panic("implement me")
}

func (srv *taskService) FlushSchedulerTaskQueue(ctx context.Context, taskID uint64) {
	//TODO implement me
	panic("implement me")
}
